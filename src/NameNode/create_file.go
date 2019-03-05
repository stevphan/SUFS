package main

import (
	"encoding/json"
	"math"
	"net/http"
	"strconv"
)

const (
	blockSize int64 = 67108864 // assuming bytes
)

func createFile(write http.ResponseWriter, req *http.Request) { //needs to return list of dataNodes per block
	var blocksRequired int

	decoder := json.NewDecoder(req.Body)
	myReq := createRequest{}
	err := decoder.Decode(&myReq)
	errorPrint(err)

	//Checks if the file exist
	if files.NumFiles > 0 {
		for i := 0; i < files.NumFiles; i++ {
			if files.MetaData[i].FileName == myReq.Filename {
				//TODO fail write (DONE?)
				myRes := responseObj{}
				js, err := convertObjectToJson(myRes)
				errorPrint(err)
				write.Header().Set("Content-Type", "application/json")
				write.Write(js)
				return
			}
		}
	}

	//Finds the blocks required
	size, err := strconv.ParseInt(myReq.Size, 10, 64)
	errorPrint(err)
	temp := float64(size)/float64(blockSize)
	temp = math.Ceil(temp)
	blocksRequired = int(temp)

	//TODO choose DN to send each block to (check size of, choose lowest) - not super important
	//Checks amount of DN vs the replication factor
	var replicationFactor int
	if numDn == 0 { //There are no DN
		//TODO fail write (DONE?)
		myRes := responseObj{}
		js, err := convertObjectToJson(myRes)
		errorPrint(err)
		write.Header().Set("Content-Type", "application/json")
		write.Write(js)
		return
	} else if numDn < repFact { //don't have enough DN for replication factor
		replicationFactor = numDn
	} else { //Have enough DN for the replication factor
		replicationFactor = repFact
	}

	//This chooses DN for each block
	j := 0 //index of the DataNode list
	myRes := responseObj{}
	myRes.Blocks = make([]createResponse, blocksRequired)
	for i := 0; i < blocksRequired; i++ {
		blockList := createResponse{}
		blockList.BlockId = myReq.Filename + "_" + strconv.Itoa(i)
		blockList.DnList = make([]string, replicationFactor)
		for k := 0; k < replicationFactor; k++ {
			blockList.DnList[k] = dnList[j]
			j++
			if j == numDn {
				j = 0
			}
		}
		myRes.Blocks[i] = blockList
	}

	//Saves the metadata of the file, nothing in DnList yet
	fileToStore := fileMetaData{}
	fileToStore.FileName = myReq.Filename
	fileToStore.NumBlocks = blocksRequired
	fileToStore.BlockLists = make([]blockList, blocksRequired)
	for i := 0; i < blocksRequired; i++ {
		blockList := blockList{}
		for j := 0; j < repFact; j++ {
			blockList.DnList[j] = ""
		}
		fileToStore.BlockLists[i] = blockList
	}
	files.NumFiles++
	files.MetaData = append(files.MetaData, fileToStore)
	writeFilesToDisk()

	//Returns myRes which is a responseObj
	js, err := convertObjectToJson(myRes)
	errorPrint(err)
	write.Header().Set("Content-Type", "application/json")
	write.Write(js)
}