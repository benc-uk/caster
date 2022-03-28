package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// ===========================================================
// Map parser and loader
// ===========================================================
func (g *Game) loadMap(name string) {
	g.mapdata = make([][]*Wall, mapSize)
	for i := range g.mapdata {
		g.mapdata[i] = make([]*Wall, mapSize)
	}

	filename := "./maps/" + name + ".map"
	log.Printf("Loading map from: %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	g.mapName = name

	// How we find monsters and items
	monsterRe := regexp.MustCompile("[a-z]")
	itemRe := regexp.MustCompile("[A-Z]")

	// Read the file line by line
	y := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Special commands: CEIL & FLOOR to set the coloring of the ceiling and floor
		if strings.HasPrefix(line, "CEIL:") || strings.HasPrefix(line, "FLOOR:") {
			re := regexp.MustCompile(`:\s*(.*?),\s*(.*?),\s*(.*?)$`)
			matches := re.FindAllStringSubmatch(line, -1)
			p1, _ := strconv.ParseFloat(matches[0][1], 64)
			p2, _ := strconv.ParseFloat(matches[0][2], 64)
			p3, _ := strconv.ParseFloat(matches[0][3], 64)

			if line[0] == 'C' {
				ceilOp.ColorM.Scale(p1, p2, p3, 1)
			} else {
				floorOp.ColorM.Scale(p1, p2, p3, 1)
			}
		}

		if strings.HasPrefix(line, "NAME:") {
			g.mapName = line[5:]
		}

		// Process the line and store it in the map
		for x, char := range line {
			if isNumeric(string(char)) {
				// Walls are numbered from 1 to 9
				g.mapdata[x][y] = newWall(x, y, string(char))
			} else {
				g.mapdata[x][y] = nil

				// Player
				if char == '*' {
					g.player.moveToCell(x, y)
				}

				// Doors are a special case, and start at 20
				if char == '#' {
					g.mapdata[x][y] = newDoor(x, y, "basic")
				}

				if monsterRe.MatchString(string(char)) {
					switch char {
					case 'g':
						g.addMonster("ghoul", x, y)
					case 's':
						g.addMonster("skeleton", x, y)
					case 't':
						g.addMonster("thing", x, y)
					case 'P':
						g.addMonster("potion", x, y)
					}
				}

				if itemRe.MatchString(string(char)) {
					switch char {
					case 'P':
						g.addItem("potion", x, y, 1.7, 0, 1.0)
					case 'B':
						g.addItem("ball", x, y, 1.7, 0, 1.0)
					}
				}
			}
		}
		y++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func isNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
