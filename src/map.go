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

			// Walls and decorations, switches etc
			if cell.Type == "w" {
				g.mapdata[cell.X][cell.Y] = newWall(cell.X, cell.Y, cell.Value)
				if len(cell.Extra) > 0 {
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
					g.mapdata[cell.X][cell.Y].metadata = append(g.mapdata[cell.X][cell.Y].metadata, cell.Extra...)
				}
			}

			// Doors
			if cell.Type == "d" {
				g.mapdata[cell.X][cell.Y] = newDoor(cell.X, cell.Y, cell.Value)
			}

			// Monsters
			if cell.Type == "m" {
				g.addMonster(cell.Value, cell.X, cell.Y)
			}

			// Items
			if cell.Type == "i" {
				g.addItem(cell.Value, cell.X, cell.Y)
			}

			// Player start point
			if cell.Type == "p" {
				g.player.x = float64(cell.X*cellSize + cellSize/2)
				g.player.y = float64(cell.Y*cellSize + cellSize/2)
				facing, _ := strconv.Atoi(cell.Value)
				g.player.setFacing(facing)
				log.Printf("Player spawn at %d,%d - Facing:%d", cell.X, cell.Y, facing)
			}
		}
	}

	g.mapName = name
}
