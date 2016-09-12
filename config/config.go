package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Title  string
	App    appInfo
	Owner  ownerInfo
	Master masetInfo
	Cron   cronInfo
}

type appInfo struct {
	Path string
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
}

type masetInfo struct {
	Server string
	Port   int
}

type cronInfo struct {
	DataPath string
}

func GetConfig() (*Config, error) {
	var config *Config

	path, _ := os.Getwd()

	//test
	//path = filepath.Dir(path)

	if _, err := toml.DecodeFile(filepath.Join(path, "config/config.toml"), &config); err != nil {
		return nil, err
	}
	if config.App.Path == "" {
		config.App.Path = path
	}
	return config, nil
}
