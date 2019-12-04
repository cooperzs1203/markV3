package mnet

import (
	"log"
	"markV3/mface"
)

/*
filter message
*/
func newMsgManager() mface.MMsgManager {
	mm := &msgManager{
		status: Serve_Status_UnStarted,
	}
	return mm
}

type msgManager struct {
	status int
	server mface.MServer

	requestChan  chan mface.MMessage
	responseChan chan mface.MMessage
}

func (mm *msgManager) SetServer(s mface.MServer) {
	mm.server = s
}

func (mm *msgManager) Load() error {
	log.Printf("[MsgManager] Load")
	mm.status = Serve_Status_Load

	mm.requestChan = make(chan mface.MMessage, mm.server.Config().MMRequestCS())
	mm.responseChan = make(chan mface.MMessage, mm.server.Config().MMResponseCS())

	return nil
}

func (mm *msgManager) Start() error {
	log.Printf("[MsgManager] Start")

	mm.status = Serve_Status_Starting

	go mm.startAcceptRequestFromCM()
	go mm.startAcceptResponseToCM()

	mm.status = Serve_Status_Running

	return nil
}

func (mm *msgManager) Stop() error {
	log.Printf("[MsgManager] Stop")
	return nil
}

func (mm *msgManager) StartEnding() error {
	log.Printf("[MsgManager] Start Ending")
	mm.status = Serve_Status_Ending

	// 1.close request channel
	close(mm.requestChan)

	// 2.wait for buffer empty
	for {
		if len(mm.requestChan) == 0 {
			break
		}
	}

	return nil
}

func (mm *msgManager) OfficialEnding() error {
	log.Printf("[MsgManager] Official Ending")

	// 1.close response channel
	close(mm.responseChan)

	// 2. wait for buffer empty
	for {
		if len(mm.responseChan) == 0 {
			break
		}
	}

	// todo:3.clean all hook functions

	mm.status = Serve_Status_Stopped


	return nil
}

func (mm *msgManager) Reload() error {
	log.Printf("[MsgManager] Reload")
	mm.status = Serve_Status_Reload
	mm.status = Serve_Status_Running
	return nil
}

// ======== private functions ============

// request
func (mm *msgManager) startAcceptRequestFromCM() {

}

// response
func (mm *msgManager) startAcceptResponseToCM() {

}
