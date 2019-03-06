package main

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
	DnList [repFact]string //stores IP of DNs it is stored at
}

type createRequest struct{
	Filename string `json:"FileName"`
	Size string		`json:"Size"`
}

type createResponse struct {
	BlockId string	`json:"BlockId"`
	DnList []string `json:"DnList"`
}

type responseObj struct { //What is returned to cli
	Blocks []createResponse
}