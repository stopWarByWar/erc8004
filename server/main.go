package main

import (
	"agent_identity/logger"
	"agent_identity/model"
	"agent_identity/server/api"
	"flag"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

var configFile = flag.String("f", "./conf.yaml", "the config file")

func main() {
	flag.Parse()

	logConf := &logger.Config{
		Level:        logrus.InfoLevel,
		ReportCaller: true,
		FilePath:     "./log/server",
	}

	_logger, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	config, err := initConf(*configFile)
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns)
	api.InitRouter(_logger, config.Mock)
	api.Run([]string{"*"}, config.Port)
}

type Config struct {
	Dns  string `yaml:"dns"`
	Port string `yaml:"port"`
	Mock bool   `yaml:"mock"`
}

func initConf(confPath string) (*Config, error) {
	config := &Config{}
	dataBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(dataBytes, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
