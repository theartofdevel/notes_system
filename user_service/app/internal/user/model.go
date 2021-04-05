package user

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	UUID     string `json:"uuid" bson:"_id"`
	Email    string `json:"email" bson:"email"`
	Password string `json:"-" bson:"password,omitempty"`
}

func (u *User) CheckPassword(password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return fmt.Errorf("password does not match")
	}
	return nil
}

type CreateUserDTO struct {
	Email          string `json:"email" bson:"email"`
	Password       string `json:"password" bson:"password"`
	RepeatPassword string `json:"repeat_password" bson:"-"`
}

func (u *CreateUserDTO) GeneratePasswordHash() error {
	pwd, err := generatePasswordHash(u.Password)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

type UpdateUserDTO struct {
	Email       string `json:"email,omitempty" bson:"email,omitempty"`
	Password    string `json:"password,omitempty" bson:"password,omitempty"`
	OldPassword string `json:"old_password,omitempty" bson:"-"`
	NewPassword string `json:"new_password,omitempty" bson:"-"`
}

func (u *UpdateUserDTO) GeneratePasswordHash() error {
	pwd, err := generatePasswordHash(u.NewPassword)
	if err != nil {
		return err
	}
	u.Password = pwd
	return nil
}

func generatePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password due to error %w", err)
	}
	return string(hash), nil
}
