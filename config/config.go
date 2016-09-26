package config

import (
	"os"
	"sync"
)

var lock sync.RWMutex

type Config struct {
	Title  string
	App    appInfo
	Owner  ownerInfo
	Master masterInfo
	Cron   cronInfo
	Aes    aesInfo
}

type appInfo struct {
	Path string
}

type ownerInfo struct {
	Name string
	Org  string `toml:"organization"`
}

type masterInfo struct {
	Server string
	Port   int
}

type cronInfo struct {
	DataPath string
}

type aesInfo struct {
	Key string
}

func GetConfig() (*Config, error) {
	path, _ := os.Getwd()
	return &Config{
		Title: "Webcron agent config",

		App: appInfo{
			Path: path,
		},

		Owner: ownerInfo{
			Name: "lianjia",
			Org:  "lianjia.com",
		},

		Master: masterInfo{
			//Server: "10.23.247.119",
			//Port:   8368,
			Server: "127.0.0.1",
			Port:   9999,
		},

		Cron: cronInfo{
			DataPath: "cron.data",
		},

		Aes: aesInfo{
			Key: "lianjia12798akljzmzmh.ahkjkljl@k",
		},
	}, nil

}

//从配置文件中读取
/*
func GetConfig() (*Config, error) {
	lock.RLock()
	defer lock.RUnlock()
	var config *Config
	path, _ := os.Getwd()

	if _, err := toml.DecodeFile(filepath.Join(path, "config/config.toml"), &config); err != nil {
		return nil, err
	}
	if config.App.Path == "" {
		config.App.Path = path
	}
	return config, nil
}
*/
