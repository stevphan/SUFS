package main

import (
	"fmt"
	"time"
)

func replicationCheck() {
	ticker := time.NewTicker(1000 * time.Millisecond)
	for range ticker.C {
		fmt.Println("For loop")
	}
}
