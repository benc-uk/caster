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

func (g *Game) addProjectile(kind string, x, y float64, angle float64, speed float64, damage int, alpha float64) {
	s := g.addSprite("effects/"+kind, x, y, angle, speed, cellSize/16.0)
	s.alpha = alpha

	id := rand.Uint64()
	g.projectiles[id] = &Projectile{
		id:     id,
		sprite: s,
		damage: damage,
	}
}

func (g *Game) updateProjectiles() {
	for id := range g.projectiles {
		sprite := g.projectiles[id].sprite

		// Animate and rotate the projectile sprite every 5 frames
		if g.ticks%5 == 0 {
			rotatedImg := ebiten.NewImageFromImage(sprite.image)
			rotateOp := &ebiten.DrawImageOptions{}
			rotateOp.GeoM.Translate(-spriteImgSizeH, -spriteImgSizeH)
			rotateOp.GeoM.Rotate(math.Pi / 4)
			rotateOp.GeoM.Translate(spriteImgSizeH, spriteImgSizeH)
			rotatedImg.Clear()
			rotatedImg.DrawImage(sprite.image, rotateOp)
			sprite.image = rotatedImg
		}

		newX := sprite.x + math.Cos(sprite.angle)*sprite.speed
		newY := sprite.y + math.Sin(sprite.angle)*sprite.speed
		if wall := g.getWallAt(newX, newY); wall != nil {
			g.removeProjectile(g.projectiles[id])
		}

		// Check if it hit a monster
		for _, m := range g.monsters {
			if m.sprite.isHit(newX, newY) {
				if m == nil || g.projectiles[id] == nil {
					continue
				}
				m.damage(g.projectiles[id].damage)
				g.removeProjectile(g.projectiles[id])
			}
		}

		// Check if it hit the player
		playerDist := sprite.getDistanceToPlayer(g.player)
		if playerDist < (g.player.size*3 + sprite.size) {
			g.player.damage(g.projectiles[id].damage)
			g.removeProjectile(g.projectiles[id])
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
