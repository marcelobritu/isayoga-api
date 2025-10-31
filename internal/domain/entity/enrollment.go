package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Enrollment struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID       primitive.ObjectID `json:"user_id" bson:"user_id"`
	ClassID      primitive.ObjectID `json:"class_id" bson:"class_id"`
	PaymentID    string             `json:"payment_id" bson:"payment_id"`
	Status       string             `json:"status" bson:"status"`
	EnrolledAt   time.Time          `json:"enrolled_at" bson:"enrolled_at"`
	CancelledAt  *time.Time         `json:"cancelled_at,omitempty" bson:"cancelled_at,omitempty"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
}

func NewEnrollment(userID, classID primitive.ObjectID) *Enrollment {
	now := time.Now()
	return &Enrollment{
		ID:        primitive.NewObjectID(),
		UserID:    userID,
		ClassID:   classID,
		Status:    "pending",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (e *Enrollment) Confirm(paymentID string) {
	e.Status = "confirmed"
	e.PaymentID = paymentID
	e.EnrolledAt = time.Now()
	e.UpdatedAt = time.Now()
}

func (e *Enrollment) Cancel() {
	e.Status = "cancelled"
	now := time.Now()
	e.CancelledAt = &now
	e.UpdatedAt = now
}

