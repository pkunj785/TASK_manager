package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type TaskList struct {
	ID     primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	Task   string             `json:"task, omitempty"`
	Status bool               `json:"status, omitempty"`
}

type User struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
}
