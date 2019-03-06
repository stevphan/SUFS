package shared

import (
	"net/http"
	"fmt"
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

	dataNodeUrl := "http://" + dataNodeIp + "/storeBlock"
	buffer, err := ConvertObjectToJsonBuffer(storeBlockReq)
	if err != nil {
		VerbosePrintln(fmt.Sprint("Error while communicating to the data node:", err))
		return false
	}

	res, err := http.Post(dataNodeUrl, "application/json", buffer)
	if err != nil {
		VerbosePrintln(fmt.Sprint("Error while communicating to the data node:", err))
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