package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Wall struct {
	x, y       int
	image      *ebiten.Image
	decoration *ebiten.Image
	isDoor     bool

	actionFunc func(g *Game)
}

func newWall(x, y int, kind string) *Wall {
	return &Wall{
		x:          x,
		y:          y,
		image:      imageCache["walls/"+kind],
		decoration: nil,
		actionFunc: func(g *Game) {},
	}
}

func newDoor(x, y int, kind string) *Wall {
	return &Wall{
		x:     x,
		y:     y,
		image: imageCache["doors/"+kind],
		actionFunc: func(g *Game) {
			game.mapdata[x][y] = nil
			playSound("door_open", 0.4, false)
		},
	}
}
