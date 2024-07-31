package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

var Config *ConfigType

type ConfigType struct {
	Port uint 
	Database string
	SqlDir string
}

func Init(filename string) (*ConfigType, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	decoder := yaml.NewDecoder(file)
	conf := new(ConfigType)
	if err := decoder.Decode(conf); err != nil {
		return nil, err
	}
	Config = conf
	return conf, err
}
