package main

import (
	"flag"
	_ "image/png"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

// 800×600, 1024×768, 1280×960, 1440×1080, 1600×1200
// Global game constants
const mapSize = 100    // Number of grid cells, maps are assumed to be square
const cellSize = 32    // Important, how many units is each grid cell in world space - DON'T CHANGE
const textureSize = 32 // Wall texture size (square)

// Used by raycasting when rendering the view
const fov = math.Pi / 3.0
const viewDistance = cellSize * 12 // How far player can see
var viewRaysRatio = 4.0            // Ratio of rays cast to screen width, higher number = less rays = faster
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
var hudMargin = 0

// Used for the map overlay view
var overlayCellSize = cellSize / 4
var overlayImage = ebiten.NewImage(mapSize*overlayCellSize, mapSize*overlayCellSize)
var overlayZoom = 5.0
var overlayShown = false

// Global texture and sprite caches for walls and floors
var wallImages []*ebiten.Image
var floorImage *ebiten.Image
var ceilImage *ebiten.Image

var game *Game
var gameFont font.Face

// ===========================================================
// Load textures & sprites etc
// ===========================================================
func init() {
	var err error
	wallImages = make([]*ebiten.Image, 25)
	spriteImages = make(map[string]*ebiten.Image, 10)

	log.Printf("Loading wall textures...")
	for i := 1; i < 25; i++ {
		img, _, err := ebitenutil.NewImageFromFile("./textures/" + strconv.Itoa(i) + ".png")
		if err == nil {

			wallImages[i] = img
		}
	}

	initSprites()

	log.Printf("Loading floor & ceiling textures...")
	floorImage, _, _ = ebitenutil.NewImageFromFile("./textures/floor.png")
	ceilImage, _, _ = ebitenutil.NewImageFromFile("./textures/ceil.png")

	initSound()

	// Font(s
	fontData, err := os.ReadFile("./fonts/morris-roman.ttf")
	if err != nil {
		log.Fatal(err)
	}
	tt, err := opentype.Parse(fontData)
	if err != nil {
		log.Fatal(err)
	}
	gameFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    100,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// ===========================================================
// Entry point
// ===========================================================
func main() {

	var flagRes string
	var flagRatio int
	var flagFull bool
	var flagVsync bool
	flag.StringVar(&flagRes, "res", "medium", "Screen resolution: tiny, small, medium, large, larger or super")
	flag.IntVar(&flagRatio, "ratio", 4, "Ray rendering ratio as a percentage of screen width")
	flag.BoolVar(&flagFull, "fullscreen", false, "Fullscreen mode (default false)")
	flag.BoolVar(&flagVsync, "vsync", false, "Enable vsync (default false)")
	flag.Parse()

	if flagRatio > 0 {
		viewRaysRatio = float64(flagRatio)
	}

	switch flagRes {
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
	viewRays = winWidth / int(viewRaysRatio)
	depthBuffer = make([]float64, viewRays)
	hudMargin = winHeight / 35

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Crypt Caster")
	ebiten.SetWindowResizable(true)
	if flagVsync {
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	} else {
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	}
	ebiten.SetFullscreen(flagFull)

	log.Printf("Resolution: %dx%d, Ray ratio: %f", winWidth, winHeight, viewRaysRatio)
	log.Printf("Starting game!")
	game = &Game{
		sprites:     make([]*Sprite, 0),
		monsters:    make(map[uint64]*Monster, 0),
		projectiles: make(map[uint64]*Projectile, 0),
		items:       make(map[uint64]*Item, 0),
	}

	game.player = newPlayer(1, 1)
	log.Printf("Player created %+v", game.player)

	game.level = 1
	game.loadMap("./maps/" + strconv.Itoa(game.level) + ".txt")
	game.mapWidth = len(game.mapdata)
	game.mapHeight = len(game.mapdata[0])
	log.Printf("Map level %d loaded", game.level)

	// HUD image cache
	hudImage = ebiten.NewImage(winWidth, winHeight)

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
