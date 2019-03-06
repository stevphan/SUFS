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
