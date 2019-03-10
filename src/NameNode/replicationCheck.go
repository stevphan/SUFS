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
	for key, value := range files.MetaData {
		for i := 0; i < len(value); i++ {
			if len(value[i].DnList) == 0 {
				log.Print("Dead Block: ", value[i].Id, "\n")
			} else if len(value[i].DnList) < replicationFactor {
				checkFailed(key, i)
			}
		}
	}

	//for i := 0; i < files.NumFiles; i++ {
	/*for i := 0; i < len(files.MetaData); i++ {
		//for j := 0; j < files.MetaData[i].NumBlocks; j++ {
		for j := 0; j < len(files.MetaData[i].BlockLists); j++ {
			if len(files.MetaData[i].BlockLists[j].DnList) == 0 {
				log.Print("Dead Block: ", files.MetaData[i].FileName, "_", j, "\n")
			} else if len(files.MetaData[i].BlockLists[j].DnList) < replicationFactor {
				checkFailed(files.MetaData[i].FileName, j, i)
			}
		}
	}*/
	log.Println("repCheck complete")
}

func updateReplicationFactor() {
	//if numDn == 0 { //There are no DN
	if len(dnList) == 0 { //There are no DN
		replicationFactor = 0
	//} else if numDn < repFact { //don't have enough DN for replication factor
	} else if len(dnList) < repFact { //don't have enough DN for replication factor
		//replicationFactor = numDn
		replicationFactor = len(dnList)
	} else { //Have enough DN for the replication factor
		replicationFactor = repFact
	}
}

func checkFailed(fileName string, blockIndex int) {
	tempBlocks := files.MetaData[fileName]
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
		dataNodeUrl := "http://" + tempBlocks[blockIndex].DnList[k] + "/replicateBlocks"
		buffer, err := shared.ConvertObjectToJsonBuffer(myReq)
		log.Print("Calling Dn at ", dataNodeUrl, "\n")
		_, err = http.Post(dataNodeUrl,"application/json", buffer)
		if err == nil {
			response = true
		}
		k++
	}
}