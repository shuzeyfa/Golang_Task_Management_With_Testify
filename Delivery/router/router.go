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

	//initialize repository, usecase and controller for task
	taskRepo := &repository.MongoTaskRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("task"),
	}

	taskUsecase := &usecase.TaskUsecase{
		Repo: taskRepo,
	}

	taskController := &controllers.TaskController{
		Control: taskUsecase,
	}

	userRepo := &repository.MongoUserRepository{
		Collection: infrastructure.Client.Database(infrastructure.DBName).Collection("user"),
	}

	UserUsecase := &usecase.UserUsecase{
		Repo: userRepo,
	}

	UserController := &controllers.UserController{
		Control: UserUsecase,
	}

	r := gin.Default()

	r.POST("/register", UserController.RegisterHandler)
	r.POST("/login", UserController.LoginUser)

	r.GET("/tasks", infrastructure.AuthMiddleware(), taskController.GetAllTask)
	r.GET("/task/:id", infrastructure.AuthMiddleware(), taskController.GetTaskByID)
	r.POST("/create", infrastructure.AuthMiddleware(), taskController.CreateTask)
	r.PUT("/update/:id", infrastructure.AuthMiddleware(), taskController.UpdateTask)
	r.DELETE("/delete/:id", infrastructure.AuthMiddleware(), taskController.DeleteTask)

	return r
}
