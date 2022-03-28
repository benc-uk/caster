package main

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var hudImage *ebiten.Image
var gameFont font.Face

func initHUD() {
	// Font(s)
	fontData, err := os.ReadFile("./fonts/morris-roman.ttf")
	if err != nil {
		log.Fatal(err)
	}

	ttFont, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}

	gameFont, err = opentype.NewFace(ttFont, &opentype.FaceOptions{
		Size:    100.0 * (float64(magicSprite) / 12.0),
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func renderTitle(screen *ebiten.Image) {

	for x := 0; x < winWidth; x += cellSize * 2 {
		for y := 0; y < winHeight; y += cellSize * 2 {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(2, 2)
			op.GeoM.Translate(float64(x), float64(y))
			screen.DrawImage(imageCache[fmt.Sprintf("walls/%s", (x+y)%9+1)], op)
		}
	}

	ebitenutil.DrawRect(screen, 0, 0, float64(winWidth), float64(winHeight), color.RGBA{0, 0, 0, 160})

	msg := "Crypt Caster"
	textRect := text.BoundString(gameFont, msg)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2.0, 2.0)
	op.Filter = ebiten.FilterLinear
	op.ColorM.Reset()
	op.ColorM.Scale(0, 0, 0, 0.5)
	op.GeoM.Translate(float64(winWidth/3)-float64(textRect.Dx())/2.0, float64(winHeight/3)-float64(textRect.Dy())/2.0)
	text.DrawWithOptions(screen, msg, gameFont, op)
	op.ColorM.Reset()
	op.ColorM.Scale(1.5, 0.3, 0.1, 1)
	op.GeoM.Translate(-(magicSprite / 2.0), -(magicSprite / 2.0))
	text.DrawWithOptions(screen, msg, gameFont, op)

	op = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(5.0, 5.0)
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(float64(winWidth)-(38*magicSprite), float64(winHeight/3)-float64(textRect.Dy())/2.0-(18*magicSprite))
	screen.DrawImage(imageCache["items/ball"], op)

	msg = fmt.Sprintf("%d. %s", titleLevelIndex+1, titleLevels[titleLevelIndex])
	textRect = text.BoundString(gameFont, msg)
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.ColorM.Reset()
	op.ColorM.Scale(0.1, 0.8, 0.1, 1)
	op.GeoM.Translate(float64(winWidth/2)-float64(textRect.Dx())/2.0, float64(winHeight/2)-float64(textRect.Dy())/2.0)
	text.DrawWithOptions(screen, msg, gameFont, op)

	msg = "Press enter to start\n  Press esc to quit"
	textRect = text.BoundString(gameFont, msg)
	op = &ebiten.DrawImageOptions{}
	op.Filter = ebiten.FilterLinear
	op.ColorM.Reset()
	op.GeoM.Translate(float64(winWidth/2)-float64(textRect.Dx())/2.0, float64((winHeight+winHeightHalf)/2)-float64(textRect.Dy())/2.0)
	text.DrawWithOptions(screen, msg, gameFont, op)

	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Version: %s", Version), 0, 0)
}

func renderPauseScreen(screen *ebiten.Image) {
	ebitenutil.DrawRect(screen, 0, 0, float64(winWidth), float64(winHeight), color.RGBA{0, 0, 0, 190})
	msg := "       Paused\n\n  Press Q to quit\nPress Esc to resume"
	pausedRect := text.BoundString(gameFont, msg)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(winWidth/2)-float64(pausedRect.Dx())/2.0, float64(winHeight/2)-float64(pausedRect.Dy())/2.0)
	text.DrawWithOptions(screen, msg, gameFont, op)
}

func renderHud(screen *ebiten.Image, g *Game) {
	// Update the HUD but only every 20 frames
	if g.ticks%20 == 0 {
		hudImage.Clear()

		healthStr := fmt.Sprintf("%3d - Health", g.player.health)
		healthOp := &ebiten.DrawImageOptions{}
		healthOp.GeoM.Translate(float64(hudMargin), float64(winHeight-hudMargin))
		healthOp.ColorM.Scale(0.0, 0.0, 0.0, 0.5)
		text.DrawWithOptions(hudImage, healthStr, gameFont, healthOp)
		healthOp.GeoM.Translate(-2.5, -2.5)
		healthOp.ColorM.Reset()
		healthOp.ColorM.Scale(0.889, 0.141, 0.188, 1)
		text.DrawWithOptions(hudImage, healthStr, gameFont, healthOp)

		manaStr := fmt.Sprintf("Mana - %3d", g.player.mana)
		manaRect := text.BoundString(gameFont, manaStr)
		manaOp := &ebiten.DrawImageOptions{}
		manaOp.GeoM.Translate(float64(winWidth-hudMargin-manaRect.Dx()), float64(winHeight-hudMargin))
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
	ebitenutil.DrawRect(overlayImage, px-1, py-1, 3, 3, color.RGBA{0, 255, 0, 255})

	// draw sprites
	for _, mon := range g.monsters {
		if mon == nil {
			continue
		}
		sx := mon.sprite.x / float64(cellSize/overlayCellSize)
		sy := mon.sprite.y / float64(cellSize/overlayCellSize)
		c := color.RGBA{255, 0, 0, 255}
		ebitenutil.DrawRect(overlayImage, sx-1, sy-1, 3, 3, c)
	}
	for _, item := range g.items {
		if item == nil {
			continue
		}
		sx := item.sprite.x / float64(cellSize/overlayCellSize)
		sy := item.sprite.y / float64(cellSize/overlayCellSize)
		c := color.RGBA{33, 33, 255, 255}
		ebitenutil.DrawRect(overlayImage, sx-1, sy-1, 3, 3, c)
	}

	// Draw the map
	for y := 0; y < mapSize; y++ {
		for x := 0; x < mapSize; x++ {
			if g.mapdata[x][y] != nil {
				c := color.RGBA{255, 255, 255, 58}
				ebitenutil.DrawRect(overlayImage, float64(x*overlayCellSize), float64(y*overlayCellSize), float64(overlayCellSize), float64(overlayCellSize), c)
			}
		}
	}

	w, h := overlayImage.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(winWidth)/float64(w)*overlayZoom, float64(winHeight)/float64(h)*overlayZoom)
	screen.DrawImage(overlayImage, op)
}
