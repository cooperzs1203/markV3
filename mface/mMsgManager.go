package mface

type MHookHandler interface {
	Method() string
	HookFunction() func(MMessage) bool
}

type MMsgManager interface {
	MBase
	SetServer(MServer)
	RequestChan() chan MMessage
	ResponseChan() chan MMessage
	AddHook(MHookHandler)
}
