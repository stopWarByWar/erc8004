package api

import (
	"agent_identity/helper"
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
	CardResponse, err := GetCardResponse(4)
	if err != nil {
		t.Errorf("GetCardResponse error: %v", err)
	}
	fmt.Println(CardResponse)
}

func TestGetCardList(t *testing.T) {
	initTest()
	agents, total, err := GetAgentList(1, 10)
	if err != nil {
		t.Errorf("GetAgentList error: %v", err)
	}
	for _, agent := range agents {
		fmt.Println(agent)
	}
	fmt.Println(total)
}

func TestGetCardResponseByTrustModel(t *testing.T) {
	initTest()
	agents, total, err := GetAgentListByFilter(1, 10, []string{"feedback", "inference-validation"}, []string{"97"})
	if err != nil {
		t.Errorf("GetAgentListByFilter error: %v", err)
	}
	for _, agent := range agents {
		fmt.Println(agent)
	}
	fmt.Println(total)
}

func TestSearchCardResponseBySkill(t *testing.T) {
	initTest()
	agents, total, err := SearchAgentListBySkill("traffic", 1, 10)
	if err != nil {
		t.Errorf("SearchAgentListBySkill error: %v", err)
	}
	for _, agent := range agents {
		fmt.Println(agent)
	}
	fmt.Println(total)
}

type Config struct {
	Dns                string `yaml:"dns"`
	Port               string `yaml:"port"`
	Mock               bool   `yaml:"mock"`
	S3Region           string `yaml:"s3_region"`
	S3BucketName       string `yaml:"s3_bucket_name"`
	AwsAccessKeyId     string `yaml:"aws_access_key_id"`
	AwsSecretAccessKey string `yaml:"aws_secret_access_key"`
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

func TestSetFeedback(t *testing.T) {
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

	helper.InitHelper(config.S3Region, config.S3BucketName, config.AwsAccessKeyId, config.AwsSecretAccessKey)

	request := UploadFeedbackRequest{
		UID:   4,
		Score: 5,
		Tag1:  "tag1",
		Tag2:  "tag2",
		FeedbackAuth: FeedbackAuth{
			AgentId:          "323",
			ClientAddress:    "0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345",
			IndexLimit:       1,
			Expiry:           1763901876,
			ChainId:          "97",
			IdentityRegistry: "0xA98A5542a1AaB336397d487e32021E0E48BEF717",
			SignerAddress:    "0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345",
		},
		FeedbackAuthSignature: "0x072eb93d9d6ee24174bae4c973f5f7d8eeb57919831e288a46e3ec1e6f11f044242efd0a9dbdd5943c3b389ac25c5d70736430bedd120166e1b4cc4b93ed91811b",
		Skill:                 "skill",
		Context:               "",
		Task:                  "",
		Capability:            "",
	}

	feedbackURI, feedbackHash, err := SetFeedback(request)
	if err != nil {
		t.Errorf("SetFeedback error: %v", err)
	}
	fmt.Println(feedbackURI)
	fmt.Println(feedbackHash)
}
