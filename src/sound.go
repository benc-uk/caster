package main

import (
	"log"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var audioCtx *audio.Context
var sounds map[string]*audio.Player

func initSound() {
	log.Printf("Loading sounds...")
	sounds = make(map[string]*audio.Player, 10)

	audioCtx = audio.NewContext(44100)

	wavDir, err := os.ReadDir("./sounds")
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range wavDir {
		f, err := os.Open("./sounds/" + file.Name())
		if err != nil {
			log.Fatal(err)
		}
		d, err := wav.DecodeWithSampleRate(44100, f)
		if err != nil {
			log.Fatal(err)
		}
		p, err := audioCtx.NewPlayer(d)
		if err != nil {
			log.Fatal(err)
		}
		sounds[strings.TrimSuffix(file.Name(), ".wav")] = p
	}
}

func playSound(sound string, volume float64, wait bool) {
	if sounds[sound] != nil {
		if sounds[sound].IsPlaying() && wait {
			return
		}
		_ = sounds[sound].Rewind()
		sounds[sound].SetVolume(volume)
		sounds[sound].Play()
	}
}
