package main

import "time"

type file struct {
	NumFiles int
	MetaData []fileMetaData
}

type fileMetaData struct {
	FileName   string
	NumBlocks  int
	BlockLists []blockList
}

type blockList struct {
	DnList []string //stores IP of DNs it is stored at
}

type dataNodeList struct {
	dnIP 	string
	dnTime time.Time
}

/*type createRequest struct{
	Filename string `json:"FileName"`
	Size string		`json:"Size"`
}

type createResponse struct {
	BlockId string	`json:"BlockId"`
	DnList []string `json:"DnList"`
}

type responseObj struct { //What is returned to cli
	Blocks []createResponse
}*/

/*type BlockReportRequest struct {
	MyIp string			`json:"MyIp"`
	BlockId []string 	`json:"BlockId"`
}*/
