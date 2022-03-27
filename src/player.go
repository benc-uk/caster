package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	x         float64 // Player position x
	y         float64 // Player position y
	cellX     int
	cellY     int
	angle     float64 // Facing angle in radians
	moveSpeed float64 // Current moving speed
	turnSpeed float64 // Current turning speed
	health    int
	mana      int

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

func newPlayer(cellX, cellY int) Player {
	return Player{
		x:                cellSize*float64(cellX) + cellSize/2,
		y:                cellSize*float64(cellY) + cellSize/2,
		angle:            0,
		moveSpeed:        cellSize / 32,
		moveSpeedMin:     cellSize / 32,
		moveSpeedMax:     cellSize / 8,
		moveSpeedAccel:   0.15,
		turnSpeed:        math.Pi / 150,
		turnSpeedMin:     math.Pi / 150,
		turnSpeedMax:     math.Pi / 40,
		turnSpeedAccel:   0.0005,
		fov:              fov,
		size:             cellSize / 16.0,
		playingFootsteps: false,
		health:           10000,
		mana:             10000,
		cellX:            cellX,
		cellY:            cellY,
	}
}

func (p *Player) move() {
	ms := p.moveSpeed
	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		ms = -p.moveSpeed
	}
	p.moveSpeed = math.Min(p.moveSpeed+p.moveSpeedAccel, p.moveSpeedMax)

	newX := p.x + math.Cos(p.angle)*ms
	newY := p.y + math.Sin(p.angle)*ms

	// Check if we're going to collide with a wall
	if wall, _, _ := p.checkWallCollision(newX, newY); wall > 0 {
		// Hit a wall so don't move
		return
	}

	// Update player position
	p.x = newX
	p.y = newY
	p.cellX = int(math.Floor(p.x / cellSize))
	p.cellY = int(math.Floor(p.y / cellSize))

	// Check items near the player we're in and pick them up
	for _, item := range game.items {
		if item.cellX != p.cellX || item.cellY != p.cellY {
			continue
		}

		item.pickUpFunc(p)
		playSound("zap", 0.3, false)
		game.removeItem(item)
	}

	// Footstep sound
	if !p.playingFootsteps {
		playSound(fmt.Sprintf("footstep_%d", rand.Intn(4)), 0.5, true)
		p.playingFootsteps = true

		time.AfterFunc(300*time.Millisecond, func() {
			p.playingFootsteps = false
		})
	}
}

func (p *Player) moveToCell(cellX, cellY int) {
	p.cellX = cellX
	p.cellY = cellY
	p.x = cellSize*float64(cellX) + cellSize/2
	p.y = cellSize*float64(cellY) + cellSize/2
}

func (p Player) checkWallCollision(x, y float64) (int, int, int) {
	if wall, x, y := game.getWallAt(x+p.size, y); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := game.getWallAt(x-p.size, y); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := game.getWallAt(x, y+p.size); wall > 0 {
		return wall, x, y
	}
	if wall, x, y := game.getWallAt(x, y-p.size); wall > 0 {
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
		wallIndex, wx, wy := game.getWallAt(cx, cy)
		if wallIndex > 0 {
			return wallIndex, wx, wy
		}
	}
	return 0, 0, 0
}

func (p Player) use() {
	wallIndex, wx, wy := p.fireRay(cellSize)
	if wallIndex > 0 {
		if wallIndex == doorWallIndex {
			game.mapdata[wx][wy] = 0
			playSound("door_open", 0.3, false)
		} else {
			playSound("grunt", 1, false)
		}
	}
}

func (p *Player) attack() {
	if p.mana <= 0 {
		return
	}
	playSound("zap", 0.3, false)

	p.mana -= 10
	if p.mana < 0 {
		p.mana = 0.0
	}
	game.addProjectile("magic_1", p.x, p.y, p.angle, 3, 40) //p.moveSpeedMax*1.1, 40)
}
