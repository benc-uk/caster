package main

import (
	"io"
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

	for _, fileEntry := range wavDir {
		file, err := os.Open("./sounds/" + fileEntry.Name())
		if err != nil {
			log.Fatal(err)
		}

		var audioStream io.Reader
		if strings.HasPrefix(fileEntry.Name(), "loop") {
			wavStream, err := wav.DecodeWithSampleRate(44100, file)
			if err != nil {
				log.Fatal(err)
			}
			audioStream = audio.NewInfiniteLoop(file, wavStream.Length())
		} else {
			audioStream, err = wav.DecodeWithSampleRate(44100, file)
			if err != nil {
				log.Fatal(err)
			}
		}

		player, err := audioCtx.NewPlayer(audioStream)
		if err != nil {
			log.Fatal(err)
		}
		sounds[strings.TrimSuffix(fileEntry.Name(), ".wav")] = player
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

func soundStartAmbience() {
	sounds["loop_ambient_1"].Play()
}

func soundStopAmbience() {
	sounds["loop_ambient_1"].Pause()
	_ = sounds["loop_menu"].Rewind()
}

func soundStopTitleScreen() {
	sounds["loop_menu"].Pause()
	_ = sounds["loop_menu"].Rewind()
}

func soundStartTitleScreen() {
	sounds["loop_menu"].SetVolume(0.5)
	sounds["loop_menu"].Play()
}
