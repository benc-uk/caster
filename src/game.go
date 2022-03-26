package main

import (
	"fmt"
	"image"
	"math"
	"sort"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Holds most core game data
type Game struct {
	mapdata     [][]int                // Map data is stored in a 2D array, 0 = empty, 1+ = wall
	mapWidth    int                    // Held as convenience
	mapHeight   int                    // Held as convenience
	player      Player                 // Player object
	level       int                    // Which level we're on
	sprites     []*Sprite              // All sprites on the map, used for depth sorting
	monsters    map[uint64]*Monster    // Monsters on the map
	projectiles map[uint64]*Projectile // Projectiles currently in the game
	items       map[uint64]*Item       // Items currently in the game
	fc          int
}

// ===========================================================
// Update loop handles inputs
// ===========================================================
func (g *Game) Update() error {
	// Update rest of game state
	g.updateMonsters()
	g.updateProjectiles()

	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.angle += g.player.turnSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.angle -= g.player.turnSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.turnSpeed = math.Min(g.player.turnSpeed+g.player.turnSpeedAccel, g.player.turnSpeedMax)
	} else {
		g.player.turnSpeed = g.player.turnSpeedMin
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.player.use()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyShift) {
		g.player.attack()
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.move()
	} else {
		g.player.moveSpeed = g.player.moveSpeedMin
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyO) {
		overlayShown = !overlayShown
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyMinus) {
		overlayZoom -= 0.3
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEqual) {
		overlayZoom += 0.3
	}

	return nil
}

// ===========================================================
// Main draw function
// ===========================================================
func (g *Game) Draw(screen *ebiten.Image) {
	g.fc++

	// Render the ceiling and floor
	floorOp := &ebiten.DrawImageOptions{}
	floorOp.GeoM.Scale(floorScaleW, floorScaleH)
	floorOp.GeoM.Translate(0, float64(winHeightHalf))
	screen.DrawImage(floorImage, floorOp)
	ceilOp := &ebiten.DrawImageOptions{}
	ceilOp.GeoM.Translate(0, 0)
	ceilOp.GeoM.Scale(floorScaleW, floorScaleH)
	screen.DrawImage(ceilImage, ceilOp)

	// Cast rays to render player's view
	for i := 0; i < viewRays; i++ {
		rayAngle := g.player.angle - g.player.fov/2 + g.player.fov*float64(i)/float64(viewRays)

		// Initialize ray and depth buffer
		t := 0.0
		depthBuffer[i] = viewDistance

		// Main ray loop
		for t = 0.0; t < viewDistance; t += rayStepT {
			// Get hit point
			cx := g.player.x + (t * math.Cos(rayAngle))
			cy := g.player.y + (t * math.Sin(rayAngle))

			// Detect collision with walls
			wallIndex, _, _ := g.getWallAt(cx, cy)
			if wallIndex == 0 {
				continue
			}

			// If wall was hit...

			// Texture mapping
			hitx := cx/cellSize - math.Floor(cx/cellSize)
			hity := cy/cellSize - math.Floor(cy/cellSize)
			texColumn := int(hitx * cellSize)
			if hitx < epi || hitx > 1-epi {
				texColumn = int(hity * cellSize)
			}

			// Get wall texture column at the hit point
			textureColStrip := wallImages[wallIndex].SubImage(image.Rect(texColumn, 0, texColumn+1, textureSize)).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}

			// Scale the height of rendered wall strip to the distance and correct fish-eye effect
			// This is the heart of the 3D effect in the game
			colHeight := (float64(winHeight) / t) * magicWall / (math.Cos(rayAngle - g.player.angle))
			// Scale and place the strip
			op.GeoM.Scale(viewRaysRatio, colHeight/textureSize)
			op.GeoM.Translate(float64(i)*float64(viewRaysRatio), float64(winHeightHalf)-colHeight/2)

			// Darken as we go further
			distScale := 1 - (t / viewDistance)
			distScale = distScale * distScale * 1.5 // The last part brightens the textures a bit

			// Draw the strip
			op.ColorM.Scale(distScale, distScale, distScale, 1)
			screen.DrawImage(textureColStrip, op)

			// Save depth in buffer
			depthBuffer[i] = t

			// Important to stop!
			break
		}
	}

	// Sprite rendering loop(s)...
	// Update sprite distances
	for i := range g.sprites {
		g.sprites[i].dist = math.Sqrt(math.Pow(g.player.x-g.sprites[i].x, 2) + math.Pow(g.player.y-g.sprites[i].y, 2))
	}
	// TODO: Can we optimize here?
	sort.Slice(g.sprites, func(i, j int) bool {
		return g.sprites[i].dist > g.sprites[j].dist
	})

	// Now render sprites
	for _, sprite := range g.sprites {
		sprite.draw(screen, g)
	}

	// Overlay map
	g.overlay(screen)

	msg := fmt.Sprintf("FPS: %0.2f\nPlayer: %f,%f", ebiten.CurrentFPS(), g.player.x, g.player.y)
	ebitenutil.DebugPrint(screen, msg)

	renderHud(screen, g)
}

// ===========================================================
// Required by ebiten
// ===========================================================
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return winWidth, winHeight
}

// ===========================================================
// Collision detection with map cells
// ===========================================================
func (g *Game) getWallAt(x, y float64) (int, int, int) {
	// NOTE: Bounds checking is not done, the map must have an outer wall
	mapCellX := int(x / cellSize)
	mapCellY := int(y / cellSize)
	return g.mapdata[mapCellX][mapCellY], mapCellX, mapCellY
}
