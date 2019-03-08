package main

import (
	"log"
	"time"
)

const (
	dnCheckFreq = 60000 //in ms, 1 min
	dnDeadAfter time.Duration = 120000000000 //2 min, in nanoseconds
)

func dataNodeDead() {
	ticker := time.NewTicker(dnCheckFreq * time.Millisecond)
	for range ticker.C {
		log.Println("Calling dnCheck")
		dnCheck()
	}
}

func dnCheck() {
	for i := 0; i < numDn; i++ {
		if time.Now().Sub(dnList[i].dnTime) > dnDeadAfter {
			removeFromDnList(i)
		}
	}
	log.Println(dnList)
	log.Println("dnCheck complete")
}
