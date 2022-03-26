package main

import (
	"log"
	"math/rand"
)

type Projectile struct {
	id     uint64
	sprite *Sprite
	damage float64
}

func (g *Game) addProjectile(kind string, x, y float64, angle float64, speed float64, damage float64) {
	s := g.addSprite(kind, x, y, angle, speed, cellSize/16.0)

	id := rand.Uint64()
	g.projectiles[id] = &Projectile{
		id:     id,
		sprite: s,
		damage: damage,
	}

	log.Printf("Added projectile: %+v", g.projectiles)
}

func (g *Game) removeProjectile(p *Projectile) {
	delete(g.projectiles, p.id)
	g.removeSprite(p.sprite)
	p.sprite = nil
	p = nil
}
