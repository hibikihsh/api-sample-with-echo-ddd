package router

import (
	"api-sample-with-echo-ddd/interface/handler"

	"github.com/labstack/echo"
)

// InitRouting routesの初期化
func InitRouting(e *echo.Echo, userHandler handler.UserHandler) {
	e.POST("/user", userHandler.Post)
	e.GET("/user/:id", userHandler.Get)
	e.GET("/users", userHandler.GetAll)
	e.PUT("/user/:id", userHandler.Put)
	e.DELETE("/user/:id", userHandler.Delete)
}
