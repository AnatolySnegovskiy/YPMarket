package models

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"market/internal/entities"
	"market/internal/system"
)

type UserModel struct {
	*gorm.DB
	entities.UserEntity
}

func NewUserModel(db *gorm.DB) *UserModel {
	u := &UserModel{
		DB:         db,
		UserEntity: entities.UserEntity{},
	}

	return u
}

func (m *UserModel) Authenticate(email string, password string) (string, error) {
	user, _ := m.getUserByEmail(email)
	if user.ID == 0 || !checkPasswordHash(password, user.Password) {
		return "", fmt.Errorf("invalid login or password")
	}

	return system.CreateToken(user.ID)
}

func (m *UserModel) Registration(email string, password string) error {
	u, _ := m.getUserByEmail(email)

	if u.ID != 0 {
		return fmt.Errorf("user already exists")
	} else {
		hp, err := hashPassword(password)
		if err != nil {
			return fmt.Errorf("failed to hash password. Error: " + err.Error())
		}
		u = &entities.UserEntity{
			Email:    email,
			Password: hp,
		}
		return m.DB.Create(u).Error
	}
}

func (m *UserModel) UserExists(email string) (bool, error) {
	var count int64
	where := "email = ?"
	err := m.DB.Model(m.UserEntity).Where(where, email).Count(&count).Error
	return count > 0, err
}

func (m *UserModel) getUserByEmail(email string) (*entities.UserEntity, error) {
	user := &entities.UserEntity{}
	where := "email = ?"
	err := m.DB.Model(m.UserEntity).Where(where, email).First(user).Error
	return user, err
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
