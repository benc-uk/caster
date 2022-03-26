package main

import (
	_ "image/png"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

// 800×600, 1024×768, 1280×960, 1440×1080, 1600×1200
// Global game constants
const mapSize = 100    // Number of grid cells, maps are assumed to be square
const cellSize = 32    // Important, how many units is each grid cell in world space - DON'T CHANGE
const textureSize = 32 // Wall texture size (square)

// Used by raycasting when rendering the view
const fov = math.Pi / 3.0
const viewDistance = cellSize * 12 // How far player can see
const viewRaysRatio = 2            // Ratio of rays cast to screen width, higher number = less rays = faster
const rayStepT = 0.3               // Ray casting step size, larger = less iterations = faster = inaccuracies/gaps
const epi = 0.01

// These are all constant but dependent on the window size, which can be changed
var winWidth = 0      // Game window width
var winHeight = 0     // Game window height
var winHeightHalf = 0 // Store half the height as we use it a lot
var floorScaleW = 0.0 // Used for drawing ceiling and floors
var floorScaleH = 0.0 // Used for drawing ceiling and floors
var viewRays = 0      // Number of rays to cast (see viewRaysRatio)
var magicWall = 0.0   // Used to scale height of walls
var magicSprite = 0.0 // Used to scale position of sprites

// Used for the map overlay view
var overlayCellSize = cellSize / 4
var overlayImage = ebiten.NewImage(mapSize*overlayCellSize, mapSize*overlayCellSize)
var overlayZoom = 5.0
var overlayShown = false

// Global texture and sprite caches for walls and floors
var wallImages []*ebiten.Image
var floorImage *ebiten.Image
var ceilImage *ebiten.Image

// ===========================================================
// Load textures & sprites etc
// ===========================================================
func init() {
	var err error
	wallImages = make([]*ebiten.Image, 10)
	spriteImages = make(map[string]*ebiten.Image, 10)

	log.Printf("Loading wall textures...")
	for i := 1; i < 10; i++ {
		wallImages[i], _, err = ebitenutil.NewImageFromFile("./textures/" + strconv.Itoa(i) + ".png")
		if err != nil {
			log.Fatal(err)
		}
	}

	initSprites()

	log.Printf("Loading floor & ceiling textures...")
	floorImage, _, _ = ebitenutil.NewImageFromFile("./textures/floor.png")
	ceilImage, _, _ = ebitenutil.NewImageFromFile("./textures/ceil.png")

	initSound()
}

// ===========================================================
// Entry point
// ===========================================================
func main() {
	res := "medium"

	if len(os.Args) > 1 {
		res = os.Args[1]
	}

	switch res {
	case "tiny":
		winWidth = 640
		winHeight = 480
		magicWall = float64(winHeight) / (float64(cellSize) / 2.6)
		magicSprite = 4
	case "small":
		winWidth = 800
		winHeight = 600
		magicWall = float64(winHeight) / (float64(cellSize) / 2.1)
		magicSprite = 5.2
	case "medium":
		winWidth = 1024
		winHeight = 768
		magicWall = float64(winHeight) / (float64(cellSize) / 1.666)
		magicSprite = 6.3
	case "large":
		winWidth = 1280
		winHeight = 960
		magicWall = float64(winHeight) / (float64(cellSize) / 1.3)
		magicSprite = 8.3
	case "larger":
		winWidth = 1440
		winHeight = 1080
		magicWall = float64(winHeight) / (float64(cellSize) / 1.16)
		magicSprite = 10
	case "super":
		winWidth = 1600
		winHeight = 1200
		magicWall = float64(winHeight) / (float64(cellSize) / 1.04)
		magicSprite = 11
	default:
		log.Fatalln("Invalid resolution provided, use: tiny, small, medium, large, larger or super")
	}

	// Set all those magic numbers
	winHeightHalf = winHeight / 2
	floorScaleW = float64(winWidth) / 10.0
	floorScaleH = float64(winHeightHalf) / 600.0
	viewRays = winWidth / viewRaysRatio
	depthBuffer = make([]float64, viewRays)

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Crypt Caster")
	ebiten.SetWindowResizable(true)

	log.Printf("Starting game!")
	g := &Game{
		sprites:     make([]*Sprite, 0),
		monsters:    make(map[uint64]*Monster, 0),
		projectiles: make(map[uint64]*Projectile, 0),
	}

	g.player = newPlayer(g)
	log.Printf("Player created %+v", g.player)

	g.level = 1
	loadMap("./maps/"+strconv.Itoa(g.level)+".txt", g)
	g.mapWidth = len(g.mapdata)
	g.mapHeight = len(g.mapdata[0])
	log.Printf("Map level %d loaded", g.level)

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
