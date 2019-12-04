package mnet

import (
	"errors"
	"fmt"
	"log"
	"markV3/mface"
	"net"
	"sync"
)

/*
manage all of connection
accept request message from connections , send it to msgManager
accept response message from msgManager , send it to connections
*/
func newConnManager() mface.MConnManager {
	cm := &connManager{
		status:          Serve_Status_UnStarted,
		server:          nil,
		requestChan:     make(chan mface.MMessage, 0),
		responseChan:    make(chan mface.MMessage, 0),
		connections:     make(map[string]mface.MConnection),
		connectionsLock: sync.RWMutex{},
	}
	return cm
}

type connManager struct {
	status int
	server mface.MServer

	requestChan  chan mface.MMessage
	responseChan chan mface.MMessage

	connections     map[string]mface.MConnection
	connectionsLock sync.RWMutex
}

func (cm *connManager) SetServer(s mface.MServer) {
	cm.server = s
}

func (cm *connManager) Server() mface.MServer {
	return cm.server
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

func (cm *connManager) RequestChan() chan mface.MMessage {
	return cm.requestChan
}

func (cm *connManager) ResponseChan() chan mface.MMessage {
	return cm.responseChan
}

func (cm *connManager) HandleNewConn(conn *net.Conn) error {
	// 1.scan max connection count
	if cm.lenOfConnections() >= cm.server.Config().CMMaxConnNumber() {
		_, _ = (*conn).Write([]byte("Server connect max"))
		_ = (*conn).Close()
		return errors.New("ConnManager got max connection count")
	}

	// 2.get new connection
	newConnection, err := newConnection(conn , cm)
	if err != nil {
		return err
	}

	// 3.add new connection into manage
	if err := cm.addNewConnection(newConnection); err != nil {
		return err
	}

	// 4.start new connection
	if err := newConnection.Start(); err != nil {
		return err
	}

	return nil
}

// ======== private functions ============

// request
func (cm *connManager) startAcceptRequestFromConn() {
	for {
		request, ok := <-cm.requestChan
		if !ok {
			if cm.status >= Serve_Status_Ending {
				break
			} else {
				continue
			}
		}

		cm.server.MsgManager().RequestChan() <- request
	}
}

// response
func (cm *connManager) startAcceptResponseToConn() {
	for {
		response, ok := <-cm.responseChan
		if !ok {
			if cm.status >= Serve_Status_Ending {
				break
			} else {
				continue
			}
		}

		if err := cm.getConnAndReplyResponse(response); err != nil {
			log.Printf("[%s] get reply response error : %+v", response.ConnID(), err)
		}
	}
}

func (cm *connManager) addNewConnection(conn mface.MConnection) error {
	cm.connectionsLock.Lock()
	defer cm.connectionsLock.Unlock()

	if _, exists := cm.connections[conn.ID()]; exists {
		return errors.New(fmt.Sprintf("[%s] exists", conn.ID()))
	}

	cm.connections[conn.ID()] = conn
	return nil
}

func (cm *connManager) getConnection(connId string) (mface.MConnection, error) {
	cm.connectionsLock.RLock()
	defer cm.connectionsLock.RUnlock()

	conn, exists := cm.connections[connId]
	if !exists {
		return nil, errors.New(fmt.Sprintf("[%s] connection not exists", connId))
	}

	return conn, nil
}

func (cm *connManager) lenOfConnections() uint64 {
	cm.connectionsLock.RLock()
	defer cm.connectionsLock.RUnlock()

	return uint64(len(cm.connections))
}

func (cm *connManager) getConnAndReplyResponse(response mface.MMessage) error {
	// 1.get conn
	conn, err := cm.getConnection(response.ConnID())
	if err != nil {
		return err
	}

	// 2. send it to conn response channel
	if conn.Status() < Serve_Status_Ending {
		conn.ResponseChan() <- response
	}

	return nil
}
