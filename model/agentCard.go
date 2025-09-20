package model

import (
	agentcard "agent_identity/agentCard"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(dns string) {
	var err error
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Error),
	})
	if err != nil {
		panic(err)
	}
}

func InsertAgentCard(agentCard agentcard.AgentCard) error {
	return db.Transaction(func(tx *gorm.DB) error {
		agentCardModel := AgentCard{
			AgentID:          agentCard.AgentID,
			AgentDomain:      agentCard.Domain,
			AgentAddress:     agentCard.AgentAddress,
			Name:             agentCard.Name,
			Description:      agentCard.Description,
			URL:              agentCard.URL,
			IconURL:          *agentCard.IconURL,
			Version:          agentCard.Version,
			DocumentationURL: *agentCard.DocumentationURL,
			FeedbackDataURI:  agentCard.FeedbackDataURI,
			ChainID:          agentCard.ChainID,
			Namespace:        agentCard.Namespace,
			Signature:        agentCard.Signature,
			UserInterface:    agentCard.UserInterface,
		}

		if err := tx.Create(&agentCardModel).Error; err != nil {
			return err
		}

		var skillTags []SkillTags
		var skills []Skill
		for _, skill := range agentCard.Skills {
			skills = append(skills, Skill{
				AgentID:     agentCard.AgentID,
				ID:          skill.ID,
				Name:        skill.Name,
				Description: *skill.Description,
			})
			for _, tag := range skill.Tags {
				skillTags = append(skillTags, SkillTags{
					AgentID: agentCard.AgentID,
					SkillID: skill.ID,
					Tag:     tag,
				})
			}
		}
		if len(skills) > 0 {
			if err := tx.Create(&skills).Error; err != nil {
				return err
			}
		}

		//todo: 如果重复如何处理
		if len(skillTags) > 0 {
			if err := tx.Create(&skillTags).Error; err != nil {
				return err
			}
		}

		provider := Provider{
			AgentID:      agentCard.AgentID,
			Organization: agentCard.Provider.Organization,
			URL:          *agentCard.Provider.URL,
		}
		if err := tx.Create(&provider).Error; err != nil {
			return err
		}

		capability := Capability{
			AgentID:                agentCard.AgentID,
			Streaming:              *agentCard.Capabilities.Streaming,
			PushNotifications:      *agentCard.Capabilities.PushNotifications,
			StateTransitionHistory: *agentCard.Capabilities.StateTransitionHistory,
		}
		if err := tx.Create(&capability).Error; err != nil {
			return err
		}

		var trustModels []TrustModel
		for _, trustModel := range agentCard.TrustModels {
			trustModels = append(trustModels, TrustModel{
				AgentID:    agentCard.AgentID,
				TrustModel: trustModel,
			})
		}

		if len(trustModels) > 0 {
			if err := tx.Create(&trustModels).Error; err != nil {
				return err
			}
		}

		var extensions []Extension
		for _, extension := range agentCard.Capabilities.Extensions {
			extensions = append(extensions, Extension{
				AgentID:     agentCard.AgentID,
				URI:         extension.URI,
				Required:    *extension.Required,
				Description: *extension.Description,
			})
		}
		if len(extensions) > 0 {
			if err := tx.Create(&extensions).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func GetAgentCard(agentID string) (*AgentCard, error) {
	var agentCards *AgentCard
	if err := db.Where("agent_id = ?", agentID).Find(&agentCards).Error; err != nil {
		return nil, err
	}
	return agentCards, nil
}

func GetAgentCards(agentIDs []string) ([]*AgentCard, error) {
	var agentCards []*AgentCard
	if err := db.Where("agent_id IN ?", agentIDs).Find(&agentCards).Error; err != nil {
		return nil, err
	}
	return agentCards, nil
}

func GetAgentCardByAddress(address string) (*AgentCard, error) {
	var agentCards *AgentCard
	if err := db.Where("agent_address = ?", address).Find(&agentCards).Error; err != nil {
		return nil, err
	}
	return agentCards, nil
}

func GetSkillsByAgentID(agentID string) ([]*Skill, error) {
	var skills []*Skill
	if err := db.Where("agent_id = ?", agentID).Find(&skills).Error; err != nil {
		return nil, err
	}
	return skills, nil
}

func GetSkillsByAgentIDs(agentIDs []string) (map[string][]*Skill, error) {
	if len(agentIDs) == 0 {
		return nil, nil
	}

	var skills []*Skill
	if err := db.Where("agent_id IN ?", agentIDs).Find(&skills).Error; err != nil {
		return nil, err
	}

	var skillsMap = make(map[string][]*Skill)
	for _, skill := range skills {
		skillsMap[skill.AgentID] = append(skillsMap[skill.AgentID], skill)
	}
	return skillsMap, nil
}

func GetSkillTagsByAgentID(agentID string) (map[string][]*SkillTags, error) {
	var skillTags []*SkillTags
	if err := db.Where("agent_id = ?", agentID).Find(&skillTags).Error; err != nil {
		return nil, err
	}
	var skillTagsMap = make(map[string][]*SkillTags)
	for _, skillTag := range skillTags {
		skillTagsMap[skillTag.SkillID] = append(skillTagsMap[skillTag.SkillID], skillTag)
	}
	return skillTagsMap, nil
}

func GetSkillTagsByAgentIDs(agentIDs []string) (map[string]map[string][]*SkillTags, error) {
	var skillTags []*SkillTags
	if err := db.Where("agent_id IN ?", agentIDs).Find(&skillTags).Error; err != nil {
		return nil, err
	}
	var agentSkillTagsMap = make(map[string]map[string][]*SkillTags)
	for _, skillTag := range skillTags {
		var skillTagsMap = make(map[string][]*SkillTags)
		skillTagsMap[skillTag.SkillID] = append(skillTagsMap[skillTag.SkillID], skillTag)
		agentSkillTagsMap[skillTag.AgentID] = skillTagsMap
	}
	return agentSkillTagsMap, nil
}

func GetProviderByAgentID(agentID string) (*Provider, error) {
	var providers *Provider
	if err := db.Where("agent_id = ?", agentID).Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

func GetProvidersByAgentIDs(agentIDs []string) (map[string]*Provider, error) {
	if len(agentIDs) == 0 {
		return nil, nil
	}

	var providers []*Provider
	if err := db.Where("agent_id IN ?", agentIDs).Find(&providers).Error; err != nil {
		return nil, err
	}

	var providersMap = make(map[string]*Provider)
	for _, provider := range providers {
		providersMap[provider.AgentID] = provider
	}
	return providersMap, nil
}

func GetTrustModelsByAgentID(agentID string) ([]*TrustModel, error) {
	var trustModels []*TrustModel
	if err := db.Where("agent_id = ?", agentID).Find(&trustModels).Error; err != nil {
		return nil, err
	}
	return trustModels, nil
}

func GetTrustModelsByAgentIDs(agentIDs []string) (map[string][]*TrustModel, error) {
	if len(agentIDs) == 0 {
		return nil, nil
	}

	var trustModels []*TrustModel
	if err := db.Where("agent_id IN ?", agentIDs).Find(&trustModels).Error; err != nil {
		return nil, err
	}

	var trustModelsMap = make(map[string][]*TrustModel)
	for _, trustModel := range trustModels {
		trustModelsMap[trustModel.AgentID] = append(trustModelsMap[trustModel.AgentID], trustModel)
	}
	return trustModelsMap, nil
}

func GetExtensionByAgentID(agentID string) ([]*Extension, error) {
	var extensions []*Extension
	if err := db.Where("agent_id = ?", agentID).Find(&extensions).Error; err != nil {
		return nil, err
	}
	return extensions, nil
}

func GetExtensionsByAgentIDs(agentIDs []string) (map[string][]*Extension, error) {
	if len(agentIDs) == 0 {
		return nil, nil
	}

	var extensions []*Extension
	if err := db.Where("agent_id IN ?", agentIDs).Find(&extensions).Error; err != nil {
		return nil, err
	}

	var extensionsMap = make(map[string][]*Extension)
	for _, extension := range extensions {
		extensionsMap[extension.AgentID] = append(extensionsMap[extension.AgentID], extension)
	}
	return extensionsMap, nil
}

func GetAgentCardsByTrustModel(page, pageSize, limit int, trustModelIDs []string) ([]*AgentCard, int64, error) {
	if len(trustModelIDs) == 0 {
		return nil, 0, nil
	}

	var agentIDs []string
	if err := db.Model(&TrustModel{}).Select("DISTINCT agent_id").Where("trust_model IN (?)", trustModelIDs).Scan(&agentIDs).Error; err != nil {
		return nil, 0, err
	}

	if len(agentIDs) == 0 {
		return nil, 0, nil
	}

	var agentCards []*AgentCard
	if err := db.Where("agent_id IN (?)", agentIDs).Offset((page - 1) * pageSize).Limit(limit).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&AgentCard{}).Where("agent_id IN (?)", agentIDs).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, total, nil
}

func GetAgentList(page, pageSize, limit int) ([]*AgentCard, int64, error) {
	var agentCards []*AgentCard
	if err := db.Offset((page - 1) * pageSize).Limit(limit).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}
	var total int64
	if err := db.Model(&AgentCard{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, total, nil
}

func GetLatestAgentRegistry() (uint64, uint64, error) {
	var agentRegistry *AgentRegistry
	err := db.Order("block_number DESC, index DESC").First(&agentRegistry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return agentRegistry.BlockNumber, agentRegistry.Index, nil
}

func CreateAgentRegistry(agentRegistries *AgentRegistry) error {
	return db.Create(&agentRegistries).Error
}

func SearchSkillsAgentCards(skill string) ([]*AgentCard, error) {
	var agentIDs []string
	if err := db.Model(&Skill{}).Select("DISTINCT agent_id").Where("name LIKE ?", "%"+skill+"%").Scan(&agentIDs).Error; err != nil {
		return nil, err
	}

	var agentCards []*AgentCard
	if err := db.Where("agent_id IN (?)", agentIDs).Find(&agentCards).Error; err != nil {
		return nil, err
	}
	return agentCards, nil
}
