package weapp

import (
	"os"
	"path/filepath"
)

type Application struct {
	container    *container
	configration *configration

	BasePath    string
	WorkingPath string
}

func NewApplication() *Application {
	app := &Application{}
	app.setBasePath()
	app.setWorkingPath()
	app.withContainer()
	app.withConfigration()
	return app
}

func (app *Application) setBasePath() {
	var err error
	if app.BasePath, err = os.Getwd(); err != nil {
		panic(err)
	}
}

func (app *Application) setWorkingPath() {
	var err error
	if app.WorkingPath, err = os.Executable(); err != nil {
		panic(err)
	} else {
		app.WorkingPath = filepath.Dir(app.WorkingPath)
	}
}

func (app *Application) withContainer() {
	app.container = newContainer()
}

func (app *Application) Container() *container {
	return app.container
}

func (app *Application) Provide(key string, instance any, callback func(instance any)) error {
	return app.container.Store(key, instance, callback)
}

func (app *Application) Load(key string) (any, bool) {
	return app.container.Load(key)
}

func (app *Application) Destory(key string) {
	app.container.Delete(key)
}

func (app *Application) withConfigration() {
	app.configration = newConfigration(filepath.Join(app.WorkingPath, "config"))
}

func (app *Application) Configration() *configration {
	return app.configration
}

func (app *Application) SetConfigPath(configPath string) {
	app.configration.configPath = configPath
}

func (app *Application) SetConfigType(configType string) {
	app.configration.configType = configType
}

func (app *Application) Configure(filename string) error {
	return app.configration.AddConfigration(filename)
}
