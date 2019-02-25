package main

import (
)

const (
	blockSize int = 64 // assuming MB
)

func createfile(fileName string, fileSize int) { //needs to return list of dataNodes per block

	//Check the list of files if the supplied file name already exist

	blocksRequired := fileSize / blockSize


}
