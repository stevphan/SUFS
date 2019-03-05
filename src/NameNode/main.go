package main

import (
	"net/http"
)

const (
	repFact = 3
	saveData = "fileData.json"
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
	readFilesFromDisk()
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
