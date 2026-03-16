package usecase

import (
	"errors"
	domain "taskmanagement/Domain"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskUsecase struct {
	Repo domain.TaskRepository
}

func (u *TaskUsecase) GetAllTask(userId primitive.ObjectID) ([]domain.Task, error) {

	tasks, err := u.Repo.GetAllTask(userId)
	if err != nil {
		return []domain.Task{}, err
	}

	return tasks, nil
}

func (u *TaskUsecase) GetTaskByID(taskId primitive.ObjectID, userId primitive.ObjectID) (domain.Task, error) {

	task, err := u.Repo.GetTaskByID(taskId, userId)
	if err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (u *TaskUsecase) CreateTask(task domain.Task, userId primitive.ObjectID) (domain.Task, error) {
	createdTask, ok := u.Repo.CreateTask(task, userId)
	if !ok {
		return domain.Task{}, errors.New("could not create task")
	}
	return createdTask, nil
}

func (u *TaskUsecase) UpdateTask(task domain.Task, userId primitive.ObjectID) (domain.Task, error) {
	updatedTask, err := u.Repo.UpdateTask(task, userId)
	if err != nil {
		return domain.Task{}, errors.New("could not update task")
	}
	return updatedTask, nil
}

func (u *TaskUsecase) DeleteTask(taskId primitive.ObjectID, userId primitive.ObjectID) error {
	err := u.Repo.DeleteTask(taskId, userId)
	if err != nil {
		return errors.New("could not delete task")
	}
	return nil
}
