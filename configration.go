package weapp

import (
	"github.com/spf13/viper"
)

type configration struct {
	configPath string
	configType string

	*viper.Viper
}

func newConfigration(configPath string) *configration {
	configration := new(configration)
	configration.Viper = viper.New()
	configration.configPath = configPath
	configration.configType = CONFIG_TYPE_JSON
	return configration
}

func (configration *configration) AddConfigration(filename string) error {
	return nil
}
