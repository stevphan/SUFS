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
	nameNodeUrl := "http://" + nameNodeAddress + "/heartBeat"
	heartbeatReq := shared.HeartbeatRequest{}
	heartbeatReq.MyIp = selfAddress
	buffer, err := shared.ConvertObjectToJsonBuffer(heartbeatReq)
	res, err := http.Post(nameNodeUrl, "application/json", buffer)
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