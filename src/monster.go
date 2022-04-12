package main

import (
	"math"
	"math/rand"
	"time"
)

type MonsterState int

const (
	MonsterStateIdle MonsterState = iota
	MonsterStateMelee
	MonsterStateAttack
	MonsterStateRecoil
	MonsterStateWander
	MonsterStateDoNothing
)

type Monster struct {
	id               uint64
	sprite           *Sprite
	health           int
	meleeDamage      int
	projectileDamage int
	state            MonsterState
	stateTicker      int
	baseSpeed        float64
	canShoot         bool
	projectileKind   string
	projectileProb   float64
	projectileSpeed  float64
	seenPlayer       bool
}

func (g *Game) addMonster(kind string, x, y int) {
	const monsterSize = float64(cellSize) / 4
	cx := float64(x)*cellSize + cellSize/2
	cy := float64(y)*cellSize + cellSize/2
	angle := rand.Float64() * 2 * math.Pi

	// We hope that 64 bit ints are unique enough
	id := rand.Uint64()
	mon := &Monster{
		id:               id,
		sprite:           g.addSprite("monsters/"+kind, cx, cy, angle, 1, monsterSize),
		health:           10,
		meleeDamage:      10,
		projectileDamage: 10,
		state:            MonsterStateIdle,
		baseSpeed:        1,
		canShoot:         false,
		projectileKind:   "slime",
		stateTicker:      1,
		projectileProb:   0.0,
		projectileSpeed:  (float64(cellSize) / 9.0),
	}

	if kind == "skeleton" {
		mon.health = 35
		mon.canShoot = true
		mon.baseSpeed = rand.Float64()*0.5 + 0.5
		mon.projectileDamage = 6
		mon.projectileKind = "bone"
		mon.projectileProb = 60.0
	}

	if kind == "thing" {
		mon.health = 100
		mon.canShoot = true
		mon.baseSpeed = rand.Float64()*0.5 + 0.5
		mon.projectileDamage = 12
		mon.projectileKind = "slime"
		mon.projectileProb = 50.0
		mon.projectileSpeed = mon.projectileSpeed * 0.8
	}

	if kind == "ghoul" {
		mon.health = 75
		mon.meleeDamage = 30
		mon.baseSpeed = rand.Float64()*0.6 + 0.3
	}

	if kind == "orc" {
		mon.health = 50
		mon.meleeDamage = 10
		mon.baseSpeed = rand.Float64()*0.8 + 1
	}

	if kind == "wiz" {
		mon.health = 50
		mon.meleeDamage = 10
		mon.baseSpeed = rand.Float64()*0.4 + 0.3
		mon.canShoot = true
		mon.projectileDamage = 18
		mon.projectileKind = "fireball"
		mon.projectileProb = 45.0
		mon.projectileSpeed = mon.projectileSpeed * 0.6
	}

	if kind == "spectre" {
		mon.health = 120
		mon.meleeDamage = 15
		mon.baseSpeed = rand.Float64()*0.6 + 1
		mon.sprite.alpha = 0.4
	}

	mon.sprite.speed = mon.baseSpeed
	g.monsters[mon.id] = mon
	g.stats.monsters++
}

func (m *Monster) checkWallCollision(x, y float64) (*Wall, float64, float64) {
	size := m.sprite.size
	if wall := game.getWallAt(x+size, y); wall != nil {
		return wall, x + size, y
	}
	if wall := game.getWallAt(x-size, y); wall != nil {
		return wall, x - size, y
	}
	if wall := game.getWallAt(x, y+size); wall != nil {
		return wall, x, y + size
	}
	if wall := game.getWallAt(x, y-size); wall != nil {
		return wall, x, y - size
	}
	return nil, 0, 0
}

func (g *Game) updateMonsters() {
	for _, mon := range g.monsters {
		sprite := mon.sprite
		playerDist := sprite.getDistanceToPlayer(g.player)
		if playerDist > viewDistance {
			continue
		}

		// Animations!
		if g.ticks%20 == 0 {
			if imageCache[sprite.kind+"-1"] != nil {
				if sprite.image == imageCache[sprite.kind+"-1"] {
					mon.sprite.image = imageCache[sprite.kind]
				} else {
					mon.sprite.image = imageCache[sprite.kind+"-1"]
				}
			}
		}

		// Handle timed state transitions
		if mon.stateTicker > 0 {
			mon.stateTicker--

			if mon.stateTicker <= 0 {
				mon.state = MonsterStateIdle
				mon.stateTicker = 0
				// random angle
				sprite.angle = rand.Float64() * 2 * math.Pi
			}
		}

		if mon.state == MonsterStateDoNothing {
		}

		if mon.state == MonsterStateIdle {
			if see, _ := mon.checkLosToPlayer(g.player); see {
				if mon.canShoot {
					mon.state = MonsterStateAttack
				} else {
					mon.state = MonsterStateMelee
				}
			}
		}

		if mon.state == MonsterStateRecoil {
			sprite.speed = -(mon.baseSpeed * 1.2)
		}

		if mon.state == MonsterStateAttack {
			var angleToPlayer float64
			var see bool
			if see, angleToPlayer = mon.checkLosToPlayer(g.player); !see {
				mon.seenPlayer = false
				mon.state = MonsterStateIdle
				continue
			}
			if rand.Float64()*10000.0 <= mon.projectileProb {
				sx := sprite.x + math.Cos(angleToPlayer)*32
				sy := sprite.y + math.Sin(angleToPlayer)*32
				game.addProjectile(mon.projectileKind, sx, sy, angleToPlayer, mon.projectileSpeed, mon.projectileDamage, 1)
				playSound("whoosh", 1, false)
			}
		}

		if mon.state == MonsterStateMelee {
			var angleToPlayer float64
			var see bool
			if see, angleToPlayer = mon.checkLosToPlayer(g.player); !see {
				mon.seenPlayer = false
				mon.state = MonsterStateIdle
				continue
			}
			sprite.speed = mon.baseSpeed
			sprite.angle = angleToPlayer
		}

		// Move the monster
		newX := sprite.x + math.Cos(sprite.angle)*sprite.speed
		newY := sprite.y + math.Sin(sprite.angle)*sprite.speed

		// Check if it hits a wall
		if wall, _, _ := mon.checkWallCollision(newX, newY); wall == nil {
			// Check if they move into the player
			if playerDist < (g.player.size*3+sprite.size) && mon.state != MonsterStateRecoil {
				playSound("monster_attack", 1, false)
				g.player.damage(mon.meleeDamage)
				mon.state = MonsterStateRecoil
				mon.stateTicker = 45
			}
		} else {
			// This weird code, stops monsters getting stuck on walls when walking towards the player
			sprite.speed = mon.baseSpeed
			mon.state = MonsterStateDoNothing
			mon.stateTicker = 5
			mon.sprite.angle += math.Pi / 4
		}
		sprite.x = newX
		sprite.y = newY
	}
}

func (m *Monster) checkLosToPlayer(p Player) (canSee bool, angle float64) {
	playerDist := m.sprite.getDistanceToPlayer(game.player)
	wall, dist, a := fireRayAt(m.sprite.x, m.sprite.y, p.x, p.y, playerDist)

	if wall == nil {
		if dist >= viewDistance {
			return false, 0
		}
		if dist <= playerDist {
			if !m.seenPlayer {
				playSound("monster_grunt", 1, false)
			}
			m.seenPlayer = true
			return true, a
		}
	}
	return false, 0
}

func (g *Game) removeMonster(m *Monster) {
	delete(g.monsters, m.id)
	g.removeSprite(m.sprite)
	m.sprite = nil
	m = nil
	g.stats.kills++
}

func (m *Monster) kill() {
	s := game.addSprite(m.sprite.kind+"-dead", m.sprite.x, m.sprite.y, 0, 0, 0)
	s.alpha = m.sprite.alpha
	game.removeMonster(m)
	time.AfterFunc(time.Millisecond*300, func() {
		game.removeSprite(s)
	})
}

func (s *Sprite) getDistanceToPlayer(p Player) float64 {
	dx := p.x - s.x
	dy := p.y - s.y
	return math.Sqrt(dx*dx + dy*dy)
}

func (m *Monster) damage(d int) {
	m.health -= d
	if m.health <= 0 {
		playSound("monster_death", 1.0, false)
		m.kill()
	} else {
		playSound("monster_hit", 1.0, false)
	}
}
