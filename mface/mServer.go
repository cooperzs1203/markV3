package mface

type MServer interface {
	MBase

	Config() MConfig

	ConnManager() MConnManager
	MsgManager() MMsgManager
	RouteManager() MRouteManager
	RunEntranceFunc(func() error)

	AddRoute(string, func(MMessage, MMessage) error) error
	AddRoutes(map[string]func(MMessage, MMessage) error) error

	AddRequestHook(func(MMessage) bool)
	AddResponseHook(func(MMessage) bool)
}
