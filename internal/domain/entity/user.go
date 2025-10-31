package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRole string

const (
	RoleStudent    UserRole = "student"
	RoleInstructor UserRole = "instructor"
	RoleAdmin      UserRole = "admin"
)

type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Role      UserRole           `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

func (u *User) IsStudent() bool {
	return u.Role == RoleStudent
}

func (u *User) IsInstructor() bool {
	return u.Role == RoleInstructor
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func NewUser(name, email string) *User {
	now := time.Now()
	return &User{
		ID:        primitive.NewObjectID(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (u *User) Update(name, email string) {
	u.Name = name
	u.Email = email
	u.UpdatedAt = time.Now()
}
