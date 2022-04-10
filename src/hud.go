package main

import (
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"os"
	"time"

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
	if hudImage == nil {
		hudImage = ebiten.NewImage(winWidth, winHeight)
		bgtiles := []string{"walls/catacombs_2", "walls/wall_vines_2", "walls/brick_brown-vines_1", "walls/snake_7", "walls/slime_6", "walls/volcanic_wall_2", "walls/cobalt_stone_9", "walls/lab-metal_1", "walls/marble_wall_5"}
		for x := 0; x < winWidth; x += cellSize * 2 {
			for y := 0; y < winHeight; y += cellSize * 2 {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(2, 2)
				op.GeoM.Translate(float64(x), float64(y))
				hudImage.DrawImage(imageCache[bgtiles[rand.Intn(len(bgtiles))]], op)
			}
		}
	}

	screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
	ebitenutil.DrawRect(screen, 0, 0, float64(winWidth), float64(winHeight), color.RGBA{0, 0, 0, 200})

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
	op.GeoM.Scale(magicSprite*0.6, magicSprite*0.6)
	op.Filter = ebiten.FilterNearest
	op.GeoM.Translate(float64(winWidth)-(35*magicSprite), float64(winHeight/3)-float64(textRect.Dy())/2.0-(15*magicSprite))
	screen.DrawImage(imageCache["hud/scroll"], op)

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
	bounds := text.BoundString(gameFont, msg)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(winWidth/2)-float64(bounds.Dx())/2.0, float64(winHeight/2)-float64(bounds.Dy())/2.0)
	text.DrawWithOptions(screen, msg, gameFont, op)
}

func renderHud(screen *ebiten.Image, g *Game) {
	// Update the HUD but only every 15 frames
	if g.ticks%hudTickInterval == 0 || forceHudUpdate {
		forceHudUpdate = false
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

		// Draw what player is holding
		items := []string{"key_red", "key_blue", "key_green"}
		for i, item := range items {
			if g.player.holding[item] <= 0 {
				continue
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(0.8*magicSprite, 0.8*magicSprite)
			op.GeoM.Translate(float64(winWidth)-28*magicSprite, (float64(i) * 8 * magicSprite))
			hudImage.DrawImage(imageCache["items/"+item], op)

			if g.player.holding[item] > 1 {
				textOp := &ebiten.DrawImageOptions{}
				textOp.GeoM.Translate(float64(winWidth)-6*magicSprite, (24*magicSprite)+(float64(i)*8*magicSprite))
				text.DrawWithOptions(hudImage, fmt.Sprintf("%d", g.player.holding[item]), gameFont, textOp)
			}
		}

		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(1.5*magicSprite, 1.5*magicSprite)
		weaponOffset := 96.0
		if g.player.justFired {
			weaponOffset = 90
			g.player.justFired = false
			op.ColorM.Scale(1, 2, 1, 1)
		}
		op.GeoM.Translate((float64(winWidth)/2.0)-(48*magicSprite), float64(winHeight)-(weaponOffset*magicSprite))
		hudImage.DrawImage(imageCache["hud/weapon_0"], op)

		screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
	} else {
		screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
	}
}

func renderGameOver(screen *ebiten.Image) {
	if hudImage == nil {
		hudImage = ebiten.NewImageFromImage(screen)
		for i := 0; i < 500; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(rand.Float64()*float64(winWidth)-100, rand.Float64()*float64(winHeight)-100)
			hudImage.DrawImage(imageCache["hud/skull"], op)
		}
		ebitenutil.DrawRect(hudImage, 0, 0, float64(winWidth), float64(winHeight), color.RGBA{0, 0, 0, 190})
		msg := "     You Have Died\n This Is Unfortunate\n\nPress Enter to restart"
		bounds := text.BoundString(gameFont, msg)
		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(0.9, 0, 0, 1)
		op.GeoM.Translate(float64(winWidth/2)-float64(bounds.Dx())/2.0, float64(winHeight/2)-float64(bounds.Dy())/2.0)
		text.DrawWithOptions(hudImage, msg, gameFont, op)
	}
	screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
}

func renderEndOfLevel(screen *ebiten.Image) {
	if hudImage == nil {

		itemPercentage := (float64(game.stats.itemsFound) / float64(game.stats.itemsTotal)) * 100.0
		monsterPercentage := (float64(game.stats.kills) / float64(game.stats.monsters)) * 100.0
		secretPercentage := 100.0
		if game.stats.secretsTotal > 0 {
			secretPercentage = (float64(game.stats.secretsFound) / float64(game.stats.secretsTotal)) * 100.0
		}
		timeTaken := game.stats.endTime.Sub(game.stats.startTime)
		timeTaken = timeTaken.Round(time.Second)

		specialMsg := ""
		if secretPercentage >= 100.0 && monsterPercentage >= 100.0 && itemPercentage >= 100.0 {
			specialMsg = "\n\nWOW! PERFECT JOB!!"
		}

		hudImage = ebiten.NewImageFromImage(screen)
		for i := 0; i < 500; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(5, 5)
			op.GeoM.Translate(rand.Float64()*float64(winWidth)-100, rand.Float64()*float64(winHeight)-100)
			if specialMsg != "" {
				hudImage.DrawImage(imageCache["hud/rainbow"], op)
			} else {
				hudImage.DrawImage(imageCache["hud/cloud"], op)
			}
		}
		ebitenutil.DrawRect(hudImage, 0, 0, float64(winWidth), float64(winHeight), color.RGBA{0, 0, 0, 190})

		msg := fmt.Sprintf("    You Escaped!\n\nMonsters Killed: %.1f %%\nItems Found: %.1f %%\nSecrets Found: %.1f %%\nTime Taken: %s%s\n\nPress Enter To Restart", monsterPercentage, itemPercentage, secretPercentage, timeTaken, specialMsg)
		bounds := text.BoundString(gameFont, msg)
		op := &ebiten.DrawImageOptions{}
		op.ColorM.Scale(0.1, 0.8, 0.2, 1)
		op.GeoM.Translate(float64(winWidth/2)-float64(bounds.Dx())/2.0, float64(winHeight/2)-float64(bounds.Dy())/2.0)
		text.DrawWithOptions(hudImage, msg, gameFont, op)
	}

	screen.DrawImage(hudImage, &ebiten.DrawImageOptions{})
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
		if !mon.seenPlayer {
			continue
		}
		sx := mon.sprite.x / float64(cellSize/overlayCellSize)
		sy := mon.sprite.y / float64(cellSize/overlayCellSize)
		c := color.RGBA{255, 0, 0, 255}
		ebitenutil.DrawRect(overlayImage, sx-1, sy-1, 3, 3, c)
	}

	// for _, item := range g.items {
	// 	if item == nil {
	// 		continue
	// 	}
	// 	// if !item.sprite.seen {
	// 	// 	continue
	// 	// }
	// 	sx := item.sprite.x / float64(cellSize/overlayCellSize)
	// 	sy := item.sprite.y / float64(cellSize/overlayCellSize)
	// 	c := color.RGBA{33, 33, 255, 255}
	// 	ebitenutil.DrawRect(overlayImage, sx-1, sy-1, 3, 3, c)
	// }

	// Draw the map
	for y := 0; y < mapSize; y++ {
		for x := 0; x < mapSize; x++ {
			if g.mapdata[x][y] != nil {
				if !g.mapdata[x][y].seen {
					continue
				}
				c := color.RGBA{255, 255, 255, 58}
				if g.mapdata[x][y].isDoor {
					c = color.RGBA{110, 50, 15, 70}
				}
				ebitenutil.DrawRect(overlayImage, float64(x*overlayCellSize), float64(y*overlayCellSize), float64(overlayCellSize), float64(overlayCellSize), c)
			}
		}
	}

	w, h := overlayImage.Size()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(winWidth)/float64(w)*overlayZoom, float64(winHeight)/float64(h)*overlayZoom)
	screen.DrawImage(overlayImage, op)
}

// ===========================================================
// Briefly flash the screen white, until the next HUD update
// ===========================================================
func screenFlashWhite(time int) {
	flashColor = []float64{1.0, 1.0, 1.0, 0.8}
	flashTimer = time
}

func screenFlashRed(time int) {
	flashColor = []float64{1.5, 0, 0, 0.8}
	flashTimer = time
}
