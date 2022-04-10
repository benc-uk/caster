package main

import (
	"flag"
	_ "image/png"
	"log"
	"math"
	"math/rand"
	"path/filepath"
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

var game *Game
var Version = "40"
var debug = false
var titleLevelIndex = 0
var titleLevels = []string{}

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
var viewRays = 0      // Number of rays to cast (see viewRaysRatio)
var magicWall = 0.0   // Used to scale height of walls
var magicSprite = 0.0 // Used to scale position of sprites
var hudMargin = 0

// Used for the map overlay view
var overlayCellSize = cellSize / 2
var overlayImage = ebiten.NewImage(mapSize*overlayCellSize, mapSize*overlayCellSize)
var overlayZoom = 5.0
var overlayShown = false

// Global/precomputed vars for drawing floor and ceiling
var floorOp *ebiten.DrawImageOptions
var ceilOp *ebiten.DrawImageOptions

var flashTimer = 0
var flashColor = []float64{1, 1, 1, 0.8}
var forceHudUpdate = false

// ===========================================================
// Load textures & sprites etc
// ===========================================================
func init() {
	// Load all textures and sprites
	loadImageCache()

	// Find all maps in the maps folder
	maps, err := filepath.Glob("maps/*.json")
	if err != nil {
		log.Fatal(err)
	}
	for _, mapFile := range maps {
		titleLevels = append(titleLevels, strings.TrimSuffix(filepath.Base(mapFile), ".json"))
	}

	// Load all sounds
	initSound()
}

// ===========================================================
// Entry point
// ===========================================================
func main() {
	rand.Seed(time.Now().UnixNano())

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
	viewRays = winWidth / int(viewRaysRatio)
	depthBuffer = make([]float64, viewRays)
	hudMargin = winHeight / 35

	// Call this after the magic numbers are set
	initHUD()

	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Crypt Caster")
	ebiten.SetWindowResizable(true)
	if flagVsync {
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
	} else {
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOffMaximum)
	}
	ebiten.SetFullscreen(flagFull)

	log.Printf("Starting game...")
	log.Printf("Resolution: %dx%d, Ray ratio: %f", winWidth, winHeight, viewRaysRatio)

	game = &Game{
		state: GameStateTitle,
	}

	game.gameOver()

	// HACK: ONLY FOR DEBUGGING/TESTING
	game.start("Entryway")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
