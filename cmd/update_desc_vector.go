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

		for _, a := range agents {
			content := strings.TrimSpace(a.Name + "\n" + a.Description)
			if content == "" {
				// 没有可用描述，删除该 agent 的向量（如果有）
				if err := model.InsertAgentVector(a.UID, a.IdentityRegistry, a.ChainID, a.Timestamps, "", nil); err != nil {
					return fmt.Errorf("delete desc_vector for agent %d failed: %w", a.UID, err)
				}
				lastUID = a.UID
				continue
			}

			if err := model.InsertAgentVector(a.UID, a.IdentityRegistry, a.ChainID, a.Timestamps, content, nil); err != nil {
				return fmt.Errorf("update desc_vector for agent %d failed: %w", a.UID, err)
			}
			lastUID = a.UID
		}
	}

	return nil
}
