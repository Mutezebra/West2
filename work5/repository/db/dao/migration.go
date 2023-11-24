package dao

import (
	"five/pkg/log"
	"five/repository/db/model"
)

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(
			&model.User{},
			&model.Group{},
			&model.Message{},
		)
	if err != nil {
		log.LogrusObj.Error(err)
	}
}
