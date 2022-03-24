package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/wav"
)

var audioCtx *audio.Context
var sounds map[string]*audio.Player

func initSound() {
	log.Printf("Loading sounds...")
	sounds = make(map[string]*audio.Player, 10)

	audioCtx = audio.NewContext(44100)

	for _, sound := range []string{"footstep", "woohoo"} {
		f, err := os.Open("./sounds/" + sound + ".wav")
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
		sounds[sound] = p
	}
}

func playSound(sound string) {
	if sounds[sound] != nil {
		if sounds[sound].IsPlaying() {
			return
		}
		_ = sounds[sound].Rewind()
		sounds[sound].Play()
	}
}
