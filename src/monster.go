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
	projectileDelay  int
}

func (g *Game) addMonster(kind string, x, y int) {
	const monsterSize = float64(cellSize) / 6.0
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
		projectileDelay:  60,
	}

	if kind == "skeleton" {
		mon.health = 35
		mon.canShoot = true
		mon.baseSpeed = rand.Float64()*0.5 + 0.5
		mon.projectileDamage = 6
		mon.projectileKind = "bone"
		mon.projectileDelay = 180
	}

	if kind == "thing" {
		mon.health = 100
		mon.canShoot = true
		mon.baseSpeed = rand.Float64()*0.5 + 0.5
		mon.projectileDamage = 12
		mon.projectileKind = "slime"
		mon.projectileDelay = 240
	}

	if kind == "ghoul" {
		mon.health = 75
		mon.meleeDamage = 25
		mon.baseSpeed = rand.Float64()*0.4 + 0.2
	}

	if kind == "orc" {
		mon.health = 50
		mon.meleeDamage = 10
		mon.baseSpeed = rand.Float64()*0.5 + 1
	}

	mon.sprite.speed = mon.baseSpeed
	g.monsters[mon.id] = mon
}

func (m *Monster) checkWallCollision(x, y float64) *Wall {
	size := m.sprite.size
	if wall := game.getWallAt(x+size, y); wall != nil {
		return wall
	}
	if wall := game.getWallAt(x-size, y); wall != nil {
		return wall
	}
	if wall := game.getWallAt(x, y+size); wall != nil {
		return wall
	}
	if wall := game.getWallAt(x, y-size); wall != nil {
		return wall
	}
	return nil
}

func (g *Game) updateMonsters() {
	for _, mon := range g.monsters {
		sprite := mon.sprite

		// Animations!
		if g.ticks%20 == 0 {
			if imageCache[sprite.kind+"-1"] != nil {
				if mon.sprite.image == imageCache[sprite.kind+"-1"] {
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
			sprite.speed = -mon.baseSpeed
		}

		if mon.state == MonsterStateAttack {
			var angleToPlayer float64
			var see bool
			if see, angleToPlayer = mon.checkLosToPlayer(g.player); !see {
				mon.state = MonsterStateIdle
				continue
			}
			sx := sprite.x + math.Cos(angleToPlayer)*32
			sy := sprite.y + math.Sin(angleToPlayer)*32
			game.addProjectile(mon.projectileKind, sx, sy, angleToPlayer, (float64(cellSize) / 9.0), mon.projectileDamage, 1)
			mon.state = MonsterStateDoNothing
			mon.stateTicker = mon.projectileDelay
		}

		if mon.state == MonsterStateMelee {
			var angleToPlayer float64
			var see bool
			if see, angleToPlayer = mon.checkLosToPlayer(g.player); !see {
				mon.state = MonsterStateIdle
				continue
			}
			sprite.speed = mon.baseSpeed
			sprite.angle = angleToPlayer
		}

		//logMessage("Monster", mon.id, "state", mon.state)

		// Move the monster
		newX := sprite.x + math.Cos(sprite.angle)*sprite.speed
		newY := sprite.y + math.Sin(sprite.angle)*sprite.speed
		// Check if the hit a wall
		if wall := mon.checkWallCollision(newX, newY); wall == nil {
			// Check if they move into the player
			if distPlayer := mon.sprite.getDistanceToPlayer(g.player); distPlayer < (g.player.size*3 + sprite.size) {
				g.player.damage(mon.meleeDamage)
				mon.state = MonsterStateRecoil
				mon.stateTicker = 45
			}
			sprite.x = newX
			sprite.y = newY
		} else {
			sprite.speed = -sprite.speed
		}
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
}

func (m *Monster) kill() {
	s := game.addSprite(m.sprite.kind+"-dead", m.sprite.x, m.sprite.y, 0, 0, 0)
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
