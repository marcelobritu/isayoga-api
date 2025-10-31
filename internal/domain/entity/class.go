package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Class struct {
	ID               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title            string             `json:"title" bson:"title"`
	Description      string             `json:"description" bson:"description"`
	InstructorID     primitive.ObjectID `json:"instructor_id" bson:"instructor_id"`
	InstructorName   string             `json:"instructor_name" bson:"instructor_name"`
	StartTime        time.Time          `json:"start_time" bson:"start_time"`
	EndTime          time.Time          `json:"end_time" bson:"end_time"`
	MaxCapacity      int                `json:"max_capacity" bson:"max_capacity"`
	CurrentEnrolled  int                `json:"current_enrolled" bson:"current_enrolled"`
	PriceInCents     int64              `json:"price_in_cents" bson:"price_in_cents"`
	Status           string             `json:"status" bson:"status"`
	CreatedAt        time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at" bson:"updated_at"`
	Version          int                `json:"version" bson:"version"`
}

func NewClass(title, description string, instructorID primitive.ObjectID, instructorName string, startTime, endTime time.Time, maxCapacity int, priceInCents int64) *Class {
	now := time.Now()
	return &Class{
		ID:              primitive.NewObjectID(),
		Title:           title,
		Description:     description,
		InstructorID:    instructorID,
		InstructorName:  instructorName,
		StartTime:       startTime,
		EndTime:         endTime,
		MaxCapacity:     maxCapacity,
		CurrentEnrolled: 0,
		PriceInCents:    priceInCents,
		Status:          "active",
		CreatedAt:       now,
		UpdatedAt:       now,
		Version:         0,
	}
}

func (c *Class) HasAvailableSpots() bool {
	return c.CurrentEnrolled < c.MaxCapacity
}

func (c *Class) IncrementEnrollment() {
	c.CurrentEnrolled++
	c.UpdatedAt = time.Now()
	c.Version++
}

func (c *Class) DecrementEnrollment() {
	if c.CurrentEnrolled > 0 {
		c.CurrentEnrolled--
		c.UpdatedAt = time.Now()
		c.Version++
	}
}

