package mnet

import (
	"errors"
	"fmt"
	"log"
	"markV3/mface"
	"sync"
)

/*
accept route register
handle message
accept route response
*/
func newRouteManager() mface.MRouteManager {
	rm := &routeManager{
		status:       Serve_Status_UnStarted,
		requestChan:  make(chan mface.MMessage, 0),
		responseChan: make(chan mface.MMessage, 0),
		routes:       make(map[string]mface.MRouteHandler),
		routesLock:   sync.RWMutex{},
	}
	return rm
}

type routeManager struct {
	status int
	server mface.MServer

	requestChan  chan mface.MMessage
	responseChan chan mface.MMessage

	routes     map[string]mface.MRouteHandler
	routesLock sync.RWMutex
}

func (rm *routeManager) SetServer(s mface.MServer) {
	rm.server = s
}

func (rm *routeManager) Load() error {
	log.Printf("[RouteManager] Load")
	rm.status = Serve_Status_Load

	rm.requestChan = make(chan mface.MMessage, rm.server.Config().RMRequestCS())
	rm.responseChan = make(chan mface.MMessage, rm.server.Config().RMResponseCS())

	return nil
}

func (rm *routeManager) Start() error {
	log.Printf("[RouteManager] Start")
	rm.status = Serve_Status_Starting

	go rm.startAcceptRequestFromMM()
	go rm.startAcceptResponseToMM()

	rm.status = Serve_Status_Running

	return nil
}

func (rm *routeManager) Stop() error {
	log.Printf("[RouteManager] Stop")
	return nil
}

func (rm *routeManager) StartEnding() error {
	log.Printf("[RouteManager] Start Ending")
	rm.status = Serve_Status_Ending

	// 1.close request channel
	close(rm.requestChan)

	// 2.wait for buffer empty
	for {
		if len(rm.requestChan) == 0 {
			break
		}
	}

	return nil
}

func (rm *routeManager) OfficialEnding() error {
	log.Printf("[RouteManager] Official Ending")

	// 1.close response channel
	close(rm.responseChan)

	// 2. wait for buffer empty
	for {
		if len(rm.responseChan) == 0 {
			break
		}
	}

	// todo:3.clean all route functions

	rm.status = Serve_Status_Stopped

	return nil
}

func (rm *routeManager) Reload() error {
	log.Printf("[RouteManager] Reload")
	rm.status = Serve_Status_Reload
	rm.status = Serve_Status_Running
	return nil
}

func (rm *routeManager) AddRouteHandle(route mface.MRouteHandler) error {
	rm.routesLock.Lock()
	defer rm.routesLock.Unlock()

	if _, ok := rm.routes[route.RouteID()]; ok {
		return errors.New(fmt.Sprintf("[%s] routeId exists", route.RouteID()))
	}

	rm.routes[route.RouteID()] = route
	return nil
}

// ========== private functions =============

// request
func (rm *routeManager) startAcceptRequestFromMM() {
	for {
		request , ok := <- rm.requestChan
		if !ok {
			break
		}

		rm.handleRequest(request)
	}
}

// response
func (rm *routeManager) startAcceptResponseToMM() {

}

func (rm *routeManager) handleRequest(request mface.MMessage) {
	// 1.get handle route
	route , err := rm.getHandleRoute(request.MsgID())
	if err != nil {
		log.Printf("get handle route error : %+v" , err)
		return
	}

	// 2.goroutine handle request
	routeHandleFunc := route.RouteHandleFunc()
	go func(routeId string , handleFunc func(request mface.MMessage , response mface.MMessage) error) {
		var response mface.MMessage
		err := handleFunc(request , response)
		if err != nil {
			log.Printf("[%s] handle request error : %+v \n %+v" , routeId, request , err)
		}
		rm.responseChan <- response
	}(route.RouteID() , routeHandleFunc)
}

func (rm *routeManager) getHandleRoute(routeId string) (mface.MRouteHandler , error) {
	rm.routesLock.RLock()
	defer rm.routesLock.RUnlock()

	route, exists := rm.routes[routeId]
	if !exists {
		return nil , errors.New(fmt.Sprintf("[%s] routeId not exists" , routeId))
	}

	return route , nil
}
