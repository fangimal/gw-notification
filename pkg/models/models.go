package models

import "time"

type Notification struct {
	UserID    int64     `json:"user_id" bson:"user_id"`
	Amount    float32   `json:"amount" bson:"amount"`
	Currency  string    `json:"currency" bson:"currency"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}
