package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// Global game constants
const mapSize = 100
const winWidth = 1024
const winHeight = 768
const winHeightHalf = winHeight / 2
const cellSize = 36
const textureSize = 32
const spriteImgSize = 32
const spriteImgSizeH = 16
const floorScaleW = winWidth / 10.0
const floorScaleH = winHeightHalf / 600.0

// Used by raycasting when rendering the view
const viewDistance = cellSize * 10
const viewRaysRatio = 4 // I think changing this now just causes problems
const colHeightScale = winHeight / (cellSize / 2)
const viewRays = winWidth / viewRaysRatio
const rayStepT = 0.3

var overlayCellSize = cellSize / 4
var overlayImage = ebiten.NewImage(mapSize*overlayCellSize, mapSize*overlayCellSize)
var overlayZoom = 1.0
var overlayShown = false

var wallImages []*ebiten.Image
var spriteImages map[string]*ebiten.Image
var floorImage *ebiten.Image
var ceilImage *ebiten.Image

var depthBuffer = make([]float64, viewRays)

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
	for _, spriteId := range []string{"ghoul", "skeleton", "thing"} {
		spriteImages[spriteId], _, err = ebitenutil.NewImageFromFile("./monsters/" + spriteId + ".png")
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Printf("Loading floor & ceiling texture...")
	floorImage, _, _ = ebitenutil.NewImageFromFile("./textures/floor.png")
	ceilImage, _, _ = ebitenutil.NewImageFromFile("./textures/ceil.png")
}

// ===========================================================
// Entry point
// ===========================================================
func main() {
	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Crypt Caster")

	log.Printf("Starting game!")
	g := &Game{}

	g.player = Player{
		x:         cellSize*1 + cellSize/2,
		y:         cellSize*1 + cellSize/2,
		angle:     0,
		moveSpeed: cellSize / 10.0,
		turnSpeed: math.Pi / 70,
		fov:       math.Pi / 3,
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
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.angle += g.player.turnSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.angle -= g.player.turnSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		ms := g.player.moveSpeed
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			ms = -g.player.moveSpeed
		}

		x := g.player.x + math.Cos(g.player.angle)*ms
		y := g.player.y + math.Sin(g.player.angle)*ms
		if g.checkCollision(x, y) > 0 {
			return nil
		}

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
			colHeight := (winHeight / t) * colHeightScale / (math.Cos(rayAngle - g.player.angle))
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

	// Sprite rendering loop
	for _, sprite := range g.sprites {
		drawSprite(screen, g, sprite)
	}

	// overlay map
	overlay(screen, g)

	msg := fmt.Sprintf("FPS: %0.2f\n", ebiten.CurrentFPS())
	ebitenutil.DebugPrint(screen, msg)

}

// ===========================================================
// Required by ebiten
// ===========================================================
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
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
