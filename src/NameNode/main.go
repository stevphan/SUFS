package main

import (
	"net/http"
)

const (
	RepFact = 3
)

type files struct {
	numFiles int
	metaData []fileMetaData
}

type fileMetaData struct {
	fileName string
	numBlocks int
	blockLists []blockList
}

type blockList struct {
	DnList [RepFact]string
}

var (
	DnList = []string{}
	NumDn int
)

func main() {
	//TODO read from disk, build the files from that
	http.HandleFunc("/create_file", createFile)
	http.HandleFunc("/get_file", getFile)
	http.HandleFunc("/block_report", blockReport)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
