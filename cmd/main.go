package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const mapSize = 10
const winWidth = 800
const winHeight = 800
const cellW = winWidth / mapSize
const cellH = winHeight / mapSize
const moveSpeed = 5.0
const viewDistance = 500
const viewRays = 800
const fov = math.Pi / 3
const rayStepT = 1

var mapColors = []color.RGBA{
	{0x00, 0x00, 0x00, 0xff},
	{95, 95, 95, 0xff},
	{190, 0, 80, 0xff},
}

type Game struct {
	mapdata [mapSize][mapSize]int
	player  Player
}

type Player struct {
	x     int
	y     int
	angle float64
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		g.player.angle += 0.06
	}

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		g.player.angle -= 0.06
	}

	if ebiten.IsKeyPressed(ebiten.KeyUp) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		ms := moveSpeed
		if ebiten.IsKeyPressed(ebiten.KeyDown) {
			ms = -moveSpeed
		}

		x := g.player.x + int(math.Cos(g.player.angle)*ms)
		y := g.player.y + int(math.Sin(g.player.angle)*ms)
		if g.checkCollision(x, y) > 0 {
			return nil
		}

		g.player.x = x
		g.player.y = y
	}

	return nil
}

func (g *Game) checkCollision(x, y int) int {
	if x < 0 || y < 0 || x > winWidth || y > winHeight {
		return 1
	}

	mapCellX := x / cellW
	mapCellY := y / cellH

	if mapCellX < 0 || mapCellX >= mapSize || mapCellY < 0 || mapCellY >= mapSize {
		return 1
	}
	return g.mapdata[mapCellX][mapCellY]
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 190, 231, 255})

	// Draw the map
	// for y := 0; y < mapSize; y++ {
	// 	for x := 0; x < mapSize; x++ {
	// 		if g.mapdata[x][y] != 0 {
	// 			ebitenutil.DrawRect(screen, float64(x*cellW), float64(y*cellH), float64(cellW), float64(cellH), mapColors[g.mapdata[x][y]])
	// 		}
	// 	}
	// }

	// Draw the player
	//ebitenutil.DrawRect(screen, float64(g.player.x-4), float64(g.player.y-4), 8, 8, color.RGBA{255, 0, 100, 255})

	// Cast rays!
	for i := 0; i < viewRays; i++ {
		rayAngle := g.player.angle - fov/2 + fov*float64(i)/viewRays
		t := 0.0
		for t = 0.0; t < viewDistance; t += rayStepT {
			x := g.player.x + int(t*math.Cos(rayAngle))
			y := g.player.y + int(t*math.Sin(rayAngle))

			colisionNum := g.checkCollision(x, y)
			if colisionNum > 0 {
				colHeight := (800 / t) * 50
				c := mapColors[colisionNum]
				d := 1 - (t / viewDistance)

				c.R = uint8(float64(c.R) * d)
				c.G = uint8(float64(c.G) * d)
				c.B = uint8(float64(c.B) * d)
				ebitenutil.DrawRect(screen, float64(i), 400-colHeight/2, 1, colHeight, c)
				ebitenutil.DrawRect(screen, float64(i), 400+colHeight/2, 1, winHeight, color.RGBA{25, 120, 10, 255})
				break
			}
		}
		if t >= viewDistance {
			ebitenutil.DrawRect(screen, float64(i), 400, 1, winHeight, color.RGBA{25, 120, 10, 255})
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Castle Caster")

	g := &Game{}
	g.mapdata = loadMap("./map.txt")

	g.player = Player{
		x:     cellW*1 + cellW/2,
		y:     cellH*1 + cellW/2,
		angle: 2.1,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}

func loadMap(filename string) [mapSize][mapSize]int {
	var mapdata [mapSize][mapSize]int
	file, err := os.Open("./map.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	y := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		for x, c := range line {
			if c == '1' {
				mapdata[x][y] = 1
			} else if c == '2' {
				mapdata[x][y] = 2
			} else {
				mapdata[x][y] = 0
			}
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return mapdata
}
