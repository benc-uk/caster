package main

import (
	"log"
	"math"
)

type Player struct {
	x         float64 // Player position x
	y         float64 // Player position y
	angle     float64 // Facing angle in radians
	moveSpeed float64 // Current moving speed
	turnSpeed float64 // Current turning speed

	game             *Game // Pointer back to the game object
	playingFootsteps bool

	// These are effectively constants, but we hold them in the player
	fov            float64 // Field of view
	size           float64 // Used for collision detection with walls
	moveSpeedMin   float64 // Base move speed
	moveSpeedMax   float64 // Fastest we can go
	moveSpeedAccel float64 // Acceleration per keypress of move speed
	turnSpeedMin   float64 // Base turn speed
	turnSpeedMax   float64 // Max turn speed
	turnSpeedAccel float64 // Acceleration per keypress of turn speed
}

func newPlayer(g *Game) Player {
	return Player{
		x:                cellSize*1 + cellSize/2,
		y:                cellSize*1 + cellSize/2,
		angle:            0,
		moveSpeed:        cellSize / 32,
		moveSpeedMin:     cellSize / 32,
		moveSpeedMax:     cellSize / 6,
		moveSpeedAccel:   0.2,
		turnSpeed:        math.Pi / 120,
		turnSpeedMin:     math.Pi / 120,
		turnSpeedMax:     math.Pi / 40,
		turnSpeedAccel:   0.001,
		fov:              fov,
		size:             cellSize / 16.0,
		game:             g,
		playingFootsteps: false,
	}
}

func (p Player) checkWallCollision(x, y float64) (int, int, int) {
	if wall, x, y := p.game.getWallAt(x+p.size, y); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := p.game.getWallAt(x-p.size, y); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := p.game.getWallAt(x, y+p.size); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := p.game.getWallAt(x, y-p.size); wall > 0 {
		return wall, x, y
	}
	return 0, 0, 0
}

func (p Player) fireRay(distance float64) (int, int, int) {
	// Fire a ray in the direction we're facing
	for t := 0.0; t < distance; t += rayStepT {
		// Get hit point
		cx := p.x + (t * math.Cos(p.angle))
		cy := p.y + (t * math.Sin(p.angle))
		// Detect collision with walls
		wallIndex, wx, wy := p.game.getWallAt(cx, cy)
		if wallIndex > 0 {
			return wallIndex, wx, wy
		}
	}
	return 0, 0, 0
}

func (p Player) use() {
	wallIndex, wx, wy := p.fireRay(cellSize)
	if wallIndex > 0 {
		if wallIndex == 9 {
			p.game.mapdata[wx][wy] = 0
			playSound("door_open", 0.3, false)
		} else {
			playSound("grunt", 1, false)
		}
	}
}

func (p Player) attack() {
	playSound("zap", 0.3, false)
	log.Printf("Player.attack()")
	p.game.addProjectile("zap", p.x, p.y, p.angle, 4, 666)
}
