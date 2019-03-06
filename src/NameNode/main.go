package main

import (
	"net/http"
)

const (
	repFact = 3
	saveData = "fileData.json"
)

var (
	dnList = []string{}
	numDn = 0
	files file
)
/*
block report send port as well
possibly use as argument during init (i.e. when booting main look for IP and port
 */

func main() {
	readFilesFromDisk()
	//for testing
	dnList = append(dnList, "localhost:8081")
	numDn++
	dnList = append(dnList, "localhost:8082")
	numDn++

	http.HandleFunc("/createFile", createFile)
	http.HandleFunc("/getFile", getFile)
	http.HandleFunc("/blockReport", blockReport)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
