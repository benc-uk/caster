package main

type Game struct {
	// Map is a 2D array of integers, 0 = empty, 1+ = wall
	mapdata   [][]int
	mapWidth  int
	mapHeight int
	player    Player
	level     int
	sprites   []Sprite
}

type Player struct {
	x         float64
	y         float64
	angle     float64
	moveSpeed float64
	turnSpeed float64
	fov       float64
}
