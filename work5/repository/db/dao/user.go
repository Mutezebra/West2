package dao

import (
	"context"
	"five/repository/db/model"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

func GetUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func (dao *UserDao) FindUserByID(id uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).First(&user, "id=?", id).Error
	return
}

func (dao *UserDao) FindUserByUserName(username string) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).First(&user, "user_name=?", username).Error
	if err != nil {
		return nil, err
	}
	return
}

func (dao *UserDao) CreateUser(user *model.User) (err error) {
	err = dao.DB.Create(&user).Error
	return
}
