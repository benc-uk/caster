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
			playSound("potion_get", 1, false)
		}
	}

	if kind == "crystal" {
		item.pickUpFunc = func(p *Player) {
			p.mana += 50
			playSound("zip_up", 1, false)
		}
	}

	if kind == "meat" {
		item.pickUpFunc = func(p *Player) {
			p.health += 25
			playSound("yum", 1, false)
		}
	}

	if kind == "apple" {
		item.pickUpFunc = func(p *Player) {
			p.health += 10
			playSound("gulp", 1, false)
		}
	}

	if kind == "key_red" || kind == "key_blue" || kind == "key_green" {
		item.pickUpFunc = func(p *Player) {
			p.holding[kind]++
			playSound("key_up", 1, false)
		}
	}

	// Very special case, these aren't items at all, but dungeon "furniture" which act like walls
	if kind == "column" || kind == "barrel" {
		item.pickUpFunc = func(p *Player) {
		}
		game.mapdata[cellX][cellY] = newInvisibleWall(cellX, cellY)
		g.stats.itemsTotal--
	}

	g.items[id] = item
	g.stats.itemsTotal++
}

func (g *Game) removeItem(i *Item) {
	screenFlashWhite(5)
	delete(g.items, i.id)
	g.removeSprite(i.sprite)
	i.sprite = nil
	i = nil
	g.stats.itemsFound++
}
