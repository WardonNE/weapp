package weapp

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	*gorm.DB
	driver string
	dsn    string
	config []gorm.Option
}

func OpenDB(driver, dsn string, config ...gorm.Option) (*Database, error) {
	db := &Database{
		driver: driver,
		dsn:    dsn,
		config: config,
	}
	if err := db.Connect(); err != nil {
		return nil, err
	}
	return db, nil
}

func (d *Database) Connect() error {
	var err error
	switch d.driver {
	case "mysql":
		d.DB, err = gorm.Open(mysql.Open(d.dsn), d.config...)
	case "sqlite":
		d.DB, err = gorm.Open(sqlite.Open(d.dsn), d.config...)
	default:
		err = fmt.Errorf("invalid driver: %s", d.driver)
	}
	return err
}
