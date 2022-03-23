package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"sort"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// 800×600, 1024×768, 1280×960, 1440×1080, 1600×1200
// Global game constants
const mapSize = 100                       // Number of grid cells, maps are assumed to be square
const winWidth = 640                      // Game window width - DON'T CHANGE
const winHeight = 480                     // Game window height - DON'T CHANGE
const winHeightHalf = winHeight / 2       // Store half the height as we use it a lot
const cellSize = 32                       // Important, how many units is each grid cell in world space - DON'T CHANGE
const textureSize = 32                    // Wall texture size (square)
const floorScaleW = winWidth / 10.0       // Used for drawing ceiling and floors
const floorScaleH = winHeightHalf / 600.0 // Used for drawing ceiling and floors
const fov = math.Pi / 3

// Used by raycasting when rendering the view
const viewDistance = cellSize * 12        // How far player can see
const viewRaysRatio = 2                   // Ratio of rays cast to screen width, higher number = less rays = faster
const rayStepT = 0.3                      // Ray casting step size, larger = less iterations = faster = inaccuracies/gaps
const viewRays = winWidth / viewRaysRatio // Number of rays to cast (see viewRaysRatio)
var magicWall = 0.0
var magicSprite = 0.0

// Used for the map overlay view
var overlayCellSize = cellSize / 4
var overlayImage = ebiten.NewImage(mapSize*overlayCellSize, mapSize*overlayCellSize)
var overlayZoom = 5.0
var overlayShown = false

// Global texture and sprite caches
var wallImages []*ebiten.Image
var spriteImages map[string]*ebiten.Image
var floorImage *ebiten.Image
var ceilImage *ebiten.Image

// Holds most core game data
type Game struct {
	mapdata   [][]int  // Map data is stored in a 2D array, 0 = empty, 1+ = wall
	mapWidth  int      // Held as convenience
	mapHeight int      // Held as convenience
	player    Player   // Player object
	level     int      // Which level we're on
	sprites   []Sprite // Any sprites on the map
}

// ===========================================================
// Load textures & sprites etc
// ===========================================================
func init() {
	var err error
	wallImages = make([]*ebiten.Image, 10)
	spriteImages = make(map[string]*ebiten.Image, 10)

	log.Printf("Loading textures...")
	for i := 1; i < 5; i++ {
		wallImages[i], _, err = ebitenutil.NewImageFromFile("./textures/" + strconv.Itoa(i) + ".png")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Loading sprites...")
	for _, spriteID := range []string{"ghoul", "skeleton", "thing"} {
		spriteImages[spriteID], _, err = ebitenutil.NewImageFromFile("./sprites/m_" + spriteID + ".png")
		if err != nil {
			log.Fatal(err)
		}
	}
	for _, spriteID := range []string{"potion", "ball"} {
		spriteImages[spriteID], _, err = ebitenutil.NewImageFromFile("./sprites/i_" + spriteID + ".png")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Loading floor & ceiling textures...")
	floorImage, _, _ = ebitenutil.NewImageFromFile("./textures/floor.png")
	ceilImage, _, _ = ebitenutil.NewImageFromFile("./textures/ceil.png")
}

// ===========================================================
// Entry point
// ===========================================================
func main() {
	winWidth = 800
	winHeight = 600
	switch winHeight {
	case 480:
		magicWall = winHeight / (cellSize / 2.6)
		magicSprite = 4
	case 600:
		magicWall = winHeight / (cellSize / 2.1)
		magicSprite = 5.2
	case 768:
		magicWall = winHeight / (cellSize / 1.666)
		magicSprite = 6.3
	case 960:
		magicWall = winHeight / (cellSize / 1.3)
		magicSprite = 8.3
	case 1080:
		magicWall = winHeight / (cellSize / 1.16)
		magicSprite = 10
	case 1200:
		magicWall = winHeight / (cellSize / 1.04)
		magicSprite = 11
	}

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Crypt Caster")
	ebiten.SetWindowResizable(true)

	log.Printf("Starting game! %f", fov)
	g := &Game{}

	g.player = Player{
		x:         cellSize*1 + cellSize/2,
		y:         cellSize*1 + cellSize/2,
		angle:     0,
		moveSpeed: cellSize / 10.0,
		turnSpeed: math.Pi / 70,
		fov:       fov,
	}

	g.level = 1
	loadMap("./maps/"+strconv.Itoa(g.level)+".txt", g)
	g.mapWidth = len(g.mapdata)
	g.mapHeight = len(g.mapdata[0])
	log.Printf("Map level %d loaded", g.level)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

// ===========================================================
// Update loop handles inputs
// ===========================================================
func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) || ebiten.IsKeyPressed(ebiten.KeyD) {
		g.player.angle += g.player.turnSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) || ebiten.IsKeyPressed(ebiten.KeyA) {
		g.player.angle -= g.player.turnSpeed
	}

	// TODO: Placeholder, replace with actual movement
	// for i := range g.sprites {
	// 	xs := g.sprites[i].x + math.Cos(g.sprites[i].angle)*g.sprites[i].speed
	// 	ys := g.sprites[i].y + math.Sin(g.sprites[i].angle)*g.sprites[i].speed
	// 	if g.checkCollision(xs, ys) > 0 {
	// 		g.sprites[i].angle += math.Pi
	// 	}
	// 	g.sprites[i].x = xs
	// 	g.sprites[i].y = ys
	// }

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyS) {
		ms := g.player.moveSpeed
		cs := cellSize / 3.0
		if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
			ms = -g.player.moveSpeed
			cs = -cs
		}

		// Check a little further ahead to see if we're going to collide with a wall
		xc := g.player.x + math.Cos(g.player.angle)*cs
		yc := g.player.y + math.Sin(g.player.angle)*cs
		if g.checkCollision(xc, yc) > 0 {
			return nil
		}

		x := g.player.x + math.Cos(g.player.angle)*ms
		y := g.player.y + math.Sin(g.player.angle)*ms
		g.player.x = x
		g.player.y = y
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
	// Render the ceiling and floor
	floorOp := &ebiten.DrawImageOptions{}
	floorOp.GeoM.Scale(floorScaleW, floorScaleH)
	floorOp.GeoM.Translate(0, winHeightHalf)
	screen.DrawImage(floorImage, floorOp)
	ceilOp := &ebiten.DrawImageOptions{}
	ceilOp.GeoM.Translate(0, 0)
	ceilOp.GeoM.Scale(floorScaleW, floorScaleH)
	screen.DrawImage(ceilImage, ceilOp)

	// Cast rays to render player's view
	for i := 0; i < viewRays; i++ {
		rayAngle := g.player.angle - g.player.fov/2 + g.player.fov*float64(i)/viewRays

		// Initialize ray and depth buffer
		t := 0.0
		depthBuffer[i] = viewDistance

		// Main ray loop
		for t = 0.0; t < viewDistance; t += rayStepT {
			// Get hit point
			cx := g.player.x + (t * math.Cos(rayAngle))
			cy := g.player.y + (t * math.Sin(rayAngle))

			// Detect collision with walls
			colisionIndex := g.checkCollision(cx, cy)
			if colisionIndex == 0 {
				continue
			}

			// If wall was hit...

			// Texture mapping
			hitx := int(cx) % cellSize
			hity := int(cy) % cellSize
			texColumn := hitx
			if hitx < 1 || hitx >= cellSize-1 {
				texColumn = hity
			}
			texColumn = texColumn * textureSize / cellSize

			// Get wall texture column at the hit point
			textureColStrip := wallImages[colisionIndex].SubImage(image.Rect(texColumn, 0, texColumn+1, textureSize)).(*ebiten.Image)

			op := &ebiten.DrawImageOptions{}

			// Scale the height of rendered wall strip to the distance and correct fish-eye effect
			// This is the heart of the 3D effect in the game
			colHeight := (winHeight / t) * magicWall / (math.Cos(rayAngle - g.player.angle))
			// Scale and place the strip
			op.GeoM.Scale(viewRaysRatio, colHeight/textureSize)
			op.GeoM.Translate(float64(i)*float64(viewRaysRatio), winHeightHalf-colHeight/2)

			// Darken as we go further
			distScale := 1 - (t / viewDistance)
			distScale = distScale * distScale

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
	// Update distances
	for i, sprite := range g.sprites {
		g.sprites[i].dist = math.Sqrt(math.Pow(g.player.x-sprite.x, 2) + math.Pow(g.player.y-sprite.y, 2))
	}
	// TODO: Can we optimize here?
	sort.Slice(g.sprites, func(i, j int) bool {
		return g.sprites[i].dist > g.sprites[j].dist
	})
	// Now render sprites
	for _, sprite := range g.sprites {
		drawSprite(screen, g, sprite)
	}

	// Overlay map
	overlay(screen, g)

	msg := fmt.Sprintf("FPS: %0.2f\n%dx%d\n\nPlayer: %f,%f", ebiten.CurrentFPS(), winWidth, winHeight, g.player.x, g.player.y)
	ebitenutil.DebugPrint(screen, msg)
}

// ===========================================================
// Required by ebiten
// ===========================================================
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return winWidth, winHeight
	//return 1024, 768
}

// ===========================================================
// Collision detection with map cells
// ===========================================================
func (g *Game) checkCollision(x, y float64) int {
	// NOTE: Bounds checking is not done, the map must have an outer wall
	mapCellX := int(x / cellSize)
	mapCellY := int(y / cellSize)
	return g.mapdata[mapCellX][mapCellY]
}

// ===========================================================
// Handle the map overlay
// ===========================================================
func overlay(screen *ebiten.Image, g *Game) {
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
		if sprite.id == "potion" {
			c = color.RGBA{0, 255, 0, 255}
		}
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
