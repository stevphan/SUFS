package main

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	repFact = 3
)

type file struct {
	numFiles int
	metaData []fileMetaData
}

type fileMetaData struct {
	fileName string
	numBlocks int
	blockLists []blockList
}

type blockList struct {
	dnList [repFact]string //stores IP of DNs it is stored at
}

var (
	dnList = []string{}
	numDn = 0
	files file
)

func main() {
	//TODO read from disk, build the files from that
	//for testing
	dnList = append(dnList, "Hello")
	numDn++
	dnList = append(dnList, "Hello1")
	numDn++
	dnList = append(dnList, "Hello2")
	numDn++
	dnList = append(dnList, "Hello3")
	numDn++

	http.HandleFunc("/create_file", createFile)
	http.HandleFunc("/get_file", getFile)
	http.HandleFunc("/block_report", blockReport)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

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
