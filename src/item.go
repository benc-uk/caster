package main

import (
	"math/rand"
)

type Item struct {
	id         uint64
	sprite     *Sprite
	pickUpFunc func(*Player)
	cellX      int
	cellY      int
}

func (g *Game) addItem(kind string, cellX, cellY int, angle float64, speed float64, damage float64) {
	x := float64(cellX)*cellSize + cellSize/2
	y := float64(cellY)*cellSize + cellSize/2
	s := g.addSprite(kind, x, y, angle, speed, cellSize/16.0)

	id := rand.Uint64()
	item := &Item{
		id:         id,
		sprite:     s,
		cellX:      cellX,
		cellY:      cellY,
		pickUpFunc: func(p *Player) {},
	}

	if kind == "potion" {
		item.pickUpFunc = func(p *Player) {
			p.mana += 10
		}
	}

	g.items[id] = item
}

func (g *Game) removeItem(i *Item) {
	delete(g.items, i.id)
	g.removeSprite(i.sprite)
	i.sprite = nil
	i = nil
}
