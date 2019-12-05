package mface

type RouteHandleFunc func(MMessage) MMessage

type MRouteHandler interface {
	RouteID() string
	RouteHandleFunc() RouteHandleFunc
}

type MRouteManager interface {
	MBase
	SetServer(MServer)
	AddRouteHandle(route MRouteHandler) error
	RequestChan() chan MMessage
	ResponseChan() chan MMessage
}
