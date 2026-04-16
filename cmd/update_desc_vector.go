package main

import (
	"agent_identity/model"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config 复用 server/main.go 中的配置结构，只保留这里需要的字段
type Config struct {
	Dns          string `yaml:"dns"`
	OpenaiAPIKey string `yaml:"openai_api_key"`
}

var configFile = flag.String("f", "./config/conf.yaml", "the config file")

func main() {
	flag.Parse()

	cfg, err := initConf(*configFile)
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	model.InitDB(cfg.Dns, cfg.OpenaiAPIKey)

	if err := updateAllAgentsDescVector(); err != nil {
		log.Fatalf("update desc_vector failed: %v", err)
	}

	log.Println("update desc_vector finished")
}

func initConf(confPath string) (*Config, error) {
	conf := &Config{}
	dataBytes, err := os.ReadFile(confPath)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(dataBytes, conf); err != nil {
		return nil, err
	}
	return conf, nil
}

// 获取所有的 agents（分批），然后更新 desc_vector
func updateAllAgentsDescVector() error {
	const batchSize = 100
	var lastUID uint64

	for {
		agents, err := model.ListAgentsBatch(lastUID, batchSize)
		if err != nil {
			return fmt.Errorf("list agents batch failed: %w", err)
		}
		if len(agents) == 0 {
			break
		}

		upserts := make([]model.AgentVectorUpsert, 0, len(agents))
		for _, a := range agents {
			content := strings.TrimSpace(a.Name + "\n" + a.Description)
			upserts = append(upserts, model.AgentVectorUpsert{
				AgentUID:         a.UID,
				IdentityRegistry: a.IdentityRegistry,
				ChainID:          a.ChainID,
				CreateTimestamp:  a.Timestamps,
				Content:          content,
				Metadata:         nil,
			})
			lastUID = a.UID
		}
		if err := model.InsertAgentVectors(upserts, batchSize); err != nil {
			return fmt.Errorf("batch update desc_vector failed: %w", err)
		}
	}

	return nil
}
