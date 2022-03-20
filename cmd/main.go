package main

import (
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

// Used by raycasting when rendering the view
const viewDistance = cellSize * 10
const viewRaysRatio = 1
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
	ebiten.SetWindowTitle("Castle Caster")

	log.Printf("Starting game!")
	g := &Game{}

	g.player = Player{
		x:         cellSize*1 + cellSize/2,
		y:         cellSize*1 + cellSize/2,
		angle:     0,
		moveSpeed: cellSize / 12.0,
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

	const scaleW = winWidth / 10.0
	const scaleH = winHeightHalf / 600.0
	floorOp := &ebiten.DrawImageOptions{}
	floorOp.GeoM.Scale(scaleW, scaleH)
	floorOp.GeoM.Translate(0, winHeightHalf)
	screen.DrawImage(floorImage, floorOp)
	ceilOp := &ebiten.DrawImageOptions{}
	ceilOp.GeoM.Translate(0, 0)
	ceilOp.GeoM.Scale(scaleW, scaleH)
	screen.DrawImage(ceilImage, ceilOp)

	// Cast rays to render player's view
	for i := 0; i < viewRays; i++ {
		rayAngle := g.player.angle - g.player.fov/2 + g.player.fov*float64(i)/viewRays
		t := 0.0
		for t = 0.0; t < viewDistance; t += rayStepT {
			cx := g.player.x + (t * math.Cos(rayAngle))
			cy := g.player.y + (t * math.Sin(rayAngle))

			colisionIndex := g.checkCollision(cx, cy)
			// If wall was hit
			if colisionIndex > 0 {
				hitx := int(cx) % cellSize
				hity := int(cy) % cellSize
				texColumn := hitx
				if hitx < 1 || hitx >= cellSize-1 {
					texColumn = hity
				}
				texColumn = texColumn * textureSize / cellSize

				// Get texture column at the hit point
				textureColStrip := wallImages[colisionIndex].SubImage(image.Rect(texColumn, 0, texColumn+1, textureSize)).(*ebiten.Image)

				op := &ebiten.DrawImageOptions{}
				colHeight := (winHeight / t) * colHeightScale / (math.Cos(rayAngle - g.player.angle))
				op.GeoM.Scale(viewRaysRatio, colHeight/textureSize)
				op.GeoM.Translate(float64(i)*float64(viewRaysRatio), winHeightHalf-colHeight/2)

				// Some visual flair, scale to darkness and filter the texture
				//op.Filter = ebiten.FilterLinear
				dist := 1 - (t / viewDistance * 1.5)
				dist *= 1.5

				op.ColorM.Scale(dist, dist, dist, 1)
				screen.DrawImage(textureColStrip, op)

				// Important to stop!
				break
			}
		}
	}

	for _, sprite := range g.sprites {
		drawSprite(screen, g, sprite)
	}

	// overlay map
	overlay(screen, g)
}

func drawSprite(screen *ebiten.Image, g *Game, sprite Sprite) {
	// direction to player
	spriteDir := math.Atan2(sprite.y-g.player.y, sprite.x-g.player.x)
	const spriteImgSize = 32
	const spriteImgSizeH = 16

	// remove unnecessary periods from the relative direction
	// I don't know what this really does
	for ; spriteDir-g.player.angle > math.Pi; spriteDir -= 2 * math.Pi {
	}
	for ; spriteDir-g.player.angle < -math.Pi; spriteDir += 2 * math.Pi {
	}

	spriteDist := math.Sqrt(math.Pow(g.player.x-sprite.x, 2) + math.Pow(g.player.y-sprite.y, 2))
	spriteScale := (1 / spriteDist) * winHeight
	hOffset := (spriteDir-g.player.angle)/g.player.fov*(winWidth) + (winWidth / 2) - (spriteImgSizeH * spriteScale) // do not forget the 3D view takes only a half of the framebuffer
	vOffset := winHeight/2.0 - (spriteImgSizeH / 2.0 * spriteScale)

	spriteOp := &ebiten.DrawImageOptions{}
	spriteOp.GeoM.Scale(spriteScale, spriteScale)
	spriteOp.GeoM.Translate(hOffset, vOffset)
	screen.DrawImage(spriteImages[sprite.id], spriteOp)
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
