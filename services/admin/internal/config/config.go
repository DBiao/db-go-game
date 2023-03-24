package config

import (
	"db-go-game/pkg/conf"
	"db-go-game/pkg/utils"
	"flag"
	"log"
)

type Config struct {
	Name        string           `yaml:"name"`
	Port        int              `yaml:"port"`
	WorkId      int              `yaml:"work-id"`
	Log         string           `yaml:"log"`
	LogicServer *conf.GrpcServer `yaml:"logic_server"`
	Etcd        *conf.Etcd       `yaml:"etcd"`
	Jaeger      *conf.Jaeger     `yaml:"jaeger"`
}

var (
	config = new(Config)
)

var confFile = flag.String("cfg", "./configs/admin.yaml", "api config")

func init() {
	flag.Parse()
	err := utils.YamlToStruct(*confFile, config)
	if err != nil {
		log.Fatalln(err)
	}
}

func GetConfig() *Config {
	return config
}
