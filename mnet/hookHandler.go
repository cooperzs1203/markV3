/**
* @Author: Cooper
* @Date: 2019/12/4 18:13
 */

package mnet

import "markV3/mface"

const (
	Hook_Method_Request  = "Request"
	Hook_Method_Response = "Response"
)

func newHookHandler(method string, hookFunc func(mface.MMessage) bool) mface.MHookHandler {
	hh := &hookHandler{
		method:   method,
		hookFunc: hookFunc,
	}
	return hh
}

type hookHandler struct {
	method   string
	hookFunc func(mface.MMessage) bool
}

func (hh *hookHandler) Method() string {
	return hh.method
}

func (hh *hookHandler) HookFunction() func(mface.MMessage) bool {
	return hh.hookFunc
}
