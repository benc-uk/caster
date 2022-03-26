package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
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

	monsterRe := regexp.MustCompile("[a-z]")
	itemRe := regexp.MustCompile("[A-Z]")

	y := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		for x, char := range line {
			i, err := strconv.Atoi(string(char))
			if err != nil {
				g.mapdata[x][y] = 0

				// Player
				if char == '*' {
					g.player.x = float64(x)*cellSize + cellSize/2
					g.player.y = float64(y)*cellSize + cellSize/2
				}

				if monsterRe.MatchString(string(char)) {
					log.Printf("Found monster %s at %d,%d", string(char), x, y)
					x := float64(x)*cellSize + cellSize/2
					y := float64(y)*cellSize + cellSize/2

					switch char {
					case 'g':
						g.addMonster("ghoul", x, y, 2.35, 0, 1.0)
					case 's':
						g.addMonster("skeleton", x, y, 1.0, 0, 1.0)
					case 't':
						g.addMonster("thing", x, y, 1.7, 0, 1.0)
					case 'P':
						g.addMonster("potion", x, y, 1.7, 0, 1.0)
					}
				}

				if itemRe.MatchString(string(char)) {
					log.Printf("Found item %s at %d,%d", string(char), x, y)
					x := float64(x)*cellSize + cellSize/2
					y := float64(y)*cellSize + cellSize/2

					switch char {
					case 'P':
						g.addMonster("potion", x, y, 1.7, 0, 1.0)
					case 'B':
						g.addMonster("ball", x, y, 1.7, 0, 1.0)
					}
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
