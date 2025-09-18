package server

import (
	"agent_identity/logger"
	"agent_identity/model"
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
	InitLogger(_logger)
	Run([]string{"*"}, config.Port)
}

type Config struct {
	Dns  string
	Port string
}

func initConf(confPath string) (*Config, error) {
	var config *Config
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
