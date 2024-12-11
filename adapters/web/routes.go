package web

import (

	"github.com/chud-lori/go-echo-temp/domain/ports"
	"github.com/labstack/echo/v4"
)
func UserRouter(controller ports.UserController, e *echo.Echo) {
	e.POST("/api/user", controller.Create)
    e.PUT("/api/user/:userId", controller.Update)
    e.DELETE("/api/user/:userId", controller.Delete)
    e.GET("/api/user/:userId", controller.FindById)
	e.GET("/api/user", controller.FindAll)

}
