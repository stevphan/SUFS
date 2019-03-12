package main

import (
	"shared"
	"log"
	"net/http"
	"time"
)

const (
	checkTime = 60000 //in ms, 1 min
)

var (
	replicationFactor int
)

func replicationCheck() {
	ticker := time.NewTicker(checkTime * time.Millisecond)
	for range ticker.C {
		log.Println("Calling repCheck")
		repCheck()
	}
}

func repCheck() {
	updateReplicationFactor()
	lock.RLock() //lock before for loop
	for key, value := range files.MetaData {
		lock.RUnlock() //unlock for reads in checkFailed
		for i := 0; i < len(value); i++ {
			if len(value[i].DnList) == 0 {
				log.Print("Dead Block: ", value[i].Id, "\n")
			} else if len(value[i].DnList) < replicationFactor {
				checkFailed(key, i)
			}
		}
		lock.RLock() //lock for for loop
	}
	lock.RUnlock()
	log.Println("repCheck complete")
}

func updateReplicationFactor() {
	if len(dnList) == 0 { //There are no DN
		replicationFactor = 0
	} else if len(dnList) < repFact { //don't have enough DN for replication factor
		replicationFactor = len(dnList)
	} else { //Have enough DN for the replication factor
		replicationFactor = repFact
	}
}

func checkFailed(fileName string, blockIndex int) {
	//tempBlocks := files.MetaData[fileName]
	tempBlocks, _ := readMap(fileName)
	log.Print("Replication check failed for ", fileName, ": blockId=", tempBlocks[blockIndex].Id, "\n")
	myReq := shared.ReplicationRequest{}
	myReq.BlockId = tempBlocks[blockIndex].Id

	var foundDn bool
	numDnNeed := replicationFactor - len(tempBlocks[blockIndex].DnList)
	i := 0
	j := 0
	for i < len(dnList) && numDnNeed > 0 {
		foundDn = false
		j = 0
		for j < len(tempBlocks[blockIndex].DnList) && !foundDn{
			if tempBlocks[blockIndex].DnList[j] == dnList[i].dnIP {
				foundDn = true
			}
			j++
		}
		if !foundDn {
			myReq.DnList = append(myReq.DnList, dnList[i].dnIP)
			numDnNeed--
		}
		i++
	}

	k := 0
	response := false
	for k < len(tempBlocks[blockIndex].DnList) && !response { //For potential error in communicating
		dataNodeUrl := "http://" + tempBlocks[blockIndex].DnList[k] + shared.PathReplication
		buffer, err := shared.ConvertObjectToJsonBuffer(myReq)
		log.Print("Calling Dn at ", dataNodeUrl, "\n")
		_, err = http.Post(dataNodeUrl,"application/json", buffer)
		if err == nil {
			response = true
		}
		k++
	}
}