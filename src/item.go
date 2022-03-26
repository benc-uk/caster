package main

import (
	"log"
	"math/rand"
)

type Item struct {
	id         uint64
	sprite     *Sprite
	pickUpFunc func(*Game)
}

func (g *Game) addItem(kind string, x, y float64, angle float64, speed float64, damage float64) {
	s := g.addSprite(kind, x, y, angle, speed, cellSize/16.0)

	id := rand.Uint64()
	g.items[id] = &Item{
		id:     id,
		sprite: s,
		pickUpFunc: func(g *Game) {
			g.player.health += 10
			log.Printf("Picked up item: %+v %s", id, s.kind)
		},
	}
}

func (g *Game) removeItem(i *Item) {
	delete(g.items, i.id)
	g.removeSprite(i.sprite)
	i.sprite = nil
	i = nil
}
