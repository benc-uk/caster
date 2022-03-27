package main

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Projectile struct {
	id     uint64
	sprite *Sprite
	damage int
}

func (g *Game) addProjectile(kind string, x, y float64, angle float64, speed float64, damage int) {
	s := g.addSprite(kind, x, y, angle, speed, cellSize/16.0)
	s.alpha = 0.5

	id := rand.Uint64()
	g.projectiles[id] = &Projectile{
		id:     id,
		sprite: s,
		damage: damage,
	}
}

func (g *Game) updateProjectiles() {
	if g.fc%20 != 0 {
		si := spriteImages["magic_1"]
		sir := ebiten.NewImage(32, 32)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-16, -16)
		op.GeoM.Rotate(math.Pi / 4)
		op.GeoM.Translate(16, 16)

		sir.Clear()
		sir.DrawImage(si, op)
		spriteImages["magic_1"] = sir
	}

	for id := range g.projectiles {
		sprite := g.projectiles[id].sprite
		if sprite.speed <= 0 {
			continue
		}

		newX := sprite.x + math.Cos(sprite.angle)*sprite.speed
		newY := sprite.y + math.Sin(sprite.angle)*sprite.speed
		if wi, _, _ := g.getWallAt(newX, newY); wi > 0 {
			g.removeProjectile(g.projectiles[id])
		}

		// Check if it hit a monster
		for _, m := range g.monsters {
			if m.sprite.isHit(newX, newY) {
				if m == nil {
					continue
				}
				m.health -= g.projectiles[id].damage
				if m.health <= 0 {
					playSound("monster_death", 1.0, false)
					g.removeMonster(m)
				} else {
					playSound("monster_hit", 1.0, false)
				}

				g.removeProjectile(g.projectiles[id])

			}
		}

		sprite.x = newX
		sprite.y = newY
	}
}

func (g *Game) removeProjectile(p *Projectile) {
	delete(g.projectiles, p.id)
	g.removeSprite(p.sprite)
	p.sprite = nil
	p = nil
}
