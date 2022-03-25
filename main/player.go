package main

import (
	"math"
)

type Player struct {
	x              float64 // Player position x
	y              float64 // Player position y
	angle          float64 // Facing angle in radians
	moveSpeed      float64 // How fast we move
	moveSpeedMin   float64 // How fast we move
	moveSpeedMax   float64 // How fast we move
	moveSpeedAccel float64 // How fast we move
	turnSpeed      float64 // How much we turn per keypress
	turnSpeedMin   float64 // How much we turn per keypress
	turnSpeedMax   float64 // How much we turn per keypress
	turnSpeedAccel float64 // How much we turn per keypress
	fov            float64 // Field of view
	size           float64 // Used for collision detection with walls
	game           *Game   // Pointer back to the game object
}

func newPlayer(g *Game) Player {
	return Player{
		x:              cellSize*1 + cellSize/2,
		y:              cellSize*1 + cellSize/2,
		angle:          0,
		moveSpeed:      cellSize / 32,
		moveSpeedMin:   cellSize / 32,
		moveSpeedMax:   cellSize / 6,
		moveSpeedAccel: 0.2,
		turnSpeed:      math.Pi / 120,
		turnSpeedMin:   math.Pi / 120,
		turnSpeedMax:   math.Pi / 40,
		turnSpeedAccel: 0.002,
		fov:            fov,
		size:           cellSize / 16.0,
		game:           g,
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

func (p Player) use() {
	const useDist = cellSize * 0.5
	const normalDoor = 9
	// const runeDoor = 8
	newX := p.x + math.Cos(p.angle)*useDist
	newY := p.y + math.Sin(p.angle)*useDist

	// Check wall that was "used"
	if wall, x, y := p.checkWallCollision(newX, newY); wall > 0 {
		p.game.mapdata[x][y] = 0
		if wall == normalDoor {

			playSound("door_open")
		} else {
			playSound("grunt")
		}
	}

}
