package main

import (
	"log"
	"math/rand"
)

type Monster struct {
	id     uint64
	sprite *Sprite
	health int
}

func (g *Game) addMonster(kind string, x, y float64, angle float64, speed float64, size float64) {
	log.Printf("!!!!!Adding monster: %s", kind)
	s := g.addSprite(kind, x, y, angle, speed, size)

	id := rand.Uint64()
	g.monsters[id] = &Monster{
		id:     id,
		sprite: s,
		health: 22,
	}

	log.Printf("Added projectile: %+v", g.projectiles)
}

func (g *Game) removeMonster(m *Monster) {
	delete(g.monsters, m.id)
	g.removeSprite(m.sprite)
	m.sprite = nil
	m = nil
}
