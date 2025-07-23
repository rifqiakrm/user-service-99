package main

import (
	"user-service/db"
	"user-service/handler"
	"user-service/repository"
	"user-service/service"

	"github.com/gin-gonic/gin"
)

// main initializes dependencies and runs the HTTP server for the user service.
func main() {
	gormDB, err := db.InitDB("user.db")

	if err != nil {
		panic(err)
	}

	userRepo := repository.NewUserRepo(gormDB)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	r := gin.Default()

	r.GET("/users", userHandler.GetAllUsers)
	r.GET("/users/:id", userHandler.GetUser)
	r.POST("/users/batch", userHandler.BatchFetchUsers)
	r.POST("/users", userHandler.CreateUser)

	_ = r.Run(":6001")
}
