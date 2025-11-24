package model

import (
	agentcard "agent_identity/agentCard"
	"errors"
	"fmt"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

			if agent.AgentCard.Provider.URL != nil {
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
			if err := tx.Where("agent_uid = ?", agentCardModel.UID).Delete(&SkillTags{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", agentCardModel.UID).Delete(&Skill{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", agentCardModel.UID).Delete(&Provider{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", agentCardModel.UID).Delete(&Capability{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", agentCardModel.UID).Delete(&TrustModel{}).Error; err != nil {
				return err
			}
			if err := tx.Where("agent_uid = ?", agentCardModel.UID).Delete(&Extension{}).Error; err != nil {
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

			capability := Capability{
				AgentUID: agentCardModel.UID,
				Streaming: func() bool {
					if agent.AgentCard.Capabilities.Streaming == nil {
						return false
					}
					return *agent.AgentCard.Capabilities.Streaming
				}(),
				PushNotifications: func() bool {
					if agent.AgentCard.Capabilities.PushNotifications == nil {
						return false
					}
					return *agent.AgentCard.Capabilities.PushNotifications
				}(),
				StateTransitionHistory: func() bool {
					if agent.AgentCard.Capabilities.StateTransitionHistory == nil {
						return false
					}
					return *agent.AgentCard.Capabilities.StateTransitionHistory
				}(),
			}
			if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&capability).Error; err != nil {
				return err
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

			var extensions []Extension
			for _, extension := range agent.AgentCard.Capabilities.Extensions {
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
		return nil
	})
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

func GetAgentList(page, pageSize int) ([]*Agent, int64, error) {
	var agentCards []*Agent
	if err := db.Where("length(name) <> 0").Offset((page - 1) * pageSize).Limit(pageSize).Find(&agentCards).Error; err != nil {
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
	err := db.Where("chain_id = ? and identity_registry = ?", chainID, registryAddress).Order("block_number DESC, index DESC").First(&agentRegistry).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}
	return agentRegistry.BlockNumber, agentRegistry.Index, nil
}

func CreateAgentRegistry(agentRegistry *AgentRegistry) error {
	return db.Transaction(func(tx *gorm.DB) error {
		agent := Agent{
			AgentID:          agentRegistry.AgentID,
			Owner:            agentRegistry.Owner,
			TokenURL:         agentRegistry.TokenURL,
			IdentityRegistry: agentRegistry.IdentityRegistry,
			ChainID:          agentRegistry.ChainID,
		}
		if err := tx.Clauses(clause.OnConflict{UpdateAll: true}).Create(&agent).Error; err != nil {
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
	err := db.Where("chain_id = ? and reputation_registry = ?", chainID, reputationRegistry).Order("block_number DESC, index DESC").First(&feedback).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var response *Response
	err = db.Where("chain_id = ? and reputation_registry = ?", chainID, reputationRegistry).Order("block_number DESC, index DESC").First(&response).Error
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
