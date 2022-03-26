package main

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

var hudImage *ebiten.Image

func renderHud(screen *ebiten.Image, g *Game) {
	// Update the HUD but only every 20 frames
	if g.fc%20 == 0 {
		hudImage.Clear()

		healthStr := fmt.Sprintf("%d - Health", g.player.health)
		healthOp := &ebiten.DrawImageOptions{}
		healthOp.GeoM.Scale(magicSprite/12, magicSprite/12)
		healthOp.GeoM.Translate(float64(hudMargin), float64(winHeight-hudMargin))
		healthOp.ColorM.Scale(0.0, 0.0, 0.0, 0.5)
		text.DrawWithOptions(hudImage, healthStr, gameFont, healthOp)
		healthOp.GeoM.Translate(-2.5, -2.5)
		healthOp.ColorM.Reset()
		healthOp.ColorM.Scale(0.889, 0.141, 0.188, 1)
		text.DrawWithOptions(hudImage, healthStr, gameFont, healthOp)

		manaStr := fmt.Sprintf("Mana - %d", g.player.mana)
		manaOp := &ebiten.DrawImageOptions{}
		manaOp.GeoM.Scale(magicSprite/12, magicSprite/12)
		manaOp.GeoM.Translate(float64(winWidth-hudMargin*12), float64(winHeight-hudMargin))
		manaOp.ColorM.Scale(0.0, 0.0, 0.0, 0.5)
		text.DrawWithOptions(hudImage, manaStr, gameFont, manaOp)
		manaOp.GeoM.Translate(-2.5, -2.5)
		manaOp.ColorM.Reset()
		manaOp.ColorM.Scale(0.094, 0.623, 0.984, 1.0)
		text.DrawWithOptions(hudImage, manaStr, gameFont, manaOp)

		screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
	} else {
		screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
	}
}

// ===========================================================
// Handle the map overlay
// ===========================================================
func (g *Game) overlay(screen *ebiten.Image) {
	if !overlayShown {
		return
	}

	overlayImage.Fill(color.RGBA{0x00, 0x00, 0x00, 0x00})
	px := g.player.x / float64(cellSize/overlayCellSize)
	py := g.player.y / float64(cellSize/overlayCellSize)

	// Draw the player
	ebitenutil.DrawRect(overlayImage, px-1, py-1, 3, 3, color.RGBA{255, 255, 255, 255})

	// draw sprites
	for _, sprite := range g.sprites {
		sx := sprite.x / float64(cellSize/overlayCellSize)
		sy := sprite.y / float64(cellSize/overlayCellSize)
		c := color.RGBA{255, 0, 0, 255}
		ebitenutil.DrawRect(overlayImage, sx-1, sy-1, 3, 3, c)
	}

	// Draw the map
	for y := 0; y < g.mapHeight; y++ {
		for x := 0; x < g.mapWidth; x++ {
			if g.mapdata[x][y] != 0 {
				ebitenutil.DrawRect(overlayImage, float64(x*overlayCellSize), float64(y*overlayCellSize), float64(overlayCellSize), float64(overlayCellSize), color.RGBA{255, 255, 255, 58})
			}
		}
	}

	w, h := overlayImage.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(winWidth)/float64(w)*overlayZoom, float64(winHeight)/float64(h)*overlayZoom)
	screen.DrawImage(overlayImage, op)
}
