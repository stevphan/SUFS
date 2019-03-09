package main

import (
	"Shared"
	"log"
	"net/http"
	"strconv"
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
	for i := 0; i < files.NumFiles; i++ {
		for j := 0; j < files.MetaData[i].NumBlocks; j++ {
			if len(files.MetaData[i].BlockLists[j].DnList) == 0 {
				log.Print("Dead Block: ", files.MetaData[i].FileName, "_", j, "\n")
			} else if len(files.MetaData[i].BlockLists[j].DnList) < replicationFactor {
				checkFailed(files.MetaData[i].FileName, j, i)
			}
		}
	}
	log.Println("repCheck complete")
}

func updateReplicationFactor() {
	if numDn == 0 { //There are no DN
		replicationFactor = 0
	} else if numDn < repFact { //don't have enough DN for replication factor
		replicationFactor = numDn
	} else { //Have enough DN for the replication factor
		replicationFactor = repFact
	}
}

func checkFailed(fileName string, blockId int, fileIndex int) {
	log.Print(fileName, "_", blockId, " replication check failed\n")
	myReq := shared.ReplicationRequest{}
	myReq.BlockId = fileName + "_" + strconv.Itoa(blockId)
	/*for i := 0; i < len(files.MetaData[fileIndex].BlockLists[blockId].DnList); i++ {
		myReq.DnList = append(myReq.DnList, files.MetaData[fileIndex].BlockLists[blockId].DnList[i])
		numDnListed++
	}*/
	//myReq.DnList = files.MetaData[fileIndex].BlockLists[blockId].DnList

	var foundDn bool
	numDnNeed := replicationFactor - len(files.MetaData[fileIndex].BlockLists[blockId].DnList)
	i := 0
	j := 0
	for i < numDn && numDnNeed > 0{
		foundDn = false
		j = 0
		for j < len(files.MetaData[fileIndex].BlockLists[blockId].DnList) && !foundDn{
			if files.MetaData[fileIndex].BlockLists[blockId].DnList[j] == dnList[i].dnIP {
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
	for k < len(files.MetaData[fileIndex].BlockLists[blockId].DnList) && !response { //For potential error in communicating
		//TODO send request to known Data Node (files.MetaData[fileIndex].BlockLists[blockId].DnList[j])
		dataNodeUrl := "http://" + files.MetaData[fileIndex].BlockLists[blockId].DnList[k] + "/replicateBlocks"
		buffer, err := shared.ConvertObjectToJsonBuffer(myReq)
		log.Print("Calling Dn at ", dataNodeUrl, "\n")
		_, err = http.Post(dataNodeUrl,"application/json", buffer)
		if err == nil {
			response = true
		}
		k++
	}
}