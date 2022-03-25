package main

type Monster struct {
	sprite *Sprite
	health int
}

func (g *Game) addMonster(id string, x, y float64, angle float64, speed float64, size float64) {
	s := g.addSprite(id, x, y, angle, speed, size)

	g.monsters = append(g.monsters, Monster{
		sprite: s,
		health: 10,
	})
}
