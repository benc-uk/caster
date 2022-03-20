package main

import (
	"bufio"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math"
	"os"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const mapSize = 100
const winWidth = 800
const winHeight = 600
const winHeightHalf = winHeight / 2
const cellSize = 40
const textureSize = 32

const viewDistance = cellSize * 8
const viewRaysRatio = 2

const colHeightScale = winHeight / 15
const viewRays = winWidth / viewRaysRatio
const fov = math.Pi / 3
const rayStepT = 0.2

var skyColour = color.RGBA{0, 190, 231, 255}
var floorColour = color.RGBA{64, 50, 27, 255}

type Game struct {
	mapdata   [][]int
	mapWidth  int
	mapHeight int
	player    Player
}

type Player struct {
	x         float64
	y         float64
	angle     float64
	moveSpeed float64
	turnSpeed float64
}

var wallImages []*ebiten.Image

func init() {
	var err error
	wallImages = make([]*ebiten.Image, 10)

	wallImages[1], _, err = ebitenutil.NewImageFromFile("./textures/1.png")
	wallImages[2], _, err = ebitenutil.NewImageFromFile("./textures/2.png")
	wallImages[3], _, err = ebitenutil.NewImageFromFile("./textures/3.png")
	wallImages[4], _, err = ebitenutil.NewImageFromFile("./textures/4.png")
	if err != nil {
		log.Fatal(err)
	}
}

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

	return nil
}

func (g *Game) checkCollision(x, y float64) int {
	// if x < 0 || y < 0 || x > winWidth || y > winHeight {
	// 	return 1
	// }

	mapCellX := int(x / cellSize)
	mapCellY := int(y / cellSize)

	// if mapCellX < 0 || mapCellX >= g.mapWidth || mapCellY < 0 || mapCellY >= g.mapHeight {
	// 	return 1
	// }
	return g.mapdata[mapCellX][mapCellY]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(skyColour)

	// Cast rays to render player's view
	for i := 0; i < viewRays; i++ {
		rayAngle := g.player.angle - fov/2 + fov*float64(i)/viewRays
		t := 0.0
		for t = 0.0; t < viewDistance; t += rayStepT {
			cx := g.player.x + (t * math.Cos(rayAngle))
			cy := g.player.y + (t * math.Sin(rayAngle))

			colisionIndex := g.checkCollision(cx, cy)
			if colisionIndex > 0 {
				colHeight := (winHeight / t) * colHeightScale / (math.Cos(rayAngle - g.player.angle))

				hitx := int(cx) % cellSize
				hity := int(cy) % cellSize
				texcoord := hitx
				if hitx == 0 || hitx == cellSize-1 {
					texcoord = hity
				}
				texcoord = texcoord * textureSize / cellSize

				// Get texture column at
				colStrip := textureColumn(wallImages[colisionIndex], texcoord)

				op := &ebiten.DrawImageOptions{}
				op.GeoM.Scale(viewRaysRatio, colHeight/textureSize)
				op.GeoM.Translate(float64(i*viewRaysRatio), winHeightHalf-colHeight/2)
				screen.DrawImage(colStrip, op)

				// draw floor
				ebitenutil.DrawRect(screen, float64(i*viewRaysRatio), winHeightHalf+colHeight/2, viewRaysRatio, winHeightHalf, floorColour)
				break
			}
		}

		if t >= viewDistance {
			ebitenutil.DrawRect(screen, float64(i*viewRaysRatio), winHeightHalf, viewRaysRatio, winHeight, floorColour)
		}
	}

	// Debug
	//overlay(screen, g)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Castle Caster")

	g := &Game{}
	g.mapdata = loadMap("./map.txt")
	g.mapWidth = len(g.mapdata)
	g.mapHeight = len(g.mapdata[0])
	log.Printf("Map size: %dx%d", g.mapWidth, g.mapHeight)

	g.player = Player{
		x:         cellSize*1 + cellSize/2,
		y:         cellSize*1 + cellSize/2,
		angle:     0,
		moveSpeed: cellSize / 12.0,
		turnSpeed: math.Pi / 70,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func textureColumn(img *ebiten.Image, x int) *ebiten.Image {
	_, h := img.Size()
	i := img.SubImage(image.Rect(x, 0, x+1, h))
	return i.(*ebiten.Image)
}

func loadMap(filename string) [][]int {
	var mapdata = make([][]int, mapSize)
	for i := range mapdata {
		mapdata[i] = make([]int, mapSize)
	}

	file, err := os.Open("./map.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	y := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)

		for x, c := range line {
			i, err := strconv.Atoi(string(c))
			if err != nil {
				mapdata[x][y] = 0
			} else {
				mapdata[x][y] = i
			}
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return mapdata
}

func overlay(screen *ebiten.Image, g *Game) {
	// Draw the player
	ebitenutil.DrawRect(screen, float64(g.player.x-4), float64(g.player.y-4), 8, 8, color.RGBA{255, 255, 255, 255})

	// Draw the map
	for y := 0; y < g.mapHeight; y++ {
		for x := 0; x < g.mapWidth; x++ {
			if g.mapdata[x][y] != 0 {
				ebitenutil.DrawRect(screen, float64(x*cellSize), float64(y*cellSize), float64(cellSize), float64(cellSize), color.RGBA{255, 255, 255, 58})
			}
		}
	}
}
