package main

import "log"

func logMessage(v ...interface{}) {
	if game.ticks%10 == 0 {
		log.Println(v...)
	}
}
