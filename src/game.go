package main

import (
	"fmt"
	"image"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type GameState int

const (
	GameStateTitle GameState = iota
	GameStatePaused
	GameStateMain
	GameStateGameOver
	GameStateEndLevel
)

// Holds most core game data
type Game struct {
	mapdata     [][]*Wall              // Map data is stored in a 2D array, 0 = empty, 1+ = wall
	player      Player                 // Player object
	sprites     []*Sprite              // All sprites on the map, used for depth sorting
	monsters    map[uint64]*Monster    // Monsters on the map
	projectiles map[uint64]*Projectile // Projectiles currently in the game
	items       map[uint64]*Item       // Items currently in the game
	ticks       int                    // Tick count
	mapName     string
	state       GameState
	stats       Stats
}

// ===========================================================
// Update loop handles inputs
// ===========================================================
func (g *Game) start(mapName string) {
	playSound("menu_start", 2, false)
	playSoundLoop("loop_ambient_1", 1)

	log.Printf("Starting level...")
	g.sprites = make([]*Sprite, 0)
	g.monsters = make(map[uint64]*Monster, 0)
	g.projectiles = make(map[uint64]*Projectile, 0)
	g.items = make(map[uint64]*Item, 0)
	g.stats = Stats{}
	g.stats.init()
	g.stats.startTime = time.Now()

	g.player = newPlayer(1, 1)
	log.Printf("Player created %+v", g.player)

	// Precompute operations for drawing floor and ceiling
	floorOp = &ebiten.DrawImageOptions{}
	floorOp.GeoM.Scale(float64(winWidth)/10.0, float64(winHeightHalf)/600.0)
	floorOp.GeoM.Translate(0.0, float64(winHeightHalf))
	ceilOp = &ebiten.DrawImageOptions{}
	ceilOp.GeoM.Scale(float64(winWidth)/10.0, float64(winHeightHalf)/600.0)

	g.mapName = mapName
	g.loadMap(mapName)
	log.Printf("Map level '%s' loaded", g.mapName)

	g.state = GameStateMain

	// HUD image cache
	hudImage = ebiten.NewImage(winWidth, winHeight)
}

// ===========================================================
// Update loop handles inputs
// ===========================================================
func (g *Game) Update() error {
	if g.state == GameStateTitle {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
			g.start(titleLevels[titleLevelIndex])
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			os.Exit(0)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			titleLevelIndex = (titleLevelIndex + 1) % len(titleLevels)
			playSound("menu_click", 1, false)
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			titleLevelIndex--
			if titleLevelIndex < 0 {
				titleLevelIndex = len(titleLevels) - 1
			}
			playSound("menu_click", 1, false)
		}

		return nil
	}

	if g.state == GameStateGameOver || g.state == GameStateEndLevel {
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) || inpututil.IsKeyJustPressed(ebiten.KeySpace) ||
			inpututil.IsKeyJustPressed(ebiten.KeyNumpadEnter) {
			g.returnToTitleScreen()
		}
		return nil
	}

	g.ticks++

	if g.state == GameStatePaused {
		if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
			g.returnToTitleScreen()
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.state = GameStateMain
		}
		return nil
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.state = GameStatePaused
	}

	// Update rest of game state
	g.updateMonsters()
	g.updateProjectiles()

	// When move keys are first pressed, reset the acceleration timer
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyDown) ||
		inpututil.IsKeyJustPressed(ebiten.KeyW) || inpututil.IsKeyJustPressed(ebiten.KeyS) {
		g.player.moveStartTime = time.Now().UnixMicro()
	}
	// Now handle the actual move as long as move keys are held
	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyW) {
		g.player.move(time.Now().UnixMicro()-g.player.moveStartTime, -1, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		g.player.move(time.Now().UnixMicro()-g.player.moveStartTime, +1, 0)
	}

	// When turn keys are first pressed, reset the acceleration timer
	if inpututil.IsKeyJustPressed(ebiten.KeyRight) || inpututil.IsKeyJustPressed(ebiten.KeyD) || inpututil.IsKeyJustPressed(ebiten.KeyLeft) || inpututil.IsKeyJustPressed(ebiten.KeyA) {
		if ebiten.IsKeyPressed(ebiten.KeyAlt) {
			g.player.moveStartTime = time.Now().UnixMicro()
		} else {
			g.player.turnStartTime = time.Now().UnixMicro()
		}
	}
	// Now handle the actual turn as long as turn keys are held
	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		if ebiten.IsKeyPressed(ebiten.KeyAlt) {
			g.player.move(time.Now().UnixMicro()-g.player.moveStartTime, +1, -1)
		} else {
			g.player.turn(time.Now().UnixMicro()-g.player.turnStartTime, -1)
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		if ebiten.IsKeyPressed(ebiten.KeyAlt) {
			g.player.move(time.Now().UnixMicro()-g.player.moveStartTime, +1, +1)
		} else {
			g.player.turn(time.Now().UnixMicro()-g.player.turnStartTime, +1)
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
		g.player.use()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyShift) {
		g.player.attack()
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyTab) {
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
	if g.state == GameStateTitle {
		renderTitle(screen)
		return
	}

	if g.state == GameStateGameOver {
		renderGameOver(screen)
		return
	}

	if g.state == GameStateEndLevel {
		renderEndOfLevel(screen)
		return
	}

	// Render the ceiling and floor
	screen.DrawImage(imageCache["other/floor"], floorOp)
	screen.DrawImage(imageCache["other/ceil"], ceilOp)

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
			wall := g.getWallAt(cx, cy)
			if wall == nil {
				continue
			}
			if wall.invisible {
				continue
			}

			// If wall was hit...
			wall.seen = true

			// Texture mapping
			hitx := cx/cellSize - math.Floor(cx/cellSize)
			hity := cy/cellSize - math.Floor(cy/cellSize)
			texColumn := int(hitx * cellSize)
			if hitx < epi || hitx > 1-epi {
				texColumn = int(hity * cellSize)
			}

			// Get wall texture column at the hit point
			textureColStrip := wall.image.SubImage(image.Rect(texColumn, 0, texColumn+1, textureSize)).(*ebiten.Image)

			// Handle decorations
			var decoStrip *ebiten.Image
			if wall.decoration != nil {
				if wall != nil && g.ticks%20 < 10 && len(wall.metadata) > 1 && wall.metadata[1] == "torch" {
					wall.decoration = imageCache["decoration/torch-1"]
				} else if wall != nil && len(wall.metadata) > 1 && wall.metadata[1] == "torch" {
					wall.decoration = imageCache["decoration/torch"]
				}
				decoStrip = wall.decoration.SubImage(image.Rect(texColumn, 0, texColumn+1, textureSize)).(*ebiten.Image)
			}

			// Scale the height of rendered wall strip to the distance and correct fish-eye effect
			// This is the heart of the 3D effect in the game
			colHeight := (float64(winHeight) / t) * magicWall / (math.Cos(rayAngle - g.player.angle))

			// Scale and place the strip
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Scale(viewRaysRatio, colHeight/textureSize)
			op.GeoM.Translate(float64(i)*float64(viewRaysRatio), float64(winHeightHalf)-colHeight/2)

			// Darken as we go further
			distScale := 1 - (t / viewDistance)
			distScale = distScale * distScale * 1.5 // The last part brightens the textures a bit

			// Draw the strip
			op.ColorM.Scale(distScale, distScale, distScale, 1)
			screen.DrawImage(textureColStrip, op)
			if decoStrip != nil {
				screen.DrawImage(decoStrip, op)
			}

			// Save depth in buffer
			depthBuffer[i] = t

			// Important to stop!
			break
		}
	}

	// Sprite rendering loop(s)...
	// Update sprite distances
	for i := range g.sprites {
		if g.sprites[i] == nil {
			continue
		}
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

	if debug {
		msg := fmt.Sprintf("FPS: %0.2f\nPlayer: %f,%f,%f\nHolding: %+v\nLevel: %s\nVer: %s", ebiten.CurrentFPS(), g.player.x, g.player.y, g.player.angle, g.player.holding, g.mapName, Version)
		ebitenutil.DebugPrint(screen, msg)
	}

	// For screen flash effects
	if flashTimer > 0 {
		flashTimer--
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Scale(float64(winWidth), float64(winHeight))
		op.ColorM.Scale(flashColor[0], flashColor[1], flashColor[2], flashColor[3])
		screen.DrawImage(imageCache["effects/flash"], op)
	}

	renderHud(screen, g)

	if g.state == GameStatePaused {
		renderPauseScreen(screen)
	}
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
func (g *Game) getWallAt(x, y float64) *Wall {
	// NOTE: Bounds checking is not done, the map must have an outer wall
	mapCellX := int(x / cellSize)
	mapCellY := int(y / cellSize)
	if mapCellX < 0 || mapCellY < 0 || mapCellX >= mapSize || mapCellY >= mapSize {
		return nil
	}

	return g.mapdata[mapCellX][mapCellY]
}

func (g *Game) returnToTitleScreen() {
	//stopAllLoops()
	playSoundLoop("loop_menu", 0.5)
	g.state = GameStateTitle
	hudImage = nil
}

func (g *Game) gameOver() {
	playSoundLoop("loop_gameover", 0.6)
	g.state = GameStateGameOver
	hudImage = nil
}

func (g *Game) endLevel() {
	playSoundLoop("loop_end", 0.6)
	g.state = GameStateEndLevel
	hudImage = nil
	g.stats.endTime = time.Now()
}

func fireRayAt(x1, y1 float64, x2, y2 float64, maxDist float64) (wall *Wall, dist float64, angle float64) {
	newAngle := math.Atan2(y2-y1, x2-x1)

	w, d := fireRayAngle(x1, y1, newAngle, maxDist)
	return w, d, newAngle
}

func fireRayAngle(x, y float64, angle float64, maxDist float64) (w *Wall, d float64) {
	// Fire a ray in the direction we're facing
	for t := 0.0; t < maxDist; t += rayStepT {
		// Get hit point
		cx := x + (t * math.Cos(angle))
		cy := y + (t * math.Sin(angle))
		// Detect collision with walls
		if wall := game.getWallAt(cx, cy); wall != nil {
			return wall, t
		}
	}
	return nil, maxDist
}
