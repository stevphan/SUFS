package main

type createFileNameNodeRequest struct {
	FileName string `json:"FileName"`
	Size     string `json:"Size"`
}

type createFileNameNodeResponse struct {
	BlockInfos []blockInfo `json:"BlockInfos"`
	Err        string      `json:"Error"`
}

type blockInfo struct {
	BlockId string   `json:"BlockId"`
	DnList  []string `json:"DataNodeList"`
}

type storeBlockRequest struct {
	Block   string   `json:"Block"` // base64 encoded block data
	DnList  []string `json:"DataNodeList"`
	BlockId string   `json:"BlockId"`
}

type storeBlockResponse struct {
	Err string `json:"Error"`
}

func makeStoreBlockRequest(block string, blockInfo blockInfo) storeBlockRequest {
	return storeBlockRequest{
		Block:   block,
		DnList:  blockInfo.DnList,
		BlockId: blockInfo.BlockId,
	}
}
