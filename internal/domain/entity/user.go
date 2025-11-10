package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserRole string

const (
	RoleStudent    UserRole = "student"
	RoleInstructor UserRole = "instructor"
	RoleAdmin      UserRole = "admin"
)

type User struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name         string             `json:"name" bson:"name"`
	Email        string             `json:"email" bson:"email"`
	PasswordHash string             `json:"-" bson:"password_hash"`
	Role         UserRole           `json:"role" bson:"role"`
	CreatedAt    time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at" bson:"updated_at"`
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

func NewUser(name, email, password string, role UserRole) (*User, error) {
	now := time.Now()
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           primitive.NewObjectID(),
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (u *User) Update(name, email string, role UserRole) {
	u.Name = name
	u.Email = email
	u.Role = role
	u.UpdatedAt = time.Now()
}

func (u *User) SetPassword(password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now()
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) CanAccessAdmin() bool {
	return u.Role == RoleAdmin || u.Role == RoleInstructor
}
