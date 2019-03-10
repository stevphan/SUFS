package main

import (
	"log"
	"net/http"
	"shared"
	"time"
)

const (
	heartTimer = 10 //in seconds
)


func send_heartbeat() { // sends heartbeat every heartTimer amount of seconds, so every 10 seconds
	ticker := time.NewTicker(heartTimer * time.Second)
	for range ticker.C {
		log.Println("sending heartbeat")
		heartbeat()
	}
}

func heartbeat() {
	nameNodeUrl := "http://" + nameNodeAddress + shared.PathHeartbeat
	heartbeatReq := shared.HeartbeatRequest{}
	heartbeatReq.MyIp = selfAddress

	buffer, err := shared.ConvertObjectToJsonBuffer(heartbeatReq)

	client := http.Client{}
	req, err := http.NewRequest(http.MethodPut, nameNodeUrl, buffer)

	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		log.Println("ERROR: ", err)
		return
	}


	if err != nil {
		log.Println("Error sending heartbeat:", err)
		return
	}
	// should I even keep response if heartbeat is one-way anyway?
	heartbeatResp := shared.HeartbeatResponse{}
	err = shared.ObjectFromResponse(res, &heartbeatResp)
	if err != nil {
		log.Println("Error parsing heartbeat response:", err)
		return
	}
	if heartbeatResp.Err != "" {
		log.Println("Error in heartbeat response:", heartbeatResp.Err)
		return
	}

	return
}