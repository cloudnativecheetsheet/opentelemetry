package config

import (
	"log"
	"userapi/utils"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	Port    string
	LogFile string
}

var Config ConfigList

func init() {
	LoadConfig()
	utils.LoggingSettings(Config.LogFile)
}

func LoadConfig() {
	cfg, err := ini.Load("config/config.ini")
	if err != nil {
		log.Fatalln(err)
	}
	Config = ConfigList{
		Port:    cfg.Section("web").Key("port").MustString("8080"),
		LogFile: cfg.Section("web").Key("logfile").String(),
	}
}
