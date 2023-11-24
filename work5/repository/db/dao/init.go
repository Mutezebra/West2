package dao

import (
	"context"
	"five/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"strings"
	"time"
)

var _db *gorm.DB

func InitMysql() {
	conf := config.Config.MySql["default"]
	conn := strings.Join([]string{conf.UserName, ":", conf.Password, "@tcp(", conf.DbHost, ":", conf.DbPort, ")/", conf.DbName, "?charset=utf8&parseTime=true"}, "")
	var ormLogger logger.Interface
	ormLogger = logger.Default
	ormLogger.LogMode(logger.Info)
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       conn,  // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 表名不加s
		},
	})
	if err != nil {
		log.Panic(err)
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(60)                  // 设置连接池
	sqlDB.SetMaxOpenConns(151)                 // 最大连接数
	sqlDB.SetConnMaxLifetime(time.Second * 30) // 最长生命周期

	_db = db
	migration()
}

// NewDBClient returns a new DB client
func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}
