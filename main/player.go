package main

import (
	"math"
)

type Player struct {
	x         float64 // Player position x
	y         float64 // Player position y
	angle     float64 // Facing angle in radians
	moveSpeed float64 // How fast we move
	turnSpeed float64 // How much we turn per keypress
	fov       float64 // Field of view
	size      float64 // Used for collision detection with walls
	game      *Game   // Pointer back to the game object
}

func newPlayer(g *Game) Player {
	return Player{
		x:         cellSize*1 + cellSize/2,
		y:         cellSize*1 + cellSize/2,
		angle:     0,
		moveSpeed: cellSize / 8.0,
		turnSpeed: math.Pi / 80,
		fov:       fov,
		size:      cellSize / 16.0,
		game:      g,
	}
}

func (p Player) checkWallCollision(x, y float64) int {
	if wall := p.game.getWallAt(x+p.size, y); wall > 0 {
		return wall
	}
	if wall := p.game.getWallAt(x-p.size, y); wall > 0 {
		return wall
	}
	if wall := p.game.getWallAt(x, y+p.size); wall > 0 {
		return wall
	}
	if wall := p.game.getWallAt(x, y-p.size); wall > 0 {
		return wall
	}
	return 0
}

func (p Player) use() {
	const useDist = cellSize
	const normalDoor = 9
	// const runeDoor = 8
	newX := p.x + math.Cos(p.angle)*useDist
	newY := p.y + math.Sin(p.angle)*useDist

	// Check wall that was "used"
	if wall := p.checkWallCollision(newX, newY); wall > 0 {
		if wall == normalDoor {
			p.game.mapdata[int(newX/cellSize)][int(newY/cellSize)] = 0
			playSound("door_open")
		} else {
			playSound("grunt")
		}
	}

}
