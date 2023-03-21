package config

import (
	"db-go-game/pkg/conf"
	"db-go-game/pkg/utils"
	"flag"
)

type Config struct {
	Name   string      `yaml:"name"`
	Port   int         `yaml:"port"`
	WorkId int         `yaml:"work-id"`
	Log    string      `yaml:"log"`
	Etcd   *conf.Etcd  `yaml:"etcd"`
	Redis  *conf.Redis `yaml:"redis"`
	Mysql  *conf.Mysql `yaml:"mysql"`
	//Jaeger *conf.Jaeger `yaml:"jaeger"`
}

var (
	config = new(Config)
)

var confFile = flag.String("cfg", "./configs/api.yaml", "api config")

func init() {
	flag.Parse()
	utils.YamlToStruct(*confFile, config)
}

func GetConfig() *Config {
	return config
}
