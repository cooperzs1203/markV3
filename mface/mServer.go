package mface

type MServer interface {
	MBase

	Config() MConfig

	ConnManager() MConnManager
	MsgManager() MMsgManager
	RouteManager() MRouteManager
	RunEntranceFunc(func() error)

	AddRoute(string, RouteHandleFunc) error
	AddRoutes(map[string]RouteHandleFunc) error

	AddRequestHook(func(MMessage) bool)
	AddResponseHook(func(MMessage) bool)
}
