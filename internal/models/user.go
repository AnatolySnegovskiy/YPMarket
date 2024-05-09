package models

import "gorm.io/gorm"

type User struct {
	gorm.Model `json:"gorm_._model"`
	Email      string  `json:"email,omitempty" gorm:"not null;unique;name:email;type:varchar(255)"`
	Password   string  `json:"password,omitempty" gorm:"not null;name:password;type:varchar(255)"`
	Balance    float64 `json:"balance,omitempty" gorm:"not null;default:0;name:balance;type:float"`
	Active     bool    `json:"active,omitempty" gorm:"not null;default:true;name:active;type:bool"`
}

func NewUserModel() *User {
	return &User{}
}
