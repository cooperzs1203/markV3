/**
* @Author: Cooper
* @Date: 2019/12/4 14:55
 */

package mnet

import "markV3/mface"

func newRouteHandler(routeId string, routeHandleFunc func(mface.MMessage, mface.MMessage) error) mface.MRouteHandler {
	rh := &routeHandler{
		routeId:         routeId,
		routeHandleFunc: routeHandleFunc,
	}
	return rh
}

type routeHandler struct {
	routeId string
	routeHandleFunc func(mface.MMessage , mface.MMessage) error
}

func (rh *routeHandler) RouteID() string {
	return rh.routeId
}

func (rh *routeHandler) RouteHandleFunc() func(mface.MMessage, mface.MMessage) error {
	return rh.routeHandleFunc
}