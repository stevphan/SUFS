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
	time.Sleep(30 * time.Second)
	ticker := time.NewTicker(dnCheckFreq * time.Millisecond)
	for range ticker.C {
		log.Println("Calling dnCheck")
		dnCheck()
	}
}

func dnCheck() {
	for i := 0; i < len(dnList); i++ {
		if time.Now().Sub(dnList[i].dnTime) > dnDeadAfter {
			removeFromDnList(i)
		}
	}
	log.Println("dnCheck complete")
}

//Removes the DN at the given index from the dnList
func removeFromDnList(dnIndex int) {
	log.Print("Removing ", dnList[dnIndex].dnIP, " from dnList\n")
	go deleteFromFiles(dnList[dnIndex].dnIP)
	dnList[dnIndex] = dnList[len(dnList)-1]
	dnList[len(dnList)-1] = dataNodeList{}
	dnList = dnList[:len(dnList)-1]
}

//Removes given IP from places a block is stored
func deleteFromFiles(ip string) {

	var foundIp bool
	var j int
	lock.RLock() //lock before for loop
	for key, value := range files.MetaData {
		lock.RUnlock() //unlock for reads
		for i := 0; i < len(value); i++ {
			j = 0
			foundIp = false
			for j < len(value[i].DnList) && !foundIp {
				if value[i].DnList[j] == ip { //remove from list
					value[i].DnList[j] = value[i].DnList[len(value[i].DnList)-1]
					value[i].DnList[len(value[i].DnList)-1] = ""
					value[i].DnList = value[i].DnList[:len(value[i].DnList)-1]
					foundIp = true
				}
				j++
			}
		}
		//files.MetaData[key] = value
		writeMap(key, value)
		lock.RLock() //lock for for loop
	}
	lock.RUnlock()
	writeFilesToDisk()
}
