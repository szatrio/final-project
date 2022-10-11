package server

import (
	// "final/server/controllers"

	"final/server/services"

	"github.com/gin-gonic/gin"
)

type Router struct {
	control *services.HandlersController
}

func UserRouther(control *services.HandlersController) *Router {
	return &Router{control: control}
}

func (r *Router) Start(port string) {
	router := gin.Default()

	///// User Handlers /////////
	router.POST("/users/register", r.control.Register_User)
	router.POST("/users/login", r.control.Login_User)
	router.PUT("/users", r.control.PUT_User)
	router.DELETE("/users", r.control.Delete_User)

	router.POST("/photos", r.control.Post_Photos)

	router.Run(port)
}
