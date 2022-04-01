package main

import (
	"log"
	"strings"

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
	if _, ok := imageCache["walls/"+kind]; !ok {
		log.Fatalf("ERROR! Wall image not found: %s", kind)
	}

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
	door := &Wall{
		x:     x,
		y:     y,
		image: imageCache["doors/"+kind],

		// Default is a locked door
		actionFunc: func(g *Game) {
			playSound("locked", 1.0, false)
		},
	}

	// Basic doors can just be opened
	if kind == "basic" {
		door.actionFunc = func(g *Game) {
			game.mapdata[x][y] = nil
			playSound("door_open", 0.4, false)
		}
	}

	if strings.HasPrefix(kind, "key") {
		door.actionFunc = func(g *Game) {
			count, holding := g.player.holding[kind]
			if holding && count > 0 {
				game.mapdata[x][y] = nil
				playSound("unlock", 1.0, false)
				g.player.holding[kind]--
			} else {
				playSound("locked", 1.0, false)
			}
		}
	}

	return door
}

func newSecretWall(x, y int, kind string) *Wall {
	return &Wall{
		x:     x,
		y:     y,
		image: imageCache["walls/"+kind],

		// Remove this wall
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

		// Remove a wall somewhere else
		actionFunc: func(g *Game) {
			wall := game.mapdata[x][y]
			if wall.metadata[0] == "pressed" {
				playSound("grunt", 1.0, false)
				return
			}
			game.mapdata[tx][ty] = nil
			playSound("switch", 1.0, false)
			wall.decoration = imageCache["decoration/switch-1"]
			wall.metadata[0] = "pressed"
		},

		metadata: []string{
			"not_pressed",
		},
	}
}
