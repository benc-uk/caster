package main

import (
	"math"
	"math/rand"
	"time"
)

type Monster struct {
	id     uint64
	sprite *Sprite
	health int
	damage int
}

func (g *Game) addMonster(kind string, x, y int) {
	const monsterSize = float64(cellSize) / 6.0

	cx := float64(x)*cellSize + cellSize/2
	cy := float64(y)*cellSize + cellSize/2

	angle := rand.Float64() * 2 * math.Pi
	speed := rand.Float64()*0.5 + 0.5

	// We hope that 64 bit ints are unique enough
	id := rand.Uint64()
	mon := &Monster{
		id:     id,
		sprite: g.addSprite("monsters/"+kind, cx, cy, angle, speed, monsterSize),
		health: 10,
		damage: 10,
	}
	if kind == "skeleton" {
		mon.health = 35
	}
	if kind == "ghoul" {
		mon.health = 75
	}

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

		if sprite.speed <= 0 {
			continue
		}

		newX := sprite.x + math.Cos(sprite.angle)*sprite.speed
		newY := sprite.y + math.Sin(sprite.angle)*sprite.speed

		if hitWall := mon.checkWallCollision(newX, newY); hitWall != nil {
			mon.sprite.angle = mon.sprite.angle + math.Pi
		}

		if dist := mon.getDistanceToPlayer(g.player); dist < (g.player.size*3 + sprite.size) {
			mon.sprite.angle = mon.sprite.angle + math.Pi
			g.player.damage(mon.damage)
			bounceX := sprite.x + math.Cos(sprite.angle)*(g.player.size+sprite.size+0.1)
			bounceY := sprite.y + math.Sin(sprite.angle)*(g.player.size+sprite.size+0.1)
			if hitWall := mon.checkWallCollision(bounceX, bounceY); hitWall == nil {
				newX = bounceX
				newY = bounceY
			}
		}

		sprite.x = newX
		sprite.y = newY
	}
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

func (m *Monster) getDistanceToPlayer(p Player) float64 {
	dx := p.x - m.sprite.x
	dy := p.y - m.sprite.y
	return math.Sqrt(dx*dx + dy*dy)
}
