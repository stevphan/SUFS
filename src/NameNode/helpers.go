package main

import (
	"bytes"
	"encoding/json"
)

func convertObjectToJsonBuffer(object interface{}) (*bytes.Buffer, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	buffer := bytes.NewBuffer(data)

	return buffer, nil
}

func convertObjectToJson(object interface{}) ([]byte, error) {
	data, err := json.Marshal(object)
	if err != nil {
		return nil, err
	}

	return data, nil
}

/*func writeFilesToDisk() {
	js, _ := convertObjectToJson(files)
	err := ioutil.WriteFile(saveData, js, 0644)
	fmt.Println(err)
}

func readFilesFromDisk() {
	tempFile, err := os.Open(saveData)
	if err != nil {
		fmt.Println(err.Error())
	}
	decoder := json.NewDecoder(tempFile)
	myFile := file{}
	decoder.Decode(&myFile)
	files = myFile
	fmt.Println(files)
	fmt.Println()
}*/