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
	defer lock.Unlock()
	lock.Lock()
	js, _ := json.MarshalIndent(files, "", " ")
	err := ioutil.WriteFile(saveData, js, 0644)
	errorPrint(err)
}

func readFilesFromDisk() {
	files.MetaData = make(map[string][]blocks)
	files.LastId = 0
	tempFile, err := os.Open(saveData)
	if err != nil {
		errorPrint(err)
		writeFilesToDisk()
		return
	}
	defer tempFile.Close()
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

func addToDnList(ip string) {
	tempDn := dataNodeList{}
	tempDn.dnIP = ip
	tempDn.dnTime = time.Now()
	dnList = append(dnList, tempDn)
	log.Print("Added ", ip, " to dnList\n")
}

func getNewBlockId() int64 {
	files.LastId++
	return files.LastId
}

func readMap(key string) ([]blocks, bool) {
	lock.RLock()
	block, ok := files.MetaData[key]
	lock.RUnlock()
	return block, ok
}

func writeMap(key string, inputValue []blocks) {
	lock.Lock()
	files.MetaData[key] = inputValue
	lock.Unlock()
}