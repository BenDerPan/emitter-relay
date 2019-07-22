package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type EmitterServerConfig struct {
	Host       string
	ChannelIn  string
	KeyIn      string
	ChannelOut string
	KeyOut     string
}

type Config struct {
	CurrentHop EmitterServerConfig
	NextHop    EmitterServerConfig
}

func (c *Config) Load(path string) error {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(f, c)
}
