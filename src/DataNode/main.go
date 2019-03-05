package main

/*
DataNodeList
BlockId
 */

import (
	"log"
	"net/http"
)

type address struct {
	s3address string
}

const s3address string = "/Users/stxv/blocks/"

func main() {
	http.HandleFunc("/getBlock", get_block)
	http.HandleFunc("/storeBlock", store_and_foward)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic(err)
	}

}

