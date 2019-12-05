/**
* @Author: Cooper
* @Date: 2019/12/4 14:55
 */

package mnet

import "markV3/mface"

func newRouteHandler(routeId string, routeHandleFunc mface.RouteHandleFunc) mface.MRouteHandler {
	rh := &routeHandler{
		routeId:         routeId,
		routeHandleFunc: routeHandleFunc,
	}
	return rh
}

type routeHandler struct {
	routeId         string
	routeHandleFunc mface.RouteHandleFunc
}

func (rh *routeHandler) RouteID() string {
	return rh.routeId
}

func (rh *routeHandler) RouteHandleFunc() mface.RouteHandleFunc {
	return rh.routeHandleFunc
}
