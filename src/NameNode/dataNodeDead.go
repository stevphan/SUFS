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
	//for i := 0; i < numDn; i++ {
	for i := 0; i < len(dnList); i++ {
		if time.Now().Sub(dnList[i].dnTime) > dnDeadAfter {
			removeFromDnList(i)
		}
	}
	//log.Println(dnList)
	log.Println("dnCheck complete")
}

//Removes the DN at the given index from the dnList
func removeFromDnList(dnIndex int) {
	log.Print("Removing ", dnList[dnIndex].dnIP, " from dnList\n")
	go deleteFromFiles(dnList[dnIndex].dnIP)
	dnList[dnIndex] = dnList[len(dnList)-1]
	dnList[len(dnList)-1] = dataNodeList{}
	dnList = dnList[:len(dnList)-1]
	//numDn--
}

//Removes given IP from places a block is stored
func deleteFromFiles(ip string) {

	var foundIp bool
	var j int
	for key, value := range files.MetaData {
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
		files.MetaData[key] = value
	}
	//for i := 0; i < files.NumFiles; i++ {
	/*for i := 0; i < len(files.MetaData); i++ {
		//for j := 0; j < files.MetaData[i].NumBlocks; j++ {
		for j := 0; j < len(files.MetaData[i].BlockLists); j++ {
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
	}*/
	writeFilesToDisk()
}
