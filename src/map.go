package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strconv"
)

type MapFileCell struct {
	X     int
	Y     int
	Type  string   `json:"t"`
	Value string   `json:"v"`
	Extra []string `json:"e"`
}

// ===========================================================
// Map parser and loader
// ===========================================================
func (g *Game) loadMap(name string) {
	// Load the map file
	data, err := ioutil.ReadFile("./maps/" + name + ".json")
	if err != nil {
		log.Fatal(err)
	}

	// Raw map holds the unmarshalled map data from JSON
	mapRaw := make([][]*MapFileCell, 0)
	err = json.Unmarshal([]byte(data), &mapRaw)
	if err != nil {
		log.Fatal(err)
	}

	// This is the real map data used by the game
	g.mapdata = make([][]*Wall, mapSize)
	for i := range g.mapdata {
		g.mapdata[i] = make([]*Wall, mapSize)
	}

	// Parse the raw map into the mapdata
	for _, cellRow := range mapRaw {
		for _, cell := range cellRow {
			g.mapdata[cell.X][cell.Y] = nil
			if cell.Type == "w" {
				g.mapdata[cell.X][cell.Y] = newWall(cell.X, cell.Y, cell.Value)
				if len(cell.Extra) > 0 {
					g.mapdata[cell.X][cell.Y].metadata = cell.Extra
					if cell.Extra[0] == "deco" {
						g.mapdata[cell.X][cell.Y].decoration = imageCache["decoration/"+cell.Extra[1]]
					}
					if cell.Extra[0] == "secret" {
						g.mapdata[cell.X][cell.Y] = newSecretWall(cell.X, cell.Y, cell.Value)
					}
					if cell.Extra[0] == "switch" {
						targetX, _ := strconv.Atoi(cell.Extra[1])
						targetY, _ := strconv.Atoi(cell.Extra[2])
						g.mapdata[cell.X][cell.Y] = newSwitchWall(cell.X, cell.Y, cell.Value, targetX, targetY)
					}
				}
			}

			if cell.Type == "d" {
				g.mapdata[cell.X][cell.Y] = newDoor(cell.X, cell.Y, cell.Value)
			}

			if cell.Type == "m" {
				g.addMonster(cell.Value, cell.X, cell.Y)
			}

			if cell.Type == "i" {
				g.addItem(cell.Value, cell.X, cell.Y)
			}

			if cell.Type == "p" {
				g.player.x = float64(cell.X*cellSize + cellSize/2)
				g.player.y = float64(cell.Y*cellSize + cellSize/2)
				f, _ := strconv.Atoi(cell.Value)
				g.player.setFacing(f)
				log.Printf("Player spawn at %d,%d, %d", cell.X, cell.Y, f)
			}
		}
	}

	g.mapName = name
}

/*func (g *Game) loadMap(name string) {
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
}*/

// func isNumeric(s string) bool {
// 	_, err := strconv.ParseFloat(s, 64)
// 	return err == nil
// }
