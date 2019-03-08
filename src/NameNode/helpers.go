package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func convertObjectToJsonBuffer(object interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(data)

	return buffer, nil
}

func convertObjectToJson(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func writeFilesToDisk() {
	js, _ := json.MarshalIndent(files, "", " ")
	err := ioutil.WriteFile(saveData, js, 0644)
	//log.Println(files)
	errorPrint(err)
}

func readFilesFromDisk() {
	tempFile, err := os.Open(saveData)
	errorPrint(err)
	decoder := json.NewDecoder(tempFile)
	myFile := file{}
	err = decoder.Decode(&myFile)
	errorPrint(err)
	files = myFile
}

func errorPrint(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

//Returns true if file name found plus the index it is in the files.metadata[]
//Returns false if file name not found plus index of -1
func findFile(fileName string) (found bool, fileIndex int) {
	if files.NumFiles > 0 {
		for i := 0; i < files.NumFiles; i++ {
			if files.MetaData[i].FileName == fileName {
				return true, i
			}
		}
	}
	return false, -1
}

func addToDnList(ip string) {
	tempDn := dataNodeList{}
	tempDn.dnIP = ip
	tempDn.dnTime = time.Now()
	dnList = append(dnList, tempDn)
	log.Print("Added ", ip, " to dnList\n")
	numDn++
}

func removeFromDnList(dnIndex int) { //TODO need to remove it from places where block is stored
	log.Print("Removing ", dnList[dnIndex].dnIP, " from dnList\n")
	go deleteFromFiles(dnList[dnIndex].dnIP)
	dnList[dnIndex] = dnList[len(dnList)-1]
	dnList[len(dnList)-1] = dataNodeList{}
	dnList = dnList[:len(dnList)-1]
	numDn--
}

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