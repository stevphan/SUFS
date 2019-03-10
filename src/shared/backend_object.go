package shared

// Create

type CreateFileNameNodeRequest struct {
	FileName string `json:"FileName"`
	Size     string `json:"Size"`
}

type CreateFileNameNodeResponse struct {
	BlockInfos []BlockInfo `json:"BlockInfos"`
	Err        string      `json:"Error"`
}

type StoreBlockRequest struct {
	Block   string   `json:"Block"` // base64 encoded block data
	DnList  []string `json:"DataNodeList"`
	BlockId string   `json:"BlockId"`
}

type StoreBlockResponse struct {
	Err string `json:"Error"`
}

// Get

type GetFileNameNodeRequest struct {
	FileName string `json:"FileName"`
}

type GetFileNameNodeResponse struct {
	BlockInfos []BlockInfo `json:"BlockInfos"`
	Err        string      `json:"Error"`
}

type GetBlockRequest struct {
	BlockId string `json:"BlockId"`
}

type GetBlockResponse struct {
	Block string `json:"Block"` // base64 encoded block data
	Err   string `json:"Error"`
}

//Block Report

type BlockReportRequest struct {
	MyIp     string   `json:"MyIp"`
	BlockIds []string `json:"BlockIds"`
}

type BlockReportResponse struct {
	Err string `json:"Error"`
}

//Heartbeat

type HeartbeatRequest struct {
	MyIp string `json:"MyIp"`
}

type HeartbeatResponse struct {
	Err string `json:"Error"`
}

//Replication Check

type ReplicationRequest struct {
	BlockId string   `json:"BlockId"`
	DnList  []string `json:"DataNodeList"`
}

type ReplicationResponse struct {
	Err string `json:"Error"`
}

// Helpers

type BlockInfo struct {
	BlockId string   `json:"BlockId"`
	DnList  []string `json:"DataNodeList"`
}

func MakeStoreBlockRequest(block string, blockInfo BlockInfo) StoreBlockRequest {
	return StoreBlockRequest{
		Block:   block,
		DnList:  blockInfo.DnList,
		BlockId: blockInfo.BlockId,
	}
}
