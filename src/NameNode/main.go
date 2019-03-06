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

func main() {
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

	http.HandleFunc("/createFile", createFile)
	http.HandleFunc("/getFile", getFile)
	http.HandleFunc("/blockReport", blockReport)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
