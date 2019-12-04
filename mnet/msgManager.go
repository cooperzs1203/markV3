package mnet

import (
	"log"
	"markV3/mface"
	"sync"
)

/*
filter message
*/
func newMsgManager() mface.MMsgManager {
	mm := &msgManager{
		status:       Serve_Status_UnStarted,
		requestChan:  make(chan mface.MMessage, 0),
		responseChan: make(chan mface.MMessage, 0),
		hooks:        make([]mface.MHookHandler, 0),
		hooksLock:    sync.RWMutex{},
	}
	return mm
}

type msgManager struct {
	status int
	server mface.MServer

	requestChan  chan mface.MMessage
	responseChan chan mface.MMessage

	hooks     []mface.MHookHandler
	hooksLock sync.RWMutex
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

	// 3.clean all hook functions
	mm.cleanHooks()

	mm.status = Serve_Status_Stopped

	return nil
}

func (mm *msgManager) Reload() error {
	log.Printf("[MsgManager] Reload")
	mm.status = Serve_Status_Reload
	// 1. ready new request and response channel
	newRequestChan := make(chan mface.MMessage, mm.server.Config().CMRequestCS())
	newResponseChan := make(chan mface.MMessage, mm.server.Config().CMResponseCS())

	// 2. close request and response channel
	close(mm.requestChan)
	close(mm.responseChan)

	// 3. wait for request and response channel empty
	for {
		if len(mm.requestChan) == 0 {
			break
		}
	}

	for {
		if len(mm.responseChan) == 0 {
			break
		}
	}

	// 4.exchange request and response channel
	mm.requestChan = newRequestChan
	mm.responseChan = newResponseChan

	mm.status = Serve_Status_Running
	return nil
}

func (mm *msgManager) RequestChan() chan mface.MMessage {
	return mm.requestChan
}

func (mm *msgManager) ResponseChan() chan mface.MMessage {
	return mm.responseChan
}

func (mm *msgManager) AddHook(hook mface.MHookHandler) {
	mm.hooksLock.Lock()
	defer mm.hooksLock.Unlock()

	mm.hooks = append(mm.hooks, hook)
}

// ======== private functions ============

// request
func (mm *msgManager) startAcceptRequestFromCM() {
	for {
		request, ok := <-mm.requestChan
		if !ok {
			if mm.status >= Serve_Status_Ending {
				break
			} else {
				continue
			}
		}

		if pass := mm.getHookAndHandleMessage(request, Hook_Method_Request); pass {
			mm.server.RouteManager().RequestChan() <- request
		} else {
			log.Printf("[%s] request message can't pass hook : %+v", request.MsgID(), request)
		}
	}
}

// response
func (mm *msgManager) startAcceptResponseToCM() {
	for {
		response, ok := <-mm.responseChan
		if !ok {
			if mm.status >= Serve_Status_Ending {
				break
			} else {
				continue
			}
		}

		if pass := mm.getHookAndHandleMessage(response, Hook_Method_Response); pass {
			mm.server.ConnManager().ResponseChan() <- response
		} else {
			log.Printf("[%s] response message can't pass hook : %+v", response.MsgID(), response)
		}
	}
}

func (mm *msgManager) getHookAndHandleMessage(msg mface.MMessage, method string) bool {
	mm.hooksLock.RLock()
	defer mm.hooksLock.RUnlock()

	pass := true
	for index := range mm.hooks {
		hook := mm.hooks[index]
		if hook.Method() != method {
			continue
		}
		hookFunc := hook.HookFunction()
		if hookFunc != nil {
			if pass = hookFunc(msg); !pass {
				break
			}
		}
	}

	return pass
}

func (mm *msgManager) cleanHooks() {
	mm.hooksLock.Lock()
	defer mm.hooksLock.Unlock()

	mm.hooks = mm.hooks[:0]
	mm.hooks = nil
}
