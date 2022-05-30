package weapp

var commandTpl = `package {{.package}}

import (
	"github.com/spf13/cobra"
	"github.com/wardonne/weapp"
)

type {{.commandName}}Command struct {
	*weapp.BaseCommand ` + "`inject:\"command:base\"`" + `
}

func (c *{{.commandName}}Command) Signation() string {
	return ""
}

func (c *{{.commandName}}Command) Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "{{.use}}",
		Short: "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	return cmd
}
`

var controllerTpl = `package {{.package}}

import (
	"github.com/gin-gonic/gin"
)

type {{.controllerName}}Controller struct {
}

func (c *{{.controllerName}}Controller) Init() {
}

func (c *{{.controllerName}}Controller) Index(ctx *gin.Context) {
}`

var migrationTpl = `package migrations

import (
	"github.com/wardonne/weapp"
)

type {{.structName}} struct {
	weapp.Migration
}

func (m *{{.structName}}) Version() string {
	return "{{.version}}"
}

func (m *{{.structName}}) Name() string {
	return "{{.name}}"
}
{{if ne .dbName "default"}}
func (m *{{.structName}}) Database() *weapp.Database {
	return m.app.DB("{{.dbName}}")
}
{{end}}

func (m *{{.structName}}) Commit() error {
	return nil
}

func (m *{{.structName}}) Rollback() error {
	return nil
}
`

var serviceTpl = `package {{.package}}

type {{.serviceName}}Service struct {
	
}

func (s *{{.serviceName}}Service) Init() {

}
`

var requestTpl = `package {{.package}}

import (
	"github.com/gin-gonic/gin"
)

type {{.requestName}}Request struct {

}

func (r *{{.requestName}}Request) Validate(ctx *gin.Context) {

}
`
