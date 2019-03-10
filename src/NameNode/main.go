package main

import (
	"log"
	"net/http"
)

const (
	repFact = 3
	saveData = "fileData.json"
)

var (
	dnList = []dataNodeList{}
	//dnList = []string{}
	//numDn = 0
	files file
)
/*
block report send port as well
possibly use as argument during init (i.e. when booting main look for IP and port
 */

func main() {
	readFilesFromDisk()
	//for testing
	/*addToDnList("localhost:8081")
	addToDnList("localhost:8082")
	addToDnList("fake.ip.1")
	addToDnList("fake.ip.2")*/

	go replicationCheck()
	go dataNodeDead()

	log.Println("NameNode Running")
	http.HandleFunc("/createFile", createFile)
	http.HandleFunc("/getFile", getFile)
	http.HandleFunc("/blockReport", blockReport)
	http.HandleFunc("/heartBeat", heartBeat)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
