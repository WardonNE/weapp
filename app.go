package weapp

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wardonne/codec"
	"github.com/wardonne/inject"
	"github.com/wardonne/weapp/utils"
)

var release = false

func Release() {
	release = true
}

func isCommandLineMode() bool {
	return len(os.Args) > 1
}

type Application struct {
	*gin.Engine
	httpAddress  string
	databases    *sync.Map
	container    *inject.Container
	configration *Configration
	rootCmd      *cobra.Command
	BasePath     string
	WorkingPath  string
	executable   string
}

func NewApplication() *Application {
	app := &Application{}
	app.setBasePath()
	app.setWorkingPath()
	app.setRootCommand()
	app.withContainer()
	app.withConfigration()
	app.withEngine()
	app.databases = new(sync.Map)
	// store app instance into container
	if err := app.container.Provide("app", app); err != nil {
		panic(err)
	}
	return app
}

func (app *Application) setBasePath() {
	if basePath, err := utils.BasePath(); err != nil {
		panic(err)
	} else {
		app.BasePath = basePath
	}
}

func (app *Application) setWorkingPath() {
	if workingPath, err := utils.WorkingPath(); err != nil {
		panic(err)
	} else {
		app.WorkingPath = workingPath
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
	if release || isCommandLineMode() {
		gin.SetMode(gin.ReleaseMode)
	}
	app.Engine = gin.New()
}

func (app *Application) ConnectDatabase(name string, driver string, dsn string) {
	db, err := OpenDB(driver, dsn)
	if err != nil {
		panic(err)
	}
	app.RegisterDatabase(name, db)
}

func (app *Application) RegisterDatabase(name string, db *Database) {
	app.databases.Store(name, db)
}

func (app *Application) setRootCommand() {
	app.rootCmd = &cobra.Command{
		Run: func(cmd *cobra.Command, args []string) {
			gin.SetMode(gin.DebugMode)
			address := app.httpAddress
			host, err := cmd.Flags().GetString("host")
			if err != nil {
				panic(err)
			}
			port, err := cmd.Flags().GetString("port")
			if err != nil {
				panic(err)
			}
			if host != "" && port != "" {
				address = fmt.Sprintf("%s:%s", host, port)
			}
			app.Engine.Run(address)
		},
	}
	app.rootCmd.PersistentFlags().StringP("host", "H", "", "http server host")
	app.rootCmd.PersistentFlags().StringP("port", "P", "", "http server port")
}

func (app *Application) AddCommand(commands ...*cobra.Command) {
	app.rootCmd.AddCommand(commands...)
}

func (app *Application) Run(addr ...string) error {
	if len(addr) > 0 {
		app.httpAddress = addr[0]
	}
	return app.rootCmd.Execute()
}
