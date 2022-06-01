package weapp

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/wardonne/codec"
	"github.com/wardonne/inject"
	"github.com/wardonne/weapp/utils"
	"gorm.io/gorm"
)

var release = false

func Release() {
	release = true
}

func isCommandLineMode() bool {
	return len(os.Args) > 1
}

type Application struct {
	engine       *gin.Engine
	httpAddress  string
	migrations   *Migrations
	container    *inject.Container
	configration *Configration
	rootCmd      *cobra.Command
	BasePath     string
	WorkingPath  string
}

func NewApplication() *Application {
	app := &Application{}
	app.withContainer()
	app.setBasePath()
	app.setWorkingPath()
	app.setRootCommand()
	app.withConfigration()
	app.withEngine()
	app.withMigrations()
	// store app instance into container
	if instance, err := app.container.LoadOrStore("app", app); err != nil {
		panic(err)
	} else {
		return instance.(*Application)
	}
}

func (app *Application) Init() {

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

func (app *Application) HTTPHost() string {
	return app.httpAddress
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

func (app *Application) MustProvide(key string, instance any) {
	if err := app.container.Provide(key, instance); err != nil {
		panic(err)
	}
}

func (app *Application) Load(key string) (any, bool) {
	return app.container.Load(key)
}

func (app *Application) MustLoad(key string) any {
	if instance, ok := app.Load(key); ok {
		return instance
	}
	panic(fmt.Errorf("invalid provider: %s", key))
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
	app.engine = gin.New()
}

func (app *Application) Router() *gin.Engine {
	return app.engine
}

func (app *Application) ConnectDatabase(name string, driver string, dsn string, config ...gorm.Option) {
	db, err := OpenDB(driver, dsn, config...)
	if err != nil {
		panic(err)
	}
	app.RegisterDatabase(name, db)
}

func (app *Application) RegisterDatabase(name string, db *Database) {
	app.MustProvide(fmt.Sprintf("db:%s", name), db)
}

func (app *Application) DB(name ...string) *Database {
	db := "default"
	if len(name) > 0 {
		db = name[0]
	}
	return app.MustLoad(fmt.Sprintf("db:%s", db)).(*Database)
}

func (app *Application) withMigrations() {
	app.migrations = &Migrations{
		migrations: make([]IMigration, 0),
	}
}

func (app *Application) Migration(migrations ...IMigration) {
	for _, migration := range migrations {
		if instance, err := app.container.LoadOrStore(fmt.Sprintf("migration:%s_%s", migration.Version(), migration.Name()), migration); err != nil {
			panic(err)
		} else {
			app.migrations.Migration(instance.(IMigration))
		}
	}
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
			app.engine.Run(address)
		},
	}
	app.rootCmd.PersistentFlags().StringP("host", "H", "", "http server host")
	app.rootCmd.PersistentFlags().StringP("port", "P", "", "http server port")
}

func (app *Application) AddCommand(commands ...ICommand) {
	for _, command := range commands {
		if instance, err := app.container.LoadOrStore(fmt.Sprintf("command:%s", command.Signation()), command); err != nil {
			panic(err)
		} else {
			app.rootCmd.AddCommand(instance.(ICommand).Command())
		}
	}
}

func (app *Application) withDefaultCommands() {
	defaultCommands := defaultCommands()
	for _, command := range defaultCommands {
		if instance, err := app.container.LoadOrStore(fmt.Sprintf("command:%s", command.Signation()), command); err != nil {
			panic(err)
		} else {
			app.rootCmd.AddCommand(instance.(ICommand).Command())
		}
	}
}

func (app *Application) Run(addr ...string) error {
	app.withDefaultCommands()

	if len(addr) > 0 {
		app.httpAddress = addr[0]
	}
	return app.rootCmd.Execute()
}
