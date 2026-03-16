package controllers

import (
	"net/http"
	domain "taskmanagement/Domain"
	usecase "taskmanagement/Usecase"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskController struct {
	Control *usecase.TaskUsecase
}

func RegisterHandler(c *gin.Context) {

	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	user, err := usecase.RegisterUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, user)
}

func LoginUser(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tokenString, err := usecase.LoginUser(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}

func getUserIDFromContext(c *gin.Context) (primitive.ObjectID, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists {
		return primitive.NilObjectID, false
	}

	userID, ok := userIDInterface.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, false
	}

	return userID, true
}

func GetTaskIdFromContext(c *gin.Context) (primitive.ObjectID, bool) {
	TaskIDInterface, exists := c.Get("task_id")
	if !exists {
		return primitive.NilObjectID, false
	}

	userID, ok := TaskIDInterface.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, false
	}

	return userID, true
}

func (ctr *TaskController) GetAllTask(c *gin.Context) {

	userId, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}

	tasks, err := ctr.Control.GetAllTask(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (ctr *TaskController) GetTaskByID(c *gin.Context) {
	TaskId, ok1 := GetTaskIdFromContext(c)
	userId, ok2 := getUserIDFromContext(c)
	if !ok1 || !ok2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}

	task, err := ctr.Control.GetTaskByID(TaskId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch task"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (ctr *TaskController) CreateTask(c *gin.Context) {
	var task domain.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	result, err := ctr.Control.CreateTask(task, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not create task"})
		return
	}

	c.JSON(http.StatusCreated, result)

}

func (ctr *TaskController) UpdateTask(c *gin.Context) {
	var task domain.Task

	if err := c.ShouldBindJSON(&task); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := getUserIDFromContext(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	result, err := ctr.Control.UpdateTask(task, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update task"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (ctr *TaskController) DeleteTask(c *gin.Context) {
	TaskId, ok1 := GetTaskIdFromContext(c)
	userId, ok2 := getUserIDFromContext(c)
	if !ok1 || !ok2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}

	err := ctr.Control.DeleteTask(TaskId, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete task"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted successfully"})
}
