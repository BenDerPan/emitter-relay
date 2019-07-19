package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type EmitterServerConfig struct {
	Host    string
	Channel string
	Key     string
}

func (c *EmitterServerConfig) ChannelIn() string {
	return fmt.Sprintf("%sIn/", c.Channel)
}

func (c *EmitterServerConfig) ChannelOut() string {
	return fmt.Sprintf("%sOut/", c.Channel)
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
