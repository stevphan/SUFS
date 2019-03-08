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

//Removes the DN at the given index from the dnList
func removeFromDnList(dnIndex int) {
	log.Print("Removing ", dnList[dnIndex].dnIP, " from dnList\n")
	go deleteFromFiles(dnList[dnIndex].dnIP)
	dnList[dnIndex] = dnList[len(dnList)-1]
	dnList[len(dnList)-1] = dataNodeList{}
	dnList = dnList[:len(dnList)-1]
	numDn--
}

//Removes given IP from places a block is stored
func deleteFromFiles(ip string) {
	k := 0
	foundIp := false
	for i := 0; i < files.NumFiles; i++ {
		for j := 0; j < files.MetaData[i].NumBlocks; j++ {
			k = 0
			foundIp = false
			for k < len(files.MetaData[i].BlockLists[j].DnList) && !foundIp {
				if files.MetaData[i].BlockLists[j].DnList[k] == ip { //remove from list
					files.MetaData[i].BlockLists[j].DnList[k] = files.MetaData[i].BlockLists[j].DnList[len(files.MetaData[i].BlockLists[j].DnList)-1]
					files.MetaData[i].BlockLists[j].DnList[len(files.MetaData[i].BlockLists[j].DnList)-1] = ""
					files.MetaData[i].BlockLists[j].DnList = files.MetaData[i].BlockLists[j].DnList[:len(files.MetaData[i].BlockLists[j].DnList)-1]
					foundIp = true
				}
				k++
			}
		}
	}
	writeFilesToDisk()
}
