package main

import (
	"image"
	"log"
	"math"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const spriteImgSize = 32
const spriteImgSizeH = 16

// Used for rendering sprites with occlusion
var depthBuffer []float64

// Global list of ALL sprites, used for depth sorting
var spriteImages map[string]*ebiten.Image

type Sprite struct {
	x     float64
	y     float64
	kind  string
	dist  float64 // distance to player, updated during the render cycle
	angle float64
	speed float64
	size  float64
	image *ebiten.Image
	alpha float64
}

func (g *Game) addSprite(kind string, x, y float64, angle float64, speed float64, size float64) *Sprite {
	if spriteImages[kind] == nil {
		log.Printf("ERROR! Sprite image not found: %s", kind)
		return nil
	}

	s := &Sprite{
		x:     x,
		y:     y,
		kind:  kind,
		angle: angle,
		speed: speed,
		size:  size,
		image: spriteImages[kind],
		alpha: 1.0,
	}

	g.sprites = append(g.sprites, s)
	return s
}

func (s *Sprite) isHit(x, y float64) bool {
	deltaX := x - s.x
	deltaY := y - s.y
	dist := math.Sqrt(deltaX*deltaX + deltaY*deltaY)
	return dist < s.size
}

// TODO: this is pretty inefficient, but keeping the sprites in a map too is hard for sorting
func (g *Game) removeSprite(sprite *Sprite) {
	for i, s := range g.sprites {
		if s == sprite {
			g.sprites = append(g.sprites[:i], g.sprites[i+1:]...)
			break
		}
	}
}

// ===========================================================
// Load sprites at startup
// ===========================================================
func initSprites() {
	var err error

	log.Printf("Loading sprites...")
	spriteDir, err := os.ReadDir("./sprites")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range spriteDir {
		name := strings.TrimSuffix(file.Name(), ".png")
		spriteImages[name], _, err = ebitenutil.NewImageFromFile("./sprites/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
	}
}

//
// ===========================================================
// Draws a sprite on the screen with correct depth
// ===========================================================
func (s *Sprite) draw(screen *ebiten.Image, g *Game) {
	if s.dist > viewDistance {
		return
	}

	// Sizing and scaling based on depth
	spriteDist := (1.0 / s.dist)
	spriteScale := spriteDist * float64(winHeight)
	// Direction to player
	spriteDir := math.Atan2(s.y-g.player.y, s.x-g.player.x)

	// I don't know what this really does, actually no fucking clue
	for ; spriteDir-g.player.angle > math.Pi; spriteDir -= 2 * math.Pi {
	}
	for ; spriteDir-g.player.angle < -math.Pi; spriteDir += 2 * math.Pi {
	}

	winWidthF := float64(winWidth)
	winHeightHalfF := float64(winHeightHalf)
	// The X coordinate of the sprite
	hOffset := (spriteDir-g.player.angle)/g.player.fov*(winWidthF) + (winWidthF / 2) - (spriteImgSizeH * spriteScale)

	// TODO: Remove this? - Crude culling
	centerX := hOffset + (spriteImgSizeH * spriteScale)
	if centerX < 0 || centerX > winWidthF {
		return
	}

	// The Y coordinate of the sprite
	// HACK: HERE BE EVIL! 😈 This only works when the screen height is 1024, I've lost DAYS trying to fix it
	vOffset := winHeightHalfF - spriteImgSize*magicWall*spriteDist*magicSprite

	// To position the sprite
	spriteOp := &ebiten.DrawImageOptions{}
	spriteOp.GeoM.Scale(spriteScale, spriteScale)
	spriteOp.GeoM.Translate(hOffset, vOffset)
	spriteOp.ColorM.Scale(1, 1, 1, s.alpha)

	// Slice the sprite image into strips and render each one
	spriteImg := spriteImages[s.kind]
	for slice := 0; slice < spriteImgSize; slice++ {
		// Each loop move the slice along with scaling taken into account
		spriteOp.GeoM.Translate(spriteScale, 0)

		// Check the depth buffer, and skip if the sprite is a wall
		depthBufferX := int(math.Floor(spriteOp.GeoM.Element(0, 2)) / viewRaysRatio)
		if depthBufferX < 0 || depthBufferX >= viewRays || depthBuffer[int(depthBufferX)] < s.dist {
			continue
		}

		// Draw the sprite slice
		sliceImg := spriteImg.SubImage(image.Rect(slice, 0, slice+1, spriteImgSize)).(*ebiten.Image)
		screen.DrawImage(sliceImg, spriteOp)
	}
}
