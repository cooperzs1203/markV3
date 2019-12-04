package mface

import "net"

type MConnManager interface {
	MBase

	SetServer(MServer)
	Server() MServer
	RequestChan() chan MMessage
	ResponseChan() chan MMessage

	HandleNewConn(*net.Conn) error
}
