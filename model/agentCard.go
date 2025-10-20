package model

import (
	agentcard "agent_identity/agentCard"
	"errors"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	gormLogger "gorm.io/gorm/logger"
)

var db *gorm.DB

func InitDB(dns string) {
	var err error
	db, err = gorm.Open(postgres.Open(dns), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Info),
	})
	if err != nil {
		panic(err)
	}
}

func InsertAgentCard(agentCard agentcard.AgentCard) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var iconURL string
		var documentationURL string
		var providerURL string

		if agentCard.IconURL != nil {
			iconURL = *agentCard.IconURL
		}

		if agentCard.DocumentationURL != nil {
			documentationURL = *agentCard.DocumentationURL
		}

		if agentCard.Provider.URL != nil {
			providerURL = *agentCard.Provider.URL
		}

		agentCardModel := AgentCard{
			AgentID:          agentCard.AgentID,
			AgentDomain:      agentCard.Domain,
			AgentAddress:     agentCard.AgentAddress,
			Name:             agentCard.Name,
			Description:      agentCard.Description,
			URL:              agentCard.URL,
			IconURL:          iconURL,
			Version:          agentCard.Version,
			DocumentationURL: documentationURL,
			FeedbackDataURI:  agentCard.FeedbackDataURI,
			ChainID:          agentCard.ChainID,
			Namespace:        agentCard.Namespace,
			Signature:        agentCard.Signature,
			UserInterface:    agentCard.UserInterface,
		}

		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&agentCardModel).Error; err != nil {
			return err
		}

		var skillTags []SkillTags
		var skills []Skill
		for _, skill := range agentCard.Skills {
			var skillDescription string
			if skill.Description != nil {
				skillDescription = *skill.Description
			}

			skills = append(skills, Skill{
				AgentID:     agentCard.AgentID,
				ID:          skill.ID,
				Name:        skill.Name,
				Description: skillDescription,
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
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&skills).Error; err != nil {
				return err
			}
		}

		if len(skillTags) > 0 {
			uniqueSkillTagsMap := make(map[string]SkillTags)
			for _, tag := range skillTags {
				key := tag.AgentID + "|" + tag.SkillID + "|" + tag.Tag
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

		provider := Provider{
			AgentID:      agentCard.AgentID,
			Organization: agentCard.Provider.Organization,
			URL:          providerURL,
		}
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&provider).Error; err != nil {
			return err
		}

		capability := Capability{
			AgentID: agentCard.AgentID,
			Streaming: func() bool {
				if agentCard.Capabilities.Streaming == nil {
					return false
				}
				return *agentCard.Capabilities.Streaming
			}(),
			PushNotifications: func() bool {
				if agentCard.Capabilities.PushNotifications == nil {
					return false
				}
				return *agentCard.Capabilities.PushNotifications
			}(),
			StateTransitionHistory: func() bool {
				if agentCard.Capabilities.StateTransitionHistory == nil {
					return false
				}
				return *agentCard.Capabilities.StateTransitionHistory
			}(),
		}
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&capability).Error; err != nil {
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
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&trustModels).Error; err != nil {
				return err
			}
		}

		var extensions []Extension
		for _, extension := range agentCard.Capabilities.Extensions {
			extensions = append(extensions, Extension{
				AgentID: agentCard.AgentID,
				URI:     extension.URI,
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

func GetAgentCardsByTrustModel(page, pageSize int, trustModelIDs []string) ([]*AgentCard, int64, error) {
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
	if err := db.Where("agent_id IN (?)", agentIDs).Offset((page - 1) * pageSize).Limit(pageSize).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	if err := db.Model(&AgentCard{}).Where("agent_id IN (?)", agentIDs).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, total, nil
}

func GetAgentList(page, pageSize int) ([]*AgentCard, int64, error) {
	var agentCards []*AgentCard
	if err := db.Offset((page - 1) * pageSize).Limit(pageSize).Find(&agentCards).Error; err != nil {
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
	err := db.Order("block DESC, index DESC").First(&agentRegistry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return agentRegistry.BlockNumber, agentRegistry.Index, nil
}

func CreateAgentRegistry(agentRegistry *AgentRegistry) error {
	var amount int64
	err := db.Model(&AgentRegistry{}).Where("block_number = ? and index = ?", agentRegistry.BlockNumber, agentRegistry.Index).Count(&amount).Error
	if err != nil || amount > 0 {
		return err
	}

	return db.Create(&agentRegistry).Error
}

func SearchSkillsAgentCards(skill string, page, pageSize int) ([]*AgentCard, int, error) {
	var agentIDs []string
	// 使用ILIKE（PostgreSQL）或LOWER函数实现不区分大小写的检索
	if err := db.Model(&Skill{}).Select("DISTINCT agent_id").Where("LOWER(name) LIKE LOWER(?)", "%"+skill+"%").Offset((page - 1) * pageSize).Limit(pageSize).Scan(&agentIDs).Error; err != nil {
		return nil, 0, err
	}

	var agentCards []*AgentCard
	if err := db.Where("agent_id IN (?)", agentIDs).Find(&agentCards).Error; err != nil {
		return nil, 0, err
	}
	return agentCards, len(agentIDs), nil
}

func GetUnInsertedAgentRegistry(limit int) ([]*AgentRegistry, error) {
	var agentRegistries []*AgentRegistry
	if err := db.Where("inserted = ?", false).Limit(limit).Find(&agentRegistries).Error; err != nil {
		return nil, err
	}
	return agentRegistries, nil
}

func UpdateAgentRegistryInserted(agentIDs []string) error {
	return db.Model(&AgentRegistry{}).Where("agent_id IN (?)", agentIDs).Update("inserted", true).Error
}

func GetLatestAgentComment() (uint64, uint64, error) {
	var agentComment *AgentComment
	err := db.Order("block DESC, index DESC").First(&agentComment).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return agentComment.Block, agentComment.Index, nil
}

func GetUnInsertedCommentAttestation(blockNumber uint64, index uint64, limit int, schemaUID, attestor string) (attestations []*Attestation, err error) {
	err = db.
		Where("block > ? or (block = ? and index > ?)", blockNumber, blockNumber, index).
		Where("schema_uid = ?", schemaUID).
		Where("attestor = ?", attestor).
		Order("block DESC, index DESC").
		Limit(limit).
		Find(&attestations).Error
	return
}

func InsertAgentComments(agentComments []*AgentComment) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&agentComments).Error
}

func InsertAuthFeedback(authFeedbacks AuthFeedback) error {
	return db.Clauses(clause.OnConflict{DoNothing: true}).Create(&authFeedbacks).Error
}

func CheckAuthFeedbackExists(clientAddress string, agentServerID string) (bool, string, error) {
	var authFeedback *AuthFeedback
	err := db.Where("agent_client_address = ? and agent_server_id = ?", clientAddress, agentServerID).First(&authFeedback).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, "", nil
	} else if err != nil {
		return false, "", err
	}
	return true, authFeedback.AgentClientID, nil
}

func GetCommentList(agentID string, page, pageSize int, isAuthorized bool) ([]*Comment, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}
	if agentID == "" {
		return nil, 0, errors.New("agentID is required")
	}

	var comments []*Comment
	query := db.Table("agent_comments ac").
		Select("ac.comment_attestation_id,ac.commenter, ac.agent_client_id, ac.comment_text, ac.score, ac.timestamps, ag.icon_url as logo, ag.name as name").
		Joins("LEFT JOIN agent_cards ag ON ac.agent_client_id = ag.agent_id").
		Where("ac.agent_server_id = ?", agentID).
		Order("ac.timestamps DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	if isAuthorized {
		query = query.Where("ac.is_authorized = ?", isAuthorized)
	}

	if err := query.Scan(&comments).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := db.Table("agent_comments").
		Where("agent_server_id = ?", agentID)
	if isAuthorized {
		countQuery = countQuery.Where("is_authorized = ?", isAuthorized)
	}
	if err := countQuery.Count(&total).Error; err != nil {
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
		account, ok := passportMap[comment.Commenter]
		comment.Passport = ok
		if comment.Name == "" {
			comment.Name = account.Name
		}
		if comment.Logo == "" {
			comment.Logo = account.Logo
		}
	}

	return comments, total, nil
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
	Logo    string `gorm:"column:avatar;type:varchar(255)"`
}
