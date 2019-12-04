package mface

type MRouteHandler interface {
	RouteID() string
	RouteHandleFunc() func(MMessage , MMessage) error
}

type MRouteManager interface {
	MBase
	SetServer(MServer)
	AddRouteHandle(route MRouteHandler) error
}
