package server

import (
	"agent_identity/logger"
	"agent_identity/model"
	"fmt"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func initTest() {
	logConf := &logger.Config{
		Level:        logrus.InfoLevel,
		ReportCaller: true,
		FilePath:     "./log/server",
	}

	_, err := logger.New(logConf)
	if err != nil {
		panic(err)
	}

	config, err := initConf("./config.yaml")
	if err != nil {
		panic(err)
	}

	model.InitDB(config.Dns)
}

func TestGetCardResponse(t *testing.T) {
	initTest()
	CardResponse, err := GetCardResponse("4")
	if err != nil {
		t.Errorf("GetCardResponse error: %v", err)
	}
	fmt.Println(CardResponse)
}

func TestGetCardList(t *testing.T) {
	initTest()
	CardList, total, err := GetCardList(1, 10, 10)
	if err != nil {
		t.Errorf("GetCardList error: %v", err)
	}
	for _, card := range CardList {
		fmt.Println(card)
	}
	fmt.Println(total)
}

func TestGetCardResponseByTrustModel(t *testing.T) {
	initTest()
	CardList, total, err := GetCardResponseByTrustModel(1, 10, 10, []string{"feedback", "inference-validation"})
	if err != nil {
		t.Errorf("GetCardResponseByTrustModel error: %v", err)
	}
	for _, card := range CardList {
		fmt.Println(card)
	}
	fmt.Println(total)
}

func TestSearchCardResponseBySkill(t *testing.T) {
	initTest()
	CardList, err := SearchCardResponseBySkill("traffic")
	if err != nil {
		t.Errorf("SearchCardResponseBySkill error: %v", err)
	}
	for _, card := range CardList {
		fmt.Println(card)
	}
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
