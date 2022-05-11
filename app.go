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
	if err := app.container.Provide("app", app); err != nil {
		panic(err)
	}
	return app
}

func (app *Application) setBasePath() {
	if basePath, err := os.Getwd(); err != nil {
		panic(err)
	} else {
		app.BasePath = basePath
	}
}

func (app *Application) setWorkingPath() {
	if workingPath, err := os.Executable(); err != nil {
		panic(err)
	} else {
		app.WorkingPath = filepath.Dir(workingPath)
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

func (app *Application) Make(key string, instance any) {

}

func (app *Application) withConfigration() {
	app.configration = new(Configration)
	app.container.Provide("config", app.configration, filepath.Join(app.BasePath, "config"))
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

func (app *Application) Config(key string) any {
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

func (app *Application) Release(release ...bool) *Application {
	releaseMode := true
	if len(release) > 0 {
		releaseMode = release[0]
	}
	if releaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	return app
}

func (app *Application) Run(addr ...string) error {
	return app.Engine.Run(addr...)
}
