package main

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var hudImage *ebiten.Image

func renderHud(screen *ebiten.Image, g *Game) {
	// Update the HUD but only every 20 frames
	if g.fc%20 == 0 {
		hudImage.Clear()

		healthStr := fmt.Sprintf("%d - Health", g.player.health)
		healthR := text.BoundString(gameFont, healthStr)
		textOp := &ebiten.DrawImageOptions{}
		textOp.GeoM.Scale(magicSprite/12, magicSprite/12)
		textOp.GeoM.Translate(float64(hudMargin), float64(winHeight-healthR.Max.Y-hudMargin))
		textOp.ColorM.Scale(0.1, 0.9, 0.1, 0.6)
		text.DrawWithOptions(hudImage, healthStr, gameFont, textOp)

		manaStr := fmt.Sprintf("Mana - %d", g.player.mana)
		manaR := text.BoundString(gameFont, manaStr)
		textOp = &ebiten.DrawImageOptions{}
		textOp.GeoM.Scale(magicSprite/12, magicSprite/12)
		textOp.GeoM.Translate(float64(winWidth-manaR.Max.X), float64(winHeight-manaR.Max.Y-hudMargin))
		textOp.ColorM.Scale(0.3, 0.2, 0.9, 0.8)
		text.DrawWithOptions(hudImage, manaStr, gameFont, textOp)
		screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
	} else {
		screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
	}
}
