package main

import (
	"encoding/base64"
	"fmt"
)

func do_not_use_encodeThenDecodeExample() {
	blockSize := 2
	fileData := []byte("hello_world")
	fmt.Printf("original: '%s' %v\n", string(fileData), fileData)

	blocks := []string{}
	byteIndex := 0

	for byteIndex < len(fileData) {
		bytesLeftCount := len(fileData) - byteIndex
		endIndex := byteIndex + min(blockSize, bytesLeftCount)

		base64Encoded := base64.StdEncoding.EncodeToString(fileData[byteIndex:endIndex])
		blocks = append(blocks, base64Encoded)

		byteIndex += blockSize
	}

	fmt.Println(blocks)

	for i, block := range blocks {
		decoded, _ := base64.StdEncoding.DecodeString(block)
		fmt.Println("decoded:", i, string(decoded))
	}
}
