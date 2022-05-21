package weapp

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func connect(driver, dsn string) (*gorm.DB, error) {
	switch driver {
	case "mysql":
		return gorm.Open(mysql.Open(dsn))
	case "sqlite":
		return gorm.Open(sqlite.Open(dsn))
	default:
		return nil, fmt.Errorf("invalid driver: `%s`", driver)
	}
}
