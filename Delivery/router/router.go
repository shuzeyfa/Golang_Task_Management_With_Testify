package router

import (
	"log"
	"taskmanagement/Delivery/controllers"
	infrastructure "taskmanagement/Infrastructure"
	repository "taskmanagement/Repository"
	usecase "taskmanagement/Usecase"

	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	infrastructure.ConnectDB()
	log.Println("App is ready!")

	//initialize repository, usecase and controller
	taskRepo := &repository.MongoTaskRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("task"),
	}

	taskUsecase := &usecase.TaskUsecase{
		Repo: taskRepo,
	}

	taskController := &controllers.TaskController{
		Control: *taskUsecase,
	}

	r := gin.Default()

	r.POST("/register", controllers.RegisterHandler)
	r.POST("/login", controllers.LoginUser)

	r.GET("/tasks", infrastructure.AuthMiddleware(), controllers.GetAllTask)
	r.GET("/task/:id", infrastructure.AuthMiddleware(), controllers.GetTaskByID)
	r.POST("/create", infrastructure.AuthMiddleware(), taskController.CreateTask)
	r.PUT("/update/:id", infrastructure.AuthMiddleware(), controllers.UpdateTask)
	r.DELETE("/delete/:id", infrastructure.AuthMiddleware(), controllers.DeleteTask)

	return r
}
