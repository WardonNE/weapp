package weapp

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/wardonne/codec"
	"github.com/wardonne/inject"
	"gorm.io/gorm"
)

type Application struct {
	*gin.Engine
	activeRouterPrefix []string

	databases *sync.Map

	container    *inject.Container
	configration *Configration

	BasePath    string
	WorkingPath string
}

func NewApplication() *Application {
	app := &Application{}
	app.setBasePath()
	app.setWorkingPath()
	app.withContainer()
	app.withConfigration()
	app.withEngine()
	app.databases = new(sync.Map)
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
	app.container = inject.NewContainer()
}

func (app *Application) Container() *inject.Container {
	return app.container
}

func (app *Application) Provide(key string, instance any) error {
	return app.container.Provide(key, instance)
}

func (app *Application) Load(key string) (any, bool) {
	return app.container.Load(key)
}

func (app *Application) withConfigration() {
	app.container.Provide("config", new(Configration), filepath.Join(app.BasePath, "config"))
}

func (app *Application) Configration() *Configration {
	return app.configration
}

func (app *Application) SetConfigPath(configPath string) {
	app.configration.configPath = configPath
}

func (app *Application) SetConfigType(configType codec.CodecType) {
	app.configration.configType = configType
}

func (app *Application) Configure(modulename, filename string) error {
	return app.configration.AddConfigration(modulename, filename)
}

func (app *Application) GetConfig(key string) any {
	return app.configration.Get(key)
}

func (app *Application) withEngine() {
	app.Engine = gin.New()
}

func (app *Application) RegisterDatabase(name string, db *gorm.DB, isDefault ...bool) {
	if len(isDefault) > 0 && isDefault[0] {
		app.databases.Store("default", db)
	}
	app.databases.Store(name, db)
}

func (app *Application) Run() error {
	return app.Engine.Run()
}
