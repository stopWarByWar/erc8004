package model

import (
	agentcard "agent_identity/agentCard"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	openai "github.com/sashabaranov/go-openai"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(dns, openaiAPIKey string) {
	var err error
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		panic(err)
	}

	baseURL := os.Getenv("OPENAI_BASE_URL")
	if baseURL != "" {
		cfg := openai.DefaultConfig(openaiAPIKey)
		cfg.BaseURL = baseURL
		openAIClient = openai.NewClientWithConfig(cfg)
		return
	}

	openAIClient = openai.NewClient(openaiAPIKey)

}

func InsertAgentCard(agent *agentcard.Agent) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var documentationURL string
		var providerURL string
		var url string
		var version string

		if agent.AgentCard != nil {
			if agent.AgentCard.DocumentationURL != nil {
				documentationURL = *agent.AgentCard.DocumentationURL
			}

			if agent.AgentCard.Provider != nil && agent.AgentCard.Provider.URL != nil {
				providerURL = *agent.AgentCard.Provider.URL
			}

			url = agent.AgentCard.URL
			version = agent.AgentCard.Version
		}

		agentCardModel := Agent{
			AgentID:          agent.AgentID,
			Owner:            agent.Owner,
			TokenURL:         agent.TokenURL,
			AgentWallet:      agent.AgentWallet,
			IdentityRegistry: agent.IdentityRegistry,
			ChainID:          agent.ChainID,
			A2AEndpoint:      agent.Endpoint,
			Type:             agent.Type,
			Name:             agent.Name,
			Description:      agent.Description,
			URL:              url,
			Image:            agent.Image,
			Version:          version,
			DocumentationURL: documentationURL,
			Timestamps:       agent.Timestamps,
			UserInterfaceURL: agent.UserInterfaceURL,
			Namespace:        agent.Namespace,
		}

		// 检查是否已存在相同 chain_id, identity_registry, agent_id 的记录
		var existingAgent Agent
		err := tx.Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", agentCardModel.ChainID, agentCardModel.IdentityRegistry, agentCardModel.AgentID).First(&existingAgent).Error
		if err == nil {
			// 清理旧的关联数据时应使用已存在记录的 UID，而不是尚未赋值的 agentCardModel.UID
			if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&SkillTags{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&Skill{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&Provider{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&Capability{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&TrustModel{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&Extension{}).Error; err != nil {
				return err
			}
		} else if err != gorm.ErrRecordNotFound {
			return err
		}
		// 不存在，插入新item
		if err := func() error {
			if existingAgent.UID != 0 {
				if err := tx.Model(&Agent{}).Where("uid = ?", existingAgent.UID).Updates(&agentCardModel).Error; err != nil {
					return err
				}
				agentCardModel.UID = existingAgent.UID
				return nil
			}
			if err := tx.Create(&agentCardModel).Error; err != nil {
				return err
			}
			return nil
		}(); err != nil {
			return err
		}

		if agent.AgentCard != nil {
			var skillTags []SkillTags
			var skills []Skill
			for _, skill := range agent.AgentCard.Skills {
				var skillDescription string
				if skill.Description != nil {
					skillDescription = *skill.Description
				}

				skills = append(skills, Skill{
					AgentUID:    agentCardModel.UID,
					ID:          skill.ID,
					Name:        skill.Name,
					Description: skillDescription,
				})
				for _, tag := range skill.Tags {
					skillTags = append(skillTags, SkillTags{
						AgentUID: agentCardModel.UID,
						SkillID:  skill.ID,
						Tag:      tag,
					})
				}
			}
			if len(skills) > 0 {
				if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&skills).Error; err != nil {
					return err
				}
			}

			if len(skillTags) > 0 {
				uniqueSkillTagsMap := make(map[string]SkillTags)
				for _, tag := range skillTags {
					key := fmt.Sprintf("%d|%s|%s", tag.AgentUID, tag.SkillID, tag.Tag)
					uniqueSkillTagsMap[key] = tag
				}
				uniqueSkillTags := make([]SkillTags, 0, len(uniqueSkillTagsMap))
				for _, tag := range uniqueSkillTagsMap {
					uniqueSkillTags = append(uniqueSkillTags, tag)
				}
				if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&uniqueSkillTags).Error; err != nil {
					return err
				}
			}

			// providerURL 已在上方声明与赋值

			if agent.AgentCard.Provider != nil {
				provider := Provider{
					AgentUID:     agentCardModel.UID,
					Organization: agent.AgentCard.Provider.Organization,
					URL:          providerURL,
				}
				if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&provider).Error; err != nil {
					return err
				}
			}

			capabilities := agent.AgentCard.Capabilities
			capability := Capability{
				AgentUID: agentCardModel.UID,
				Streaming: func() bool {
					if capabilities.Streaming == nil {
						return false
					}
					return *capabilities.Streaming
				}(),
				PushNotifications: func() bool {
					if capabilities.PushNotifications == nil {
						return false
					}
					return *capabilities.PushNotifications
				}(),
				StateTransitionHistory: func() bool {
					if capabilities.StateTransitionHistory == nil {
						return false
					}
					return *capabilities.StateTransitionHistory
				}(),
			}
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&capability).Error; err != nil {
				return err
			}

			var extensions []Extension
			for _, extension := range capabilities.Extensions {
				extensions = append(extensions, Extension{
					AgentUID: agentCardModel.UID,
					URI:      extension.URI,
					Required: func() bool {
						if extension.Required == nil {
							return false
						}
						return *extension.Required
					}(),
					Description: func() string {
						if extension.Description == nil {
							return ""
						}
						return *extension.Description
					}(),
				})
			}
			if len(extensions) > 0 {
				if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&extensions).Error; err != nil {
					return err
				}
			}
		}

		var trustModels []TrustModel
		for _, trustModel := range agent.SupportedTrust {
			trustModels = append(trustModels, TrustModel{
				AgentUID:   agentCardModel.UID,
				TrustModel: trustModel,
			})
		}

		if len(trustModels) > 0 {
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&trustModels).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func GetTokenURL(chainID, identityRegistry, agentID string) (string, error) {
	var agentRegistry AgentRegistry
	err := db.Model(&AgentRegistry{}).
		Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", chainID, identityRegistry, agentID).
		Order("block_number DESC, index DESC").
		First(&agentRegistry).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return agentRegistry.TokenURL, nil
}

func GetAgentByUID(uid uint64) (*Agent, error) {
	var agentCards *Agent
	if err := db.Where("uid = ?", uid).Find(&agentCards).Error; err != nil {
		return nil, err
	}
	return agentCards, nil
}

func GetAgentUID(chainID string, identityRegistry string, agentID string) (uint64, error) {
	var agentUID Agent
	err := db.Model(&Agent{}).Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, identityRegistry, agentID).First(&agentUID).Error
	if err != nil {
		return 0, err
	}
	return agentUID.UID, nil
}

func GetAgentsByUIDs(uids []uint64) ([]*Agent, error) {
	if len(uids) == 0 {
		return []*Agent{}, nil
	}

	var agentCards []*Agent
	if err := db.Where("uid IN ?", uids).Find(&agentCards).Error; err != nil {
		return nil, err
	}
	return agentCards, nil
}

func GetAgentByOwner(owner string) ([]*Agent, int64, error) {
	var agentCards []*Agent
	var total int64
	if err := db.Where("owner = ?", owner).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, total, nil
}

func GetAgentByAgentWallet(address string) ([]*Agent, int64, error) {
	var agentCards []*Agent
	var total int64
	if err := db.Where("agent_wallet = ?", address).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, total, nil
}

func GetSkillsByAgentUID(uid uint64) ([]*Skill, error) {
	var skills []*Skill
	if err := db.Where("agent_uid = ?", uid).Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func GetSkillsByAgentUIDs(uids []uint64) (map[uint64][]*Skill, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var skills []*Skill
	if err := db.Where("agent_uid IN ?", uids).Find(&skills).Error; err != nil {
		return nil, err
	}

	var skillsMap = make(map[uint64][]*Skill)
	for _, skill := range skills {
		skillsMap[skill.AgentUID] = append(skillsMap[skill.AgentUID], skill)
	}
	return skillsMap, nil
}

func GetSkillTagsByAgentUID(uid uint64) (map[string][]*SkillTags, error) {
	var skillTags []*SkillTags
	if err := db.Where("agent_uid = ?", uid).Find(&skillTags).Error; err != nil {
		return nil, err
	}
	var skillTagsMap = make(map[string][]*SkillTags)
	for _, skillTag := range skillTags {
		skillTagsMap[skillTag.SkillID] = append(skillTagsMap[skillTag.SkillID], skillTag)
	}
	return skillTagsMap, nil
}

func GetSkillTagsByAgentUIDs(uids []uint64) (map[uint64]map[string][]*SkillTags, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var skillTags []*SkillTags
	if err := db.Where("agent_uid IN ?", uids).Find(&skillTags).Error; err != nil {
		return nil, err
	}
	var agentSkillTagsMap = make(map[uint64]map[string][]*SkillTags)
	for _, skillTag := range skillTags {
		var skillTagsMap = make(map[string][]*SkillTags)
		skillTagsMap[skillTag.SkillID] = append(skillTagsMap[skillTag.SkillID], skillTag)
		agentSkillTagsMap[skillTag.AgentUID] = skillTagsMap
	}
	return agentSkillTagsMap, nil
}

func GetProviderByAgentUID(uid uint64) (*Provider, error) {
	var providers *Provider
	if err := db.Where("agent_uid = ?", uid).Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func GetProvidersByAgentUIDs(uids []uint64) (map[uint64]*Provider, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var providers []*Provider
	if err := db.Where("agent_uid IN ?", uids).Find(&providers).Error; err != nil {
		return nil, err
	}

	var providersMap = make(map[uint64]*Provider)
	for _, provider := range providers {
		providersMap[provider.AgentUID] = provider
	}
	return providersMap, nil
}

func GetTrustModelsByAgentUID(uid uint64) ([]*TrustModel, error) {
	var trustModels []*TrustModel
	if err := db.Where("agent_uid = ?", uid).Find(&trustModels).Error; err != nil {
		return nil, err
	}
	return trustModels, nil
}

func GetTrustModelsByAgentUIDs(uids []uint64) (map[uint64][]*TrustModel, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var trustModels []*TrustModel
	if err := db.Where("agent_uid IN ?", uids).Find(&trustModels).Error; err != nil {
		return nil, err
	}

	var trustModelsMap = make(map[uint64][]*TrustModel)
	for _, trustModel := range trustModels {
		trustModelsMap[trustModel.AgentUID] = append(trustModelsMap[trustModel.AgentUID], trustModel)
	}
	return trustModelsMap, nil
}

func GetExtensionByAgentUID(uid uint64) ([]*Extension, error) {
	var extensions []*Extension
	if err := db.Where("agent_uid = ?", uid).Find(&extensions).Error; err != nil {
		return nil, err
	}
	return extensions, nil
}

func GetExtensionsByAgentIDs(uids []uint64) (map[uint64][]*Extension, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var extensions []*Extension
	if err := db.Where("agent_uid IN ?", uids).Find(&extensions).Error; err != nil {
		return nil, err
	}

	var extensionsMap = make(map[uint64][]*Extension)
	for _, extension := range extensions {
		extensionsMap[extension.AgentUID] = append(extensionsMap[extension.AgentUID], extension)
	}
	return extensionsMap, nil
}

func GetAgentsByTrustModel(page, pageSize int, trustModelIDs []string) ([]*Agent, int64, error) {
	if len(trustModelIDs) == 0 {
		return nil, 0, nil
	}

	var agentUIDs []uint64
	// trustModelIDs 可能被识别成 'feedback,inference-validation,tee-attestation' 一个字符串而不是字符串数组, 这里做拆分
	ids := make([]string, 0)
	for _, v := range trustModelIDs {
		for _, id := range strings.Split(v, ",") {
			if id != "" {
				ids = append(ids, id)
			}
		}
	}
	if err := db.Model(&TrustModel{}).Select("agent_uid").Where("trust_model IN ?", ids).Scan(&agentUIDs).Error; err != nil {
		return nil, 0, err
	}

	if len(agentUIDs) == 0 {
		return nil, 0, nil
	}

	var agentCards []*Agent
	if err := db.Where("uid IN (?)", agentUIDs).Offset((page - 1) * pageSize).Limit(pageSize).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&Agent{}).Where("uid IN (?)", agentUIDs).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, total, nil
}

func GetAgentsByFilter(page, pageSize int, trustModelIDs, chainIDs []string) ([]*Agent, int64, error) {
	// 构建基础查询的辅助函数
	buildBaseQuery := func() *gorm.DB {
		query := db.Model(&Agent{}).Distinct("agents.uid")
		if len(trustModelIDs) > 0 {
			query = query.Joins("INNER JOIN trust_models ON agents.uid = trust_models.agent_uid").
				Where("trust_models.trust_model IN (?)", trustModelIDs)
		}
		if len(chainIDs) > 0 {
			query = query.Where("agents.chain_id IN (?)", chainIDs)
		}
		return query
	}

	// 查询总数（创建新的查询对象）
	var total int64
	countQuery := buildBaseQuery()
	if err := countQuery.Select("COUNT(DISTINCT agents.uid)").Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// 如果总数为 0，直接返回
	if total == 0 {
		return []*Agent{}, 0, nil
	}

	// 分页查询 agent_uid（创建新的查询对象）
	var agentUIDs []uint64
	dataQuery := buildBaseQuery()
	if err := dataQuery.Select("agents.uid").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&agentUIDs).Error; err != nil {
		return nil, 0, err
	}

	var agentCards []*Agent
	if err := db.Where("uid IN (?)", agentUIDs).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}

	return agentCards, total, nil
}

func GetAgentList(page, pageSize int) ([]*Agent, int64, error) {
	var agentCards []*Agent
	if err := db.
		Where("length(name) <> 0").
		Order("uid DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}
	var total int64
	if err := db.Model(&Agent{}).Where("length(name) <> 0").Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, total, nil
}

func GetLatestAgentRegistry(chainID string, registryAddress string) (uint64, uint64, error) {
	var agentRegistry *AgentRegistry
	err := db.Where("chain_id = ? and (identity_registry = ? or identity_registry = ?)", chainID, common.HexToAddress(registryAddress).String(), common.HexToAddress(registryAddress).Hex()).Order("block_number DESC, index DESC").First(&agentRegistry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return agentRegistry.BlockNumber, agentRegistry.Index, nil
}

func CreateAgentRegistry(agentRegistry *AgentRegistry) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var existingAgent Agent
		err := tx.Where("agent_id = ? AND chain_id = ? AND identity_registry = ?",
			agentRegistry.AgentID, agentRegistry.ChainID, agentRegistry.IdentityRegistry).
			First(&existingAgent).Error
		if err == nil {
			return nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		agent := Agent{
			AgentID:          agentRegistry.AgentID,
			Owner:            agentRegistry.Owner,
			TokenURL:         agentRegistry.TokenURL,
			IdentityRegistry: agentRegistry.IdentityRegistry,
			ChainID:          agentRegistry.ChainID,
		}
		if err := tx.Create(&agent).Error; err != nil {
			return err
		}
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&agentRegistry).Error; err != nil {
			return err
		}
		return nil
	})
}

func GetLatestFeedbackAndResponse(chainID string, reputationRegistry string) (uint64, uint64, error) {
	var feedback *Feedback
	err := db.Where("chain_id = ? and (reputation_registry = ? or reputation_registry = ?)", chainID, common.HexToAddress(reputationRegistry).String(), common.HexToAddress(reputationRegistry).Hex()).Order("block_number DESC, index DESC").First(&feedback).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var response *Response
	err = db.Where("chain_id = ? and (reputation_registry = ? or reputation_registry = ?)", chainID, common.HexToAddress(reputationRegistry).String(), common.HexToAddress(reputationRegistry).Hex()).Order("block_number DESC, index DESC").First(&response).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	if response.BlockNumber > feedback.BlockNumber {
		return response.BlockNumber, response.Index, nil
	} else if response.BlockNumber == feedback.BlockNumber && response.Index > feedback.Index {
		return response.BlockNumber, response.Index, nil
	}
	return feedback.BlockNumber, feedback.Index, nil
}

func CreateFeedback(feedback *Feedback) error {
	var amount int64
	err := db.Model(&Feedback{}).Where("tx_hash = ?", feedback.TxHash).Count(&amount).Error
	if err != nil || amount > 0 {
		return err
	}

	return db.Omit("uid").Create(&feedback).Error
}

func UpdateFeedbackRevoked(chainID string, agentID string, clientAddress string, feedbackIndex uint64) error {
	return db.Model(&Feedback{}).Where("chain_id = ? and agent_id = ? and client_address = ? and feedback_index = ?", chainID, agentID, clientAddress, feedbackIndex).Update("revoked", true).Error
}

func GetFeedbackUIDAndAgentUID(chainID string, agentID string, clientAddress string, feedbackIndex uint64) (uint64, uint64, error) {
	var feedback *Feedback
	err := db.Where("chain_id = ? and agent_id = ? and client_address = ? and feedback_index = ?", chainID, agentID, clientAddress, feedbackIndex).First(&feedback).Error
	if err != nil {
		return 0, 0, err
	}
	return feedback.UID, feedback.AgentUID, nil
}

func CreateResponse(response *Response) error {
	var amount int64
	err := db.Model(&Response{}).Where("block_number = ? and index = ?", response.BlockNumber, response.Index).Count(&amount).Error
	if err != nil || amount > 0 {
		return err
	}

	return db.Omit("uid").Create(&response).Error
}

func UpdateAgentTokenURL(chainID, identityRegistry, agentID, tokenURL string, blockNumber uint64, index uint64) error {
	updateMap := map[string]interface{}{
		"token_url":    tokenURL,
		"inserted":     false,
		"block_number": blockNumber,
		"index":        index,
	}
	return db.Model(&AgentRegistry{}).Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, identityRegistry, agentID).Updates(updateMap).Error
}

func TransferOwnerShip(chainID, identityRegistry, agentID, newOwner string, blockNumber uint64, index uint64) error {
	updateAgentRegistryMap := map[string]interface{}{
		"owner":        newOwner,
		"block_number": blockNumber,
		"index":        index,
	}

	updateAgentCardMap := map[string]interface{}{
		"owner": newOwner,
	}
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&AgentRegistry{}).Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, identityRegistry, agentID).Updates(updateAgentRegistryMap).Error; err != nil {
			return err
		}
		if err := tx.Model(&Agent{}).Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, identityRegistry, agentID).Updates(updateAgentCardMap).Error; err != nil {
			return err
		}
		return nil
	})
}

func SearchAgentsByName(name string, page, pageSize int) ([]*Agent, int, error) {
	var agents []*Agent
	query := db.Where("LOWER(name) LIKE LOWER(?)", "%"+name+"%").Offset((page - 1) * pageSize).Limit(pageSize)

	if err := query.Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	return agents, int(count), nil
}

func FilterSearchAgentsByName(name string, page, pageSize int, trustModelIDs, chainIDs []string) ([]*Agent, int64, error) {

	buildBaseQuery := func() *gorm.DB {
		query := db.Model(&Agent{}).Distinct("agents.uid").
			Where("LOWER(agents.name) LIKE LOWER(?)", "%"+name+"%")
		if len(trustModelIDs) > 0 {
			query = query.Joins("INNER JOIN trust_models ON agents.uid = trust_models.agent_uid").
				Where("trust_models.trust_model IN (?)", trustModelIDs)
		}
		if len(chainIDs) > 0 {
			query = query.Where("agents.chain_id IN (?)", chainIDs)
		}
		return query
	}

	// 查询总数（创建新的查询对象）
	var count int64
	countQuery := buildBaseQuery()
	if err := countQuery.Select("COUNT(DISTINCT agents.uid)").Scan(&count).Error; err != nil {
		return nil, 0, err
	}

	if count == 0 {
		return []*Agent{}, 0, nil
	}

	// 分页查询 agent_uid（创建新的查询对象）
	var agentUIDs []uint64
	dataQuery := buildBaseQuery()
	if err := dataQuery.Select("agents.uid").
		Limit(pageSize).
		Offset((page - 1) * pageSize).
		Scan(&agentUIDs).Error; err != nil {
		return nil, 0, err
	}

	var agents []*Agent
	if err := db.Where("uid IN (?)", agentUIDs).Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, count, nil
}

func SearchAgentsBySkill(skill string, page, pageSize int) ([]*Agent, int, error) {
	var agentUIDs []uint64

	query := db.Model(&Skill{}).
		Select("DISTINCT agent_uid").
		Where("LOWER(name) LIKE LOWER(?)", "%"+skill+"%")

	if err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Scan(&agentUIDs).
		Error; err != nil {
		return nil, 0, err
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	var agents []*Agent
	if err := db.Where("uid IN (?)", agentUIDs).Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, int(count), nil
}

func GetUnInsertedAgentRegistry(chainID string, identityRegistry string, limit int) ([]*AgentRegistry, error) {
	var agentRegistries []*AgentRegistry
	if err := db.Where("chain_id = ? and identity_registry = ? and inserted = ? and length(token_url) <> 0", chainID, identityRegistry, false).Limit(limit).Find(&agentRegistries).Error; err != nil {
		return nil, err
	}
	return agentRegistries, nil
}

func UpdateAgentRegistryInserted(agentIDs []string) error {
	return db.Model(&AgentRegistry{}).Where("agent_id IN (?)", agentIDs).Update("inserted", true).Error
}

func GetUnInsertedCommentAttestation(blockNumber uint64, index uint64, limit int, schemaUID string) (attestations []*Attestation, err error) {
	err = db.
		Where("block > ? or (block = ? and index > ?)", blockNumber, blockNumber, index).
		Where("schema_uid = ?", schemaUID).
		Order("block DESC, index DESC").
		Limit(limit).
		Find(&attestations).Error
	return
}

func checkPassport(addrs []string) (map[string]SimplePassportAccount, error) {
	passportMap := make(map[string]SimplePassportAccount)
	var accounts []SimplePassportAccount

	err := db.Table("passport_accounts").Select("address, twitter_name, avatar").Where("address IN (?) and length(twitter_name) > 0", addrs).Scan(&accounts).Error
	if err != nil {
		return nil, err
	}
	for _, account := range accounts {
		passportMap[account.Address] = account
	}
	return passportMap, nil
}

type SimplePassportAccount struct {
	Address string `gorm:"column:address;type:varchar(255)"`
	Name    string `gorm:"column:twitter_name;type:varchar(255)"`
	Avatar  string `gorm:"column:avatar;type:varchar(255)"`
}

func GetLatestAgentComment(chainID string) (uint64, uint64, error) {
	var agentComment *AgentComment
	err := db.Where("chain_id = ?", chainID).Order("block DESC, index DESC").First(&agentComment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return agentComment.Block, agentComment.Index, nil
}

func CreateAgentComments(agentComments []*AgentComment) error {
	if len(agentComments) == 0 {
		return nil
	}
	return db.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "comment_attestation_id"}},
			DoNothing: true,
		},
	).Create(&agentComments).Error
}

func GetCommentsByAgentUID(uid uint64, page, pageSize int) ([]*Comment, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}

	var comments []*Comment
	if err := db.Table("agent_comments ac").
		Select("ac.comment_attestation_id,ac.commenter, ac.agent_uid, ac.comment_text, ac.score, ac.timestamps").
		Where("ac.agent_uid = ? and ac.revoked = ?", uid, false).
		Order("ac.timestamps DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	var total int64

	if err := db.Model(&AgentComment{}).Where("agent_uid = ? and revoked = ?", uid, false).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var commentAddrs []string
	for _, comment := range comments {
		commentAddrs = append(commentAddrs, comment.Commenter)
	}

	passportMap, err := checkPassport(commentAddrs)
	if err != nil {
		return nil, 0, err
	}
	for _, comment := range comments {
		if account, ok := passportMap[comment.Commenter]; ok {
			comment.Name = account.Name
			comment.Avatar = account.Avatar
			comment.Passport = true
		} else {
			comment.Passport = false
		}

	}
	return comments, total, nil
}

func GetFeedbacksByAgentUID(uid uint64, page, pageSize int) ([]*FeedbackResp, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}
	var feedbacks []*FeedbackResp
	if err := db.Table("feedbacks f").
		Select("f.uid, f.agent_uid, f.chain_id, f.agent_id, f.reputation_registry, f.client_address, f.feedback_index, f.score, f.tag1, f.tag2, f.feedback_uri, f.feedback_hash, f.tx_hash, f.timestamps").
		Where("f.agent_uid = ? and f.revoked = ?", uid, false).
		Order("f.timestamps DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&Feedback{}).Where("agent_uid = ? and revoked = ?", uid, false).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var feedbackAddrs []string
	for _, feedback := range feedbacks {
		feedbackAddrs = append(feedbackAddrs, feedback.ClientAddress)
	}

	passportMap, err := checkPassport(feedbackAddrs)
	if err != nil {
		return nil, 0, err
	}
	for _, feedback := range feedbacks {
		if account, ok := passportMap[feedback.ClientAddress]; ok {
			feedback.Name = account.Name
			feedback.Avatar = account.Avatar
			feedback.Passport = true
		} else {
			feedback.Passport = false
		}
	}
	return feedbacks, total, nil
}

func CreateMetadata(metadata *Metadata) error {
	err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "identity_registry"},
			{Name: "agent_id"},
			{Name: "key"},
		},
		DoNothing: true,
	}).Create(metadata).Error

	if err != nil {
		return err
	}

	// 查询记录（无论是新插入的还是已存在的）
	var existing Metadata
	err = db.Where("chain_id = ? AND identity_registry = ? AND agent_id = ? AND key = ?",
		metadata.ChainID, metadata.IdentityRegistry, metadata.AgentID, metadata.Key).
		First(&existing).Error

	if err != nil {
		return err
	}

	// 检查是否有变化
	hasChange := existing.Value != metadata.Value ||
		existing.Block != metadata.Block ||
		existing.Index != metadata.Index ||
		existing.TxHash != metadata.TxHash

	if !hasChange {
		// 没有变化，什么都不做
		return nil
	}

	// 有变化，更新记录
	updateMap := map[string]interface{}{
		"value":   metadata.Value,
		"block":   metadata.Block,
		"index":   metadata.Index,
		"tx_hash": metadata.TxHash,
	}
	return db.Model(&Metadata{}).
		Where("chain_id = ? AND identity_registry = ? AND agent_id = ? AND key = ?",
			metadata.ChainID, metadata.IdentityRegistry, metadata.AgentID, metadata.Key).
		Updates(updateMap).Error
}

func GetMetadata(chainID string, identityRegistry string, agentID string) ([]Metadata, error) {
	var metadatas []Metadata
	err := db.Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", chainID, identityRegistry, agentID).Find(&metadatas).Error
	if err != nil {
		return nil, err
	}
	return metadatas, nil
}

// InsertAgentVector 插入或更新向量数据（upsert）
// 如果 agent_uid 已存在，则更新所有字段；否则插入新记录
// 如果 content 为空，则删除已存在的向量记录（如果有）
func InsertAgentVector(agentUID uint64, identityRegistry, chainID string, createTimestamp uint64, content string, metadata map[string]interface{}) error {
	// 检查是否已存在向量记录
	var existing AgentVector
	checkErr := db.Where("agent_uid = ?", agentUID).First(&existing).Error
	recordExists := checkErr == nil

	// 如果 content 为空，删除已存在的记录（如果有）
	if len(content) == 0 {
		if recordExists {
			// 记录存在，删除它
			return db.Where("agent_uid = ?", agentUID).Delete(&AgentVector{}).Error
		} else if errors.Is(checkErr, gorm.ErrRecordNotFound) {
			// 记录不存在，无需操作
			return nil
		} else {
			// 查询出错
			return checkErr
		}
	}

	// content 不为空，生成 embedding
	embedding, err := textToEmbedding(content)
	if err != nil {
		return fmt.Errorf("failed to convert content to embedding: %w", err)
	}

	agentVector := AgentVector{
		AgentUID:         agentUID,
		IdentityRegistry: identityRegistry,
		ChainID:          chainID,
		CreateTimestamp:  createTimestamp,
		Embedding:        formatVector(embedding),
		Content:          content,
	}

	if metadata != nil {
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		agentVector.Metadata = string(metadataJSON)
	} else {
		// 当 metadata 为 nil 时，设置为 JSON null，避免空字符串导致 PostgreSQL jsonb 类型错误
		agentVector.Metadata = "null"
	}

	// 使用 upsert：基于之前的查询结果，存在则更新，不存在则插入
	if recordExists {
		// 记录已存在，更新所有字段
		updateMap := map[string]interface{}{
			"identity_registry": agentVector.IdentityRegistry,
			"chain_id":          agentVector.ChainID,
			"create_timestamp":  agentVector.CreateTimestamp,
			"embedding":         agentVector.Embedding,
			"content":           agentVector.Content,
			"metadata":          agentVector.Metadata,
		}
		return db.Model(&AgentVector{}).Where("agent_uid = ?", agentUID).Updates(updateMap).Error
	} else if errors.Is(checkErr, gorm.ErrRecordNotFound) {
		// 记录不存在，插入新记录
		return db.Create(&agentVector).Error
	} else {
		// 查询出错（这种情况不应该发生，因为 content 为空时已经处理了）
		return checkErr
	}
}

func SearchSimilarVectors(desc string, limit int, threshold float64, filters *VectorSearchFilters) ([]uint64, error) {
	if desc == "" {
		return nil, nil
	}

	embedding, err := textToEmbedding(desc)
	if err != nil {
		return nil, fmt.Errorf("failed to convert content to embedding: %w", err)
	}

	queryVector := formatVector(embedding)

	if limit <= 0 {
		limit = 10
	}

	if threshold < 0 || threshold > 1 {
		threshold = 0.0
	}

	// 构建WHERE条件
	whereConditions := []string{}
	args := []interface{}{}
	needJoinTrustModels := false

	if filters != nil {
		// 如果提供了 TrustModel 过滤条件，需要通过 JOIN trust_models 表
		if len(filters.TrustModel) > 0 {
			needJoinTrustModels = true
			whereConditions = append(whereConditions, "tm.trust_model IN (?)")
			args = append(args, filters.TrustModel)
		}
		// 使用 AgentVector 表中的字段进行过滤
		if len(filters.IdentityRegistry) > 0 {
			whereConditions = append(whereConditions, "av.identity_registry IN (?)")
			args = append(args, filters.IdentityRegistry)
		}
		if len(filters.ChainID) > 0 {
			whereConditions = append(whereConditions, "av.chain_id IN (?)")
			args = append(args, filters.ChainID)
		}
	}

	// 构建 SQL 查询
	// 如果需要进行 TrustModel 过滤，则 JOIN trust_models 表
	var sql string
	if needJoinTrustModels {
		sql = `
			SELECT DISTINCT av.agent_uid, 
			       1 - (av.embedding <=> ?::vector) as similarity
			FROM agent_vectors av
			INNER JOIN trust_models tm ON av.agent_uid = tm.agent_uid
		`
	} else {
		sql = `
			SELECT av.agent_uid, 
			       1 - (av.embedding <=> ?::vector) as similarity
			FROM agent_vectors av
		`
	}

	whereClause := ""
	if len(whereConditions) > 0 {
		whereClause = " WHERE " + strings.Join(whereConditions, " AND ") + " AND "
	} else {
		whereClause = " WHERE "
	}

	sql += whereClause + `
		1 - (av.embedding <=> ?::vector) >= ?
		ORDER BY av.embedding <=> ?::vector
		LIMIT ?
	`

	allArgs := []interface{}{queryVector}
	allArgs = append(allArgs, args...)
	allArgs = append(allArgs, queryVector, threshold, queryVector, limit)

	var results []struct {
		AgentUID   uint64  `gorm:"column:agent_uid"`
		Similarity float64 `gorm:"column:similarity"`
	}

	err = db.Raw(sql, allArgs...).Scan(&results).Error
	if err != nil {
		return nil, err
	}

	agentUIDs := make([]uint64, len(results))
	for i, r := range results {
		agentUIDs[i] = r.AgentUID
	}

	return agentUIDs, nil
}

type VectorSearchFilters struct {
	TrustModel       []string
	IdentityRegistry []string
	ChainID          []string
}

// DeleteAgentVector 删除向量
func DeleteAgentVector(uid uint64) error {
	return db.Where("uid = ?", uid).Delete(&AgentVector{}).Error
}

// DeleteAgentVectorsByAgentUID 根据agent_uid删除所有相关向量
func DeleteAgentVectorsByAgentUID(agentUID uint64) error {
	return db.Where("agent_uid = ?", agentUID).Delete(&AgentVector{}).Error
}

// UpdateAgentVector 更新向量数据
func UpdateAgentVector(uid uint64, content string, metadata map[string]interface{}) error {
	embedding, err := textToEmbedding(content)
	if err != nil {
		return fmt.Errorf("failed to convert content to embedding: %w", err)
	}
	updates := map[string]interface{}{}

	if embedding != nil {
		updates["embedding"] = formatVector(embedding)
	}
	if content != "" {
		updates["content"] = content
	}
	if metadata != nil {
		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return fmt.Errorf("failed to marshal metadata: %w", err)
		}
		updates["metadata"] = string(metadataJSON)
	}

	if len(updates) == 0 {
		return nil
	}

	return db.Model(&AgentVector{}).Where("uid = ?", uid).Updates(updates).Error
}

// formatVector 将 []float32 转换为 PostgreSQL vector 格式的字符串
func formatVector(vec []float32) string {
	if len(vec) == 0 {
		return "[]"
	}

	var builder strings.Builder
	builder.WriteString("[")
	for i, v := range vec {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("%.6f", v))
	}
	builder.WriteString("]")
	return builder.String()
}

func GetAgentAmountForEachChain() (map[string]int64, error) {
	type chainAmount struct {
		ChainID string
		Amount  int64
	}

	var rows []chainAmount
	if err := db.Model(&Agent{}).
		Select("chain_id, COUNT(*) as amount").
		Group("chain_id").
		Scan(&rows).Error; err != nil {
		return nil, err
	}

	amounts := make(map[string]int64, len(rows))
	for _, row := range rows {
		amounts[row.ChainID] = row.Amount
	}
	return amounts, nil
}
