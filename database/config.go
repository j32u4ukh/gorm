package database

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Database DatabaseConfig `yaml:"Database"`
}

type DatabaseConfig struct {
	UserName string `yaml:"UserName"`
	Password string `yaml:"Password"`
	Server   string `yaml:"Server"`
	Port     int    `yaml:"Port"`
	Name     string `yaml:"Name"`
}

var configs map[string]*Config

func init() {
	configs = make(map[string]*Config)
}

func NewConfig(path string) (*Config, error) {
	var c *Config
	var ok bool
	if c, ok = configs[path]; ok {
		return c, nil
	}
	var b []byte
	var err error
	b, err = ioutil.ReadFile(path)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read file %s.", path)
	}
	c = &Config{}
	err = yaml.Unmarshal(b, c)
	if err != nil {
		return nil, errors.Wrapf(err, "讀取 Config 時發生錯誤(path: %s)", path)
	}
	configs[path] = c
	return c, nil
}

func (c *Config) GetDatabase() DatabaseConfig {
	return c.Database
}
