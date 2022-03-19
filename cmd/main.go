package main

import (
	"image/color"
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const mapSize = 6
const winWidth = 800
const winHeight = 800

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
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{30, 150, 20, 255})

	rectW := winWidth / mapSize
	rectH := winHeight / mapSize

	// Draw the map
	for y := 0; y < mapSize; y++ {
		for x := 0; x < mapSize; x++ {
			if g.mapdata[x][y] != 0 {
				ebitenutil.DrawRect(screen, float64(x*rectW), float64(y*rectH), float64(rectW), float64(rectH), mapColors[g.mapdata[x][y]])
			}
		}
	}

	// Draw the player
	ebitenutil.DrawRect(screen, float64(g.player.x-4), float64(g.player.y-4), 8, 8, color.RGBA{0, 0, 0, 255})

	// Draw the player's view
	for t := 0.0; t < 500; t += 2.0 {
		x := g.player.x + int(t*math.Cos(g.player.angle))
		y := g.player.y + int(t*math.Sin(g.player.angle))

		mapCellX := x / rectW
		mapCellY := y / rectH
		if mapCellX < 0 || mapCellX >= mapSize || mapCellY < 0 || mapCellY >= mapSize {
			break
		}
		if g.mapdata[mapCellX][mapCellY] != 0 {
			break
		}

		ebitenutil.DrawRect(screen, float64(x), float64(y), 2, 2, color.RGBA{255, 255, 255, 255})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}

func main() {
	ebiten.SetWindowSize(winWidth, winHeight)
	ebiten.SetWindowTitle("Castle Caster")

	g := &Game{}
	g.mapdata[1][1] = 1
	g.mapdata[2][1] = 1
	g.mapdata[3][1] = 2
	g.mapdata[4][1] = 1
	g.mapdata[4][2] = 1
	g.mapdata[4][4] = 2

	g.player = Player{
		x:     40,
		y:     30,
		angle: 2.1,
	}

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
