package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Wall struct {
	x, y       int
	image      *ebiten.Image
	decoration *ebiten.Image
	metadata   []string

	actionFunc func(g *Game)
}

func newWall(x, y int, kind string) *Wall {
	return &Wall{
		x:          x,
		y:          y,
		image:      imageCache["walls/"+kind],
		decoration: nil,
		actionFunc: func(g *Game) {
			playSound("grunt", 1.0, false)
		},
	}
}

func newDoor(x, y int, kind string) *Wall {
	w := &Wall{
		x:     x,
		y:     y,
		image: imageCache["doors/"+kind],
		actionFunc: func(g *Game) {
			playSound("locked", 1.0, false)
		},
	}

	if kind == "basic" {
		w.actionFunc = func(g *Game) {
			game.mapdata[x][y] = nil
			playSound("door_open", 0.4, false)
		}
	}

	return w
}

func newSecretWall(x, y int, kind string) *Wall {
	return &Wall{
		x:     x,
		y:     y,
		image: imageCache["walls/"+kind],
		actionFunc: func(g *Game) {
			game.mapdata[x][y] = nil
			playSound("secret", 1.0, false)
		},
	}
}

func newSwitchWall(x, y int, kind string, tx, ty int) *Wall {
	return &Wall{
		x:          x,
		y:          y,
		image:      imageCache["walls/"+kind],
		decoration: imageCache["decoration/switch"],
		actionFunc: func(g *Game) {
			game.mapdata[tx][ty] = nil
			playSound("switch", 1.0, false)
		},
	}
}
