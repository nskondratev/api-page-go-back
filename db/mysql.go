package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"time"
)

type MysqlDBConfig struct {
	ConnectionString string
}

func NewGorm(conf *MysqlDBConfig) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", conf.ConnectionString)

	if err != nil {
		return db, err
	}

	db.DB().SetConnMaxLifetime(time.Minute * 3)
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(20)

	db.LogMode(true)

	return db, nil
}
