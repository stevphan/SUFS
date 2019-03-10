package shared

import (
	"fmt"
	"net/http"
)

func StoreSingleBlock(storeBlockReq StoreBlockRequest) bool {
	for _, dataNodeIp := range storeBlockReq.DnList {
		success := StoreSingleToDataNode(storeBlockReq, dataNodeIp)
		if success {
			return true
		}
	}

	return false
}

func StoreSingleToDataNode(storeBlockReq StoreBlockRequest, dataNodeIp string) bool {
	VerbosePrintln(fmt.Sprintf("Attempting to save block to data node '%s'", dataNodeIp))

	buffer, err := ConvertObjectToJsonBuffer(storeBlockReq)
	if err != nil {
		VerbosePrintln(fmt.Sprint("Error while communicating to the data node:", err))
		return false
	}

	url := "http://" + dataNodeIp + PathBlock
	req, err := http.NewRequest(http.MethodPut, url, buffer)
	if err != nil {
		VerbosePrintln(fmt.Sprintln("Error creating request to data node:", err))
		return false
	}

	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}
	res, err := client.Do(req)
	if err != nil {
		VerbosePrintln(fmt.Sprintln("Error while communicating with the data node:", err))
		return false
	}

	storeBlockRes := StoreBlockResponse{}
	err = ObjectFromResponse(res, &storeBlockRes)
	CheckErrorAndFatal("Unable to parse response", err)

	if storeBlockRes.Err != "" {
		VerbosePrintln(fmt.Sprint("Error from data node:", storeBlockRes.Err))
		return false
	}

	return true
}
