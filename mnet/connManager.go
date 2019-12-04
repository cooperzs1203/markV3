package mnet

import (
	"log"
	"markV3/mface"
)

/*
manage all of connection
accept request message from connections , send it to msgManager
accept response message from msgManager , send it to connections
*/
func newConnManager() mface.MConnManager {
	cm := &connManager{
		status:       Serve_Status_UnStarted,
		server:       nil,
		requestChan:  make(chan mface.MMessage, 0),
		responseChan: make(chan mface.MMessage, 0),
	}
	return cm
}

type connManager struct {
	status int
	server mface.MServer

	requestChan  chan mface.MMessage
	responseChan chan mface.MMessage
}

func (cm *connManager) SetServer(s mface.MServer) {
	cm.server = s
}

func (cm *connManager) Load() error {
	log.Printf("[ConnManager] Load")
	cm.status = Serve_Status_Load

	cm.requestChan = make(chan mface.MMessage, cm.server.Config().CMRequestCS())
	cm.responseChan = make(chan mface.MMessage, cm.server.Config().CMResponseCS())

	return nil
}

func (cm *connManager) Start() error {
	log.Printf("[ConnManager] Start")
	cm.status = Serve_Status_Starting

	go cm.startAcceptRequestFromConn()
	go cm.startAcceptResponseToConn()

	cm.status = Serve_Status_Running

	return nil
}

func (cm *connManager) Stop() error {
	log.Printf("[ConnManager] Stop")

	return nil
}

func (cm *connManager) StartEnding() error {
	log.Printf("[ConnManager] Start Ending")
	cm.status = Serve_Status_Ending

	// todo:1.close all connections of request goroutine

	// 2.close request channel
	close(cm.requestChan)

	// 3.wait for buffer empty
	for {
		if len(cm.requestChan) == 0 {
			break
		}
	}

	return nil
}

func (cm *connManager) OfficialEnding() error {
	log.Printf("[ConnManager] Official Ending")

	// 1.close response channel
	close(cm.responseChan)

	// 2. wait for buffer empty
	for {
		if len(cm.responseChan) == 0 {
			break
		}
	}

	// todo:3.close all connections of response goroutine

	// todo:4.stop all of connection , and clean []connections

	cm.status = Serve_Status_Stopped

	return nil
}

func (cm *connManager) Reload() error {
	log.Printf("[ConnManager] Reload")
	cm.status = Serve_Status_Reload

	// 1. ready new request and response channel
	newRequestChan := make(chan mface.MMessage, cm.server.Config().CMRequestCS())
	newResponseChan := make(chan mface.MMessage, cm.server.Config().CMResponseCS())

	// 2. close request and response channel 
	close(cm.requestChan)
	close(cm.responseChan)

	// 3. wait for request and response channel empty
	for {
		if len(cm.requestChan) == 0 {
			break
		}
	}

	for {
		if len(cm.responseChan) == 0 {
			break
		}
	}

	// 4.exchange request and response channel
	cm.requestChan = newRequestChan
	cm.responseChan = newResponseChan

	cm.status = Serve_Status_Running

	return nil
}

// ======== private functions ============

// request
func (cm *connManager) startAcceptRequestFromConn() {

}

// response
func (cm *connManager) startAcceptResponseToConn() {

}
