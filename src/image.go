package main

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

var imageCache map[string]*ebiten.Image

const gfxDir = "./gfx"

func loadImageCache() {
	imageCache = make(map[string]*ebiten.Image)

	imageDirEntry, err := os.ReadDir(gfxDir)
	if err != nil {
		log.Fatalln(err)
	}

	for _, subDir := range imageDirEntry {
		imageDirEntry, err := os.ReadDir(gfxDir + "/" + subDir.Name())
		if err != nil {
			log.Fatalln(err)
		}

		log.Printf("Loading %s images", subDir.Name())

		for _, file := range imageDirEntry {
			filename := gfxDir + "/" + subDir.Name() + "/" + file.Name()
			entryname := subDir.Name() + "/" + strings.TrimSuffix(file.Name(), ".png")
			if err != nil {
				log.Fatalln(err)
			}

			img, _, err := ebitenutil.NewImageFromFile(filename)
			if err != nil {
				log.Fatalln(err)
			}
			imageCache[entryname] = img
		}
	}
}
