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

func (g *Game) addItem(kind string, cellX, cellY int) {
	x := float64(cellX)*cellSize + cellSize/2
	y := float64(cellY)*cellSize + cellSize/2
	s := g.addSprite("items/"+kind, x, y, 0, 0, cellSize/16.0)

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
			p.mana += 25
			playSound("zip_up", 1, false)
		}
	}

	if kind == "ball" {
		item.pickUpFunc = func(p *Player) {
			p.mana += 50
			playSound("zip_up", 1, false)
		}
	}

	if kind == "key_red" || kind == "key_blue" || kind == "key_green" {
		item.pickUpFunc = func(p *Player) {
			p.holding[kind]++
			playSound("woohoo", 1, false)
		}
	}

	g.items[id] = item
}

func (g *Game) removeItem(i *Item) {
	screenFlashWhite(5)
	delete(g.items, i.id)
	g.removeSprite(i.sprite)
	i.sprite = nil
	i = nil
}
