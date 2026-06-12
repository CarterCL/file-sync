package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type SyncConfig struct {
	DefaultUa string      `yaml:"default-ua" mapstructure:"default-ua"`
	SyncTasks []*SyncTask `yaml:"sync-tasks" mapstructure:"sync-tasks"`
}

type SyncTask struct {
	Tag       string      `yaml:"tag"`
	FilePairs []*FilePair `yaml:"file-pairs" mapstructure:"file-pairs"`
}

type FilePair struct {
	Url        string   `yaml:"url"`
	Path       string   `yaml:"path"`
	Convert    bool     `yaml:"convert"`
	UA         string   `yaml:"ua"`
	Extensions []string `yaml:"extensions"`
}

func InitConfig(configFile string) *SyncConfig {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("配置文件读取失败:", err)
		os.Exit(1)
	}

	var syncConfig SyncConfig
	if err := viper.Unmarshal(&syncConfig); err != nil {
		fmt.Println("配置反序列化失败:", err)
		os.Exit(1)
	}

	if syncConfig.DefaultUa == "" {
		syncConfig.DefaultUa = "file-sync/1.0"
	}

	return &syncConfig
}
