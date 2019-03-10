package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

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
	//if files.NumFiles > 0 {
	if len(files.MetaData) > 0 {
		//for i := 0; i < files.NumFiles; i++ {
		for i := 0; i < len(files.MetaData); i++ {
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
	//numDn++
}