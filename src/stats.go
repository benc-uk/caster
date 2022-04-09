package main

import "time"

type Stats struct {
	monsters     int
	kills        int
	itemsTotal   int
	itemsFound   int
	secretsTotal int
	secretsFound int
	startTime    time.Time
	endTime      time.Time
}

func (s *Stats) init() {
	s.monsters = 0
	s.kills = 0
	s.itemsTotal = 0
	s.itemsFound = 0
	s.secretsTotal = 0
	s.secretsFound = 0
}
