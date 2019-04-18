package testutils

import (
	"github.com/jinzhu/gorm"
	"github.com/nskondratev/api-page-go-back/db"
	"os"
)

func NewGormTestDB() (*gorm.DB, error) {
	connStr := os.Getenv("DB_CONNECTION_STRING")
	if len(connStr) < 1 {
		connStr = "api_page:api_page@/api_page_test?charset=utf8&parseTime=True"
	}
	return db.NewGorm(&db.MysqlDBConfig{
		ConnectionString: connStr,
	})
}
