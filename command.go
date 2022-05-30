package weapp

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"github.com/spf13/cobra"
	"github.com/wardonne/weapp/utils"
)

type ICommand interface {
	Command() *cobra.Command
	Signation() string
}

type BaseCommand struct {
	*Component `inject:"component"`
}

func (c *BaseCommand) Command() *cobra.Command {
	return nil
}

func (c *BaseCommand) Signation() string {
	return "base"
}

func defaultCommands() []ICommand {
	return []ICommand{
		new(MigrateCommand),
		new(MakeControllerCommand),
		new(MakeServiceCommand),
		new(MakeCommandCommand),
		new(MakeRequestCommand),
		new(MakeMigrationCommand),
		new(ListProviderCommand),
	}
}

type MigrateCommand struct {
	*BaseCommand `inject:"command:base"`
}

func (c *MigrateCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate [names...]",
		Short: "run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			names := args
			rollback, err := cmd.Flags().GetBool("rollback")
			if err != nil {
				return err
			}
			if rollback {
				c.migrations.Rollback(names...)
			} else {
				c.migrations.Commit(names...)
			}
			return nil
		},
	}
	cmd.PersistentFlags().BoolP("rollback", "r", false, "rollback migrations")
	return cmd
}

func (c *MigrateCommand) Signation() string {
	return "migrate"
}

type MakeMigrationCommand struct {
	*BaseCommand `inject:"command:base"`
}

func (c *MakeMigrationCommand) Signation() string {
	return "make:migration"
}

func (c *MakeMigrationCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:migration {name}",
		Short: "create a migration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			db, err := cmd.Flags().GetString("db")
			if err != nil {
				return err
			}
			outputDirpath, err := cmd.Flags().GetString("dist")
			if err != nil {
				return err
			}
			var data = map[string]string{
				"package":    filepath.Base(outputDirpath),
				"structName": utils.ToCamel(name),
				"dbName":     db,
				"name":       name,
				"version":    time.Now().Format("20060102150405"),
			}

			t := template.Must(template.New("migration").Parse(migrationTpl))
			if err := os.MkdirAll(outputDirpath, 0777); err != nil {
				return err
			}
			f, err := os.OpenFile(filepath.Join(c.BasePath, outputDirpath, name+".go"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			return t.Execute(f, data)
		},
	}
	cmd.PersistentFlags().String("db", "default", "database connection name")
	cmd.PersistentFlags().StringP("dist", "d", filepath.Join("database", "migrations"), "output directory")
	return cmd
}

type MakeControllerCommand struct {
	*BaseCommand `inject:"command:base"`
}

func (c *MakeControllerCommand) Signation() string {
	return "make:controller"
}

func (c *MakeControllerCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:controller {name}",
		Short: "create a controller file",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			outputDirpath, err := cmd.Flags().GetString("dist")
			if err != nil {
				return err
			}
			var data = map[string]string{
				"package":        filepath.Base(outputDirpath),
				"controllerName": utils.ToCamel(name),
			}
			t := template.Must(template.New("controller").Parse(controllerTpl))
			if err := os.MkdirAll(outputDirpath, 0777); err != nil {
				return err
			}
			f, err := os.OpenFile(filepath.Join(outputDirpath, name+".go"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			return t.Execute(f, data)
		},
	}
	cmd.PersistentFlags().StringP("dist", "d", filepath.Join(c.BasePath, "app", "controllers"), "output directory")
	return cmd
}

type MakeServiceCommand struct {
	*BaseCommand `inject:"command:base"`
}

func (c *MakeServiceCommand) Signation() string {
	return "make:service"
}

func (c *MakeServiceCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:service {name}",
		Short: "create a service file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			outputDirpath, err := cmd.Flags().GetString("dist")
			if err != nil {
				return err
			}
			var data = map[string]string{
				"package":     filepath.Base(outputDirpath),
				"serviceName": utils.ToCamel(name),
			}
			t := template.Must(template.New("service").Parse(serviceTpl))
			if err := os.MkdirAll(outputDirpath, 0777); err != nil {
				return err
			}
			f, err := os.OpenFile(filepath.Join(c.BasePath, outputDirpath, name+".go"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			return t.Execute(f, data)
		},
	}
	cmd.PersistentFlags().StringP("dist", "d", filepath.Join("app", "services"), "output directory")
	return cmd
}

type MakeCommandCommand struct {
	*BaseCommand `inject:"command:base"`
}

func (c *MakeCommandCommand) Signation() string {
	return "make:command"
}

func (c *MakeCommandCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:command {name}",
		Short: "create a command file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			outputDirpath, err := cmd.Flags().GetString("dist")
			if err != nil {
				return err
			}
			var data = map[string]string{
				"package":     filepath.Base(outputDirpath),
				"commandName": utils.ToCamel(name),
			}
			t := template.Must(template.New("command").Parse(commandTpl))
			if err := os.MkdirAll(outputDirpath, 0777); err != nil {
				return err
			}
			f, err := os.OpenFile(filepath.Join(c.BasePath, outputDirpath, name+".go"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			return t.Execute(f, data)
		},
	}
	cmd.PersistentFlags().StringP("dist", "d", filepath.Join("app", "console", "commands"), "output directory")
	return cmd
}

type MakeRequestCommand struct {
	*BaseCommand `inject:"command:base"`
}

func (c *MakeRequestCommand) Signation() string {
	return "make:request"
}

func (c *MakeRequestCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "make:request {name}",
		Short: "create a request file",
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			outputDirpath, err := cmd.Flags().GetString("dist")
			if err != nil {
				return err
			}
			data := map[string]string{
				"package":     filepath.Base(outputDirpath),
				"requestName": utils.ToCamel(name),
			}
			t := template.Must(template.New("request").Parse(requestTpl))
			if err := os.MkdirAll(outputDirpath, 0777); err != nil {
				return err
			}
			f, err := os.OpenFile(filepath.Join(c.BasePath, outputDirpath, name+".go"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0777)
			if err != nil {
				return err
			}
			return t.Execute(f, data)
		},
	}
	cmd.PersistentFlags().StringP("dist", "d", filepath.Join("app", "requests"), "output directory")
	return cmd
}

type ListProviderCommand struct {
	*BaseCommand `inject:"command:base"`
}

func (c *ListProviderCommand) Signation() string {
	return "list:provider"
}

func (c *ListProviderCommand) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list:provider",
		Short: "list all providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("application: ", c.Application)
			c.container.Each(func(key, value any) bool {
				fmt.Printf("[%s] %s\r\n", key.(string), reflect.TypeOf(value).String())
				return true
			})
			return nil
		},
	}
	return cmd
}
