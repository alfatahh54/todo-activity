package routes

import "github.com/gin-gonic/gin"

type MainRouteType struct {
	Engine *gin.Engine
}

var MainRoute MainRouteType

func Router(r *gin.Engine) {
	MainRoute.Engine = r
}

func (r *MainRouteType) NewRoute(method string, endpoint string, handlers ...gin.HandlerFunc) {
	for {
		if r.Engine != nil {
			break
		}
	}
	r.Engine.Handle(method, endpoint, handlers...)
}
