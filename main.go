package main

import (
	"final/db"
	"final/server"
	"final/server/services"
)

// @title My Grams APP
// @description Final Project Golang
// @version v1.0
// @termsOfService http://swagger.io/terms/
// @BasePath /
// @host localhost:52000
// @contact.name Hans Parson
// @contact.email hansparson013@gmail.com

func main() {

	port := ":4000"
	db := db.ConnectGorm()

	serviceController := services.User_DB_Controller(db)
	user_service := server.UserRouther(serviceController)
	user_service.Start(port)
}
