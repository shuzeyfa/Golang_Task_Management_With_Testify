package repository

import (
	"context"
	"errors"
	domain "taskmanagement/Domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoTaskRepository struct {
	Collection *mongo.Collection
}

func (r *MongoTaskRepository) GetAllTask(userID primitive.ObjectID) ([]domain.Task, error) {

	filter := bson.M{"user_id": userID}

	cursor, err := r.Collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var tasks []domain.Task
	if err := cursor.All(context.Background(), &tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *MongoTaskRepository) GetTaskByID(taskID primitive.ObjectID, userID primitive.ObjectID) (domain.Task, error) {

	var task domain.Task
	err := r.Collection.FindOne(context.Background(), bson.M{"_id": taskID, "user_id": userID}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return domain.Task{}, errors.New("Task not found")
		}

		return domain.Task{}, errors.New("Server error")

	}

	return task, nil
}

func (r *MongoTaskRepository) CreateTask(task domain.Task, userID primitive.ObjectID) (domain.Task, bool) {

	task.UserId = userID
	_, err := r.Collection.InsertOne(context.Background(), task)
	if err != nil {
		return domain.Task{}, false
	}

	return task, true
}

func (r *MongoTaskRepository) UpdateTask(task domain.Task, userID primitive.ObjectID) (domain.Task, bool) {

	filter := bson.M{"_id": task.ID, "user_id": userID}
	err := r.Collection.FindOneAndUpdate(context.Background(), filter, bson.M{"$set": task}).Decode(&task)
	if err != nil {
		return domain.Task{}, false
	}

	return task, true
}

func (r *MongoTaskRepository) DeleteTask(taskID primitive.ObjectID, userID primitive.ObjectID) error {

	filter := bson.M{"_id": taskID, "user_id": userID}
	result, err := r.Collection.DeleteOne(context.Background(), filter)
	if err != nil || result.DeletedCount == 0 {
		return err
	}

	return nil
}
