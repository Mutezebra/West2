package model

import (
	"errors"
	"five/consts"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	SearchRecordNotExist = errors.New("search record not exist")
)

type User struct {
	gorm.Model
	UserName       string `gorm:"unique"`
	NickName       string
	Avatar         string
	PasswordDigest string
}

func (user *User) SetPassword(password string) (err error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), consts.PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(hashPassword)
	return nil
}

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	if err != nil {
		return false
	}
	return true
}
