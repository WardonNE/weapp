package weapp

import (
	"fmt"

	"gorm.io/gorm"
)

type IModel interface {
	Database() string
}

type Model struct {
	*gorm.DB

	app *Application
}

func (m *Model) Init(app *Application) {
	m.app = app
	connectionName := m.Database()
	db, ok := app.databases.Load(connectionName)
	if ok {
		m.DB = db.(*gorm.DB)
	}
	panic(fmt.Errorf("invalid database connection `%s`", connectionName))
}

func (m *Model) Database() string {
	return "default"
}
