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
	seen       bool
	isDoor     bool
	invisible  bool

	actionFunc func(g *Game)
}

func (w *Wall) getCenter() (float64, float64) {
	return float64(w.x*cellSize + cellSize/2), float64(w.y*cellSize + cellSize/2)
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
		x:      x,
		y:      y,
		image:  imageCache["doors/"+kind],
		isDoor: true,

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
	game.stats.secretsTotal++
	return &Wall{
		x:     x,
		y:     y,
		image: imageCache["walls/"+kind],

		// Remove this wall
		actionFunc: func(g *Game) {
			game.mapdata[x][y] = nil
			playSound("secret", 1.0, false)
			game.stats.secretsFound++
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

func newExitWall(x, y int, kind string) *Wall {
	return &Wall{
		x:          x,
		y:          y,
		image:      imageCache["walls/"+kind],
		decoration: imageCache["decoration/exit"],

		// Remove this wall
		actionFunc: func(g *Game) {
			g.endLevel()
		},
	}
}

func newInvisibleWall(x, y int) *Wall {
	return &Wall{
		x:          x,
		y:          y,
		image:      nil,
		decoration: nil,
		invisible:  true,
	}
}
