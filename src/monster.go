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
}

const monsterSize = cellSize / 4

func (g *Game) addMonster(kind string, x, y int) {
	cx := float64(x)*cellSize + cellSize/2
	cy := float64(y)*cellSize + cellSize/2

	angle := rand.Float64() * 2 * math.Pi
	speed := rand.Float64()*0.5 + 0.5

	s := g.addSprite(kind, cx, cy, angle, speed, monsterSize)

	id := rand.Uint64()
	mon := &Monster{
		id:     id,
		sprite: s,
		health: 100,
	}
	if kind == "skeleton" {
		mon.health = 35
	}
	if kind == "ghoul" {
		mon.health = 75
	}

	g.monsters[id] = mon
}

func (m *Monster) checkWallCollision(x, y float64) (int, int, int) {
	size := m.sprite.size
	if wall, x, y := game.getWallAt(x+size, y); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := game.getWallAt(x-size, y); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := game.getWallAt(x, y+size); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := game.getWallAt(x, y-size); wall > 0 {
		return wall, x, y
	}
	return 0, 0, 0
}

func (g *Game) updateMonsters() {
	for _, mon := range g.monsters {
		sprite := mon.sprite

		// Animations!
		if g.ticks%20 == 0 {
			if spriteImages[sprite.kind+"-1"] != nil {
				if mon.sprite.image == spriteImages[sprite.kind+"-1"] {
					mon.sprite.image = spriteImages[sprite.kind]
				} else {
					mon.sprite.image = spriteImages[sprite.kind+"-1"]
				}
			}
		}

		if sprite.speed <= 0 {
			continue
		}

		newX := sprite.x + math.Cos(sprite.angle)*sprite.speed
		newY := sprite.y + math.Sin(sprite.angle)*sprite.speed
		if wi, _, _ := g.getWallAt(newX, newY); wi > 0 {
			mon.sprite.angle = mon.sprite.angle + math.Pi
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
