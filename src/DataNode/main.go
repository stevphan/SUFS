package main

/*
DataNodeList
BlockId
 */

import (
	"log"
	"net/http"
	"os"
)

var (
	s3address string
	directory = "/blocks/"
)

func main() {
	/*
	ensure directory exists
	 */
	addr := ""
	if len(os.Args) >= 1 {
		addr = os.Args[1]
	}

	if addr == "" {
		//ipAddress, _ :=  http.Get("http://169.254.169.254/latest/meta-data/public-ipv4)")
		//body, _ := ioutil.ReadAll(ipAddress.Body)
	}

	// accessible var so other store_and_forward can check if self in DnList
	s3address = addr

	http.HandleFunc("/getBlock", get_block)
	http.HandleFunc("/storeBlock", store_and_foward)

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panic(err)
	}

}

