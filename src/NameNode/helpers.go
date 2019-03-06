package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
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