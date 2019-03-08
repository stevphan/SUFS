package main

import (
	"Shared"
	"log"
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
				//TODO figure out what happens in this case, maybe after two fails delete file
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
	for i < numDn && numDnNeed > 0{
		foundDn = false
		for j := 0; j < len(files.MetaData[fileIndex].BlockLists[blockId].DnList); j++ {
			if files.MetaData[fileIndex].BlockLists[blockId].DnList[j] == dnList[i] {
				foundDn = true
			}
		}
		if !foundDn {
			myReq.DnList = append(myReq.DnList, dnList[i])
			numDnNeed--
		}
		i++
	}

	j := 0
	response := false
	for j < len(files.MetaData[fileIndex].BlockLists[blockId].DnList) && !response { //For potential error in communicating
		//TODO send request to known Data Node (files.MetaData[fileIndex].BlockLists[blockId].DnList[j])
		//if gotten response, response = true
		j++
	}
	//Don't care if never communicate, will eventually fix
	//TODO receive response from Data Node
}