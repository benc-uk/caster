package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	x      float64 // Player position x
	y      float64 // Player position y
	cellX  int
	cellY  int
	angle  float64 // Facing angle in radians
	health int
	mana   int

	playingFootsteps bool

	// These are effectively constants, but we hold them in the player
	fov  float64 // Field of view
	size float64 // Used for collision detection with walls

	// These handle movement and turning
	moveStartTime int64
	moveFunc      func(int64) float64
	turnStartTime int64
	turnFunc      func(int64) float64

	holding map[string]int
}

func newPlayer(cellX, cellY int) Player {
	return Player{
		x:                cellSize*float64(cellX) + cellSize/2,
		y:                cellSize*float64(cellY) + cellSize/2,
		angle:            0,
		moveStartTime:    0.0,
		turnStartTime:    0.0,
		fov:              fov,
		size:             cellSize / 16.0,
		playingFootsteps: false,
		health:           100,
		mana:             100,
		cellX:            cellX,
		cellY:            cellY,
		holding:          map[string]int{},

		moveFunc: func(t int64) float64 {
			min := float64(cellSize) / 50.0
			max := float64(cellSize) / 14.0
			return math.Min(min+math.Pow(float64(t)/250000, 2), max)
		},

		turnFunc: func(t int64) float64 {
			min := math.Pi / 280.0
			max := math.Pi / 70.0
			return math.Min(min+math.Pow(float64(t)/800_000, 3), max)
		},
	}
}

func (p *Player) turn(t int64, direction float64) {
	p.angle = p.angle + p.turnFunc(t)*direction
}

func (p *Player) move(t int64, direction float64) {
	// Invoke the move function
	speed := p.moveFunc(t)

	if ebiten.IsKeyPressed(ebiten.KeyDown) || ebiten.IsKeyPressed(ebiten.KeyS) {
		speed = -speed
	}

	newX := p.x + math.Cos(p.angle)*speed
	newY := p.y + math.Sin(p.angle)*speed

	// Check if we're going to collide with a wall
	if wall := p.checkWallCollision(newX, newY); wall != nil {
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

func (p *Player) setFacing(facing int) {
	// facing: 0 = up, 1 = right, 2 = down, 3 = left
	p.angle = math.Pi / 2 * float64(facing-1)
}

func (p *Player) checkWallCollision(x, y float64) *Wall {
	if wall := game.getWallAt(x+p.size, y); wall != nil {
		return wall
	}
	if wall := game.getWallAt(x-p.size, y); wall != nil {
		return wall
	}
	if wall := game.getWallAt(x, y+p.size); wall != nil {
		return wall
	}
	if wall := game.getWallAt(x, y-p.size); wall != nil {
		return wall
	}
	return nil
}

func (p Player) fireRay(distance float64) *Wall {
	// Fire a ray in the direction we're facing
	for t := 0.0; t < distance; t += rayStepT {
		// Get hit point
		cx := p.x + (t * math.Cos(p.angle))
		cy := p.y + (t * math.Sin(p.angle))
		// Detect collision with walls
		wall := game.getWallAt(cx, cy)
		if wall != nil {
			return wall
		}
	}
	return nil
}

func (p Player) use() {
	wall := p.fireRay(cellSize)

	if wall != nil {
		wall.actionFunc(game)
	}
}

func (p *Player) attack() {
	if p.mana <= 0 {
		return
	}

	playSound("zap", 0.3, false)

	p.mana -= 5
	if p.mana < 0 {
		p.mana = 0.0
	}

	sx := p.x + ((cellSize / 3) * math.Cos(p.angle))
	sy := p.y + ((cellSize / 3) * math.Sin(p.angle))
	game.addProjectile("magic_1", sx, sy, p.angle, (float64(cellSize) / 9.0), 40)
}

// damage the player
func (p *Player) damage(amount int) {
	screenFlashRed(10)
	p.health -= amount
	if p.health <= 0 {
		p.health = 0
		playSound("scream", 1, false)
		titleScreen = true
	}
	playSound("pain", 1, false)
}
