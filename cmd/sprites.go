package main

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const spriteImgSize = 32
const spriteImgSizeH = 16

// Used for rendering sprites with occlusion
var depthBuffer = make([]float64, viewRays)

type Sprite struct {
	x     float64
	y     float64
	id    string
	dist  float64
	scale float64
	dir   float64
}

// ===========================================================
// Draws a sprite on the screen with correct depth
// ===========================================================
func drawSprite(screen *ebiten.Image, g *Game, sprite Sprite) {
	if sprite.dist > viewDistance {
		return
	}

	// Sizing and scaling based on depth
	spriteDist := (1.0 / sprite.dist)
	spriteScale := spriteDist * winHeight
	// Direction to player
	spriteDir := math.Atan2(sprite.y-g.player.y, sprite.x-g.player.x)

	// I don't know what this really does, actually no fucking clue
	for ; spriteDir-g.player.angle > math.Pi; spriteDir -= 2 * math.Pi {
	}
	for ; spriteDir-g.player.angle < -math.Pi; spriteDir += 2 * math.Pi {
	}

	// The X coordinate of the sprite
	hOffset := (spriteDir-g.player.angle)/g.player.fov*(winWidth) + (winWidth / 2) - (spriteImgSizeH * spriteScale)

	// TODO: Remove this? - Crude culling
	centerX := hOffset + (spriteImgSizeH * spriteScale)
	if centerX < 0 || centerX > winWidth {
		return
	}

	// The Y coordinate of the sprite
	// HACK: This only works when the screen height is 1024, I've lost DAYS trying to fix it
	vOffset := winHeight/2.0 - (spriteImgSize*spriteDist*colHeightScale)*5.3

	// To position the sprite
	spriteOp := &ebiten.DrawImageOptions{}
	spriteOp.GeoM.Scale(spriteScale, spriteScale)
	spriteOp.GeoM.Translate(hOffset, vOffset)

	// Slice the sprite image into strips and render each one
	spriteImg := spriteImages[sprite.id]
	for slice := 0; slice < spriteImgSize; slice++ {
		// Each loop move the slice along with scaling taken into account
		spriteOp.GeoM.Translate(spriteScale, 0)

		// Check the depth buffer, and skip if the sprite is a wall
		depthBufferX := math.Floor(spriteOp.GeoM.Element(0, 2)) / viewRaysRatio
		if depthBufferX < 0 || depthBufferX >= viewRays || depthBuffer[int(depthBufferX)] < sprite.dist {
			continue
		}

		// Draw the sprite slice
		sliceImg := spriteImg.SubImage(image.Rect(slice, 0, slice+1, spriteImgSize)).(*ebiten.Image)
		screen.DrawImage(sliceImg, spriteOp)
	}
}
