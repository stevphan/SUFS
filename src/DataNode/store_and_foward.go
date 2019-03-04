package main

import "net/http"

/*
Block
BlockId
DataNodeList
 */

 // passed in JSON payloads

type saveBlockRequest struct{
	Block string `json:"Block"`
	DataNodeList[] string `json:"DataNodeList"`
	BlockId string `json:"BlockId"`

}

func store_and_foward(write http.ResponseWriter, req *http.Request)  { // returns block requested from the current DN



	// TODO: changearguments to fit correct type
	// TODO: check for valid filename
	// TODO: check for valid save location

}
