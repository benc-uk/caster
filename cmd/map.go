package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
)

// ===========================================================
// Map parser and loader
// ===========================================================
func loadMap(filename string, g *Game) {
	g.mapdata = make([][]int, mapSize)
	for i := range g.mapdata {
		g.mapdata[i] = make([]int, mapSize)
	}

	log.Printf("Loading map from %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	y := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		for x, c := range line {
			i, err := strconv.Atoi(string(c))
			if err != nil {
				g.mapdata[x][y] = 0

				// Player
				if c == 'P' {
					g.player.x = float64(x)*cellSize + cellSize/2
					g.player.y = float64(y)*cellSize + cellSize/2
				}

				// Monsters
				if c == 'g' || c == 's' || c == 't' {
					s := Sprite{
						x: float64(x)*cellSize + cellSize/2,
						y: float64(y)*cellSize + cellSize/2,
					}
					switch c {
					case 'g':
						s.id = "ghoul"
					case 's':
						s.id = "skeleton"
					case 't':
						s.id = "thing"
					}
					g.sprites = append(g.sprites, s)
				}
			} else {
				// Walls
				g.mapdata[x][y] = i
			}
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}
