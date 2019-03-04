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
	http.HandleFunc("/storeBlock", store_and_foward)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}

}

