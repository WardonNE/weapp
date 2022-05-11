package weapp

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"github.com/wardonne/codec"
)

type Configration struct {
	configPath  string
	configType  codec.CodecType
	configCodec codec.ICodec

	*viper.Viper
}

func (configration *Configration) Init(configPath string) {
	configration.Viper = viper.New()
	configration.configPath = configPath
	configration.configType = codec.CODEC_TYPE_JSON
	configration.configCodec = codec.NewCodecFactory().Get(codec.CODEC_TYPE_JSON)
}

func (configration *Configration) SetCodecType(typ codec.CodecType) {
	configration.configType = typ
	configration.configCodec = codec.NewCodecFactory().Get(typ)
}

func (configration *Configration) AddConfigration(modulename, filename string) error {
	codec := codec.NewCodecFactory().Get(codec.CodecType(configration.configType))
	f, err := os.Open(filepath.Join(configration.configPath, filename))
	if err != nil {
		return err
	}
	defer f.Close()
	var configs interface{}
	if err := codec.DecodeFromReader(f, &configs); err != nil {
		return err
	}
	configration.Viper.MergeConfigMap(map[string]any{
		modulename: configs,
	})
	return nil
}
