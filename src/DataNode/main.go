package main

/*
DataNodeList
BlockId
 */

import (
	"log"
	"net/http"
)

var verbose = false


func main() {
	http.HandleFunc("/getBlock", get_block)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}

}

