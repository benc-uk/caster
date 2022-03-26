package main

import (
	"math"
	"math/rand"
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
	g.monsters[id] = &Monster{
		id:     id,
		sprite: s,
		health: 100,
	}
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
	for id := range g.monsters {
		sprite := g.monsters[id].sprite
		if sprite.speed <= 0 {
			continue
		}

		newX := sprite.x + math.Cos(sprite.angle)*sprite.speed
		newY := sprite.y + math.Sin(sprite.angle)*sprite.speed
		if wi, _, _ := g.getWallAt(newX, newY); wi > 0 {
			g.monsters[id].sprite.angle = g.monsters[id].sprite.angle + math.Pi
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
