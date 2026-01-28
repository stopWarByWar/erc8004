package model

import (
	agentcard "agent_identity/agentCard"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func GetLatestAgent(chainID string, identityRegistry string) (uint64, uint64, error) {
	var agent *Agent
	err := db.
		Where("chain_id = ? and identity_registry = ?", chainID, common.HexToAddress(identityRegistry).String()).
		Order("block_number DESC, index DESC").
		First(&agent).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return 0, 0, nil
	}
	return agent.BlockNumber, agent.Index, nil
}

func CreateAgent(agent *Agent) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var amount int64
		err := tx.Model(&Agent{}).Where("agent_id = ? AND chain_id = ? AND identity_registry = ?",
			agent.AgentID, agent.ChainID, agent.IdentityRegistry).
			Count(&amount).Error
		if err != nil {
			return err
		}
		if amount > 0 {
			return nil
		}

		// 创建新记录，忽略 uid 字段（让数据库自动生成主键）
		if err := tx.Omit("uid").Create(&agent).Error; err != nil {
			return err
		}
		return nil
	})
}

func UpdateAgentTokenURL(chainID, identityRegistry, agentID, agentURI string, blockNumber uint64, index uint64) error {
	updateMap := map[string]interface{}{
		"agent_uri":    agentURI,
		"inserted":     false,
		"block_number": blockNumber,
		"index":        index,
	}
	return db.
		Model(&Agent{}).
		Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, common.HexToAddress(identityRegistry).String(), agentID).
		Updates(updateMap).
		Error
}

func TransferOwnerShip(chainID, identityRegistry, agentID, newOwner string, blockNumber uint64, index uint64) error {
	updateMap := map[string]interface{}{
		"owner":        newOwner,
		"block_number": blockNumber,
		"index":        index,
	}

	return db.Model(&Agent{}).
		Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, common.HexToAddress(identityRegistry).String(), agentID).
		Updates(updateMap).Error
}

func GetUnInsertedAgents(chainID string, identityRegistry string, limit int) (agents []*Agent, err error) {
	err = db.Where("chain_id = ? and identity_registry = ? and inserted = ? and length(agent_uri) <> 0", chainID, identityRegistry, false).Limit(limit).Find(&agents).Error
	return agents, err
}

func UpdateAgent(chainID, identityRegistry, agentID string, agentProfile *agentcard.TokenURLResponse) (agent *Agent, err error) {
	var updatedAgent Agent
	err = db.Transaction(func(tx *gorm.DB) error {
		// 检查是否已存在相同 chain_id, identity_registry, agent_id 的记录
		var existingAgent Agent
		err := tx.Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", chainID, common.HexToAddress(identityRegistry).String(), agentID).First(&existingAgent).Error
		if err != nil {
			return err
		}

		// 清理旧的关联数据时应使用已存在记录的 UID，而不是尚未赋值的 agentCardModel.UID
		if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&Service{}).Error; err != nil {
			return err
		}
		if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&OASFSkill{}).Error; err != nil {
			return err
		}
		if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&OASFDomain{}).Error; err != nil {
			return err
		}
		if err := tx.Where("agent_uid = ?", existingAgent.UID).Delete(&TrustModel{}).Error; err != nil {
			return err
		}

		var services []Service
		var oasfSkills []OASFSkill
		var oasfDomains []OASFDomain
		for _, service := range agentProfile.Services {
			var version string
			if service.Version != nil {
				version = *service.Version
			}
			services = append(services, Service{
				AgentUID:    existingAgent.UID,
				ServiceName: service.Name,
				Endpoint:    service.Endpoint,
				Version:     version,
			})

			if service.Name == "oasf" {
				for _, skill := range service.Skills {
					newOasfSkill := OASFSkill{
						AgentUID:  existingAgent.UID,
						SkillName: skill,
					}
					if service.Version != nil {
						newOasfSkill.Version = *service.Version
					}
					oasfSkills = append(oasfSkills, newOasfSkill)
				}
				for _, domain := range service.Domains {
					newOasfDomain := OASFDomain{
						AgentUID: existingAgent.UID,
						Domain:   domain,
					}
					if service.Version != nil {
						newOasfDomain.Version = *service.Version
					}
					oasfDomains = append(oasfDomains, newOasfDomain)
				}
			}
		}
		if len(services) > 0 {
			if err := tx.Create(&services).Error; err != nil {
				return err
			}
		}

		if len(oasfSkills) > 0 {
			if err := tx.Create(&oasfSkills).Error; err != nil {
				return err
			}
		}
		if len(oasfDomains) > 0 {
			if err := tx.Create(&oasfDomains).Error; err != nil {
				return err
			}
		}

		var trustModels []TrustModel
		for _, trustModel := range agentProfile.SupportedTrust {
			trustModels = append(trustModels, TrustModel{
				AgentUID:   existingAgent.UID,
				TrustModel: trustModel,
			})
		}

		if len(trustModels) > 0 {
			if err := tx.Create(&trustModels).Error; err != nil {
				return err
			}
		}

		updateMap := map[string]interface{}{
			"type":         agentProfile.Type,
			"name":         agentProfile.Name,
			"description":  agentProfile.Description,
			"image":        agentProfile.Image,
			"x402_support": agentProfile.X402Support,
			"active":       agentProfile.Active,
			"inserted":     true,
		}
		if err := tx.Model(&Agent{}).Where("uid = ?", existingAgent.UID).Updates(updateMap).Error; err != nil {
			return err
		}

		// 查询更新后的 agent
		if err := tx.Where("uid = ?", existingAgent.UID).First(&updatedAgent).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &updatedAgent, nil
}

func UpdateAgentInserted(agentUIDs []uint64) error {
	return db.Model(&Agent{}).Where("uid IN (?)", agentUIDs).Update("inserted", true).Error
}

func UpdateAgentWallet(chainID string, identityRegistry string, agentID string, agentWallet string) error {
	return db.Model(&Agent{}).Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", chainID, identityRegistry, agentID).Update("agent_wallet", agentWallet).Error
}

func GetAgentUID(chainID string, identityRegistry string, agentID string) (uint64, error) {
	var agentUID Agent
	err := db.Model(&Agent{}).Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, identityRegistry, agentID).First(&agentUID).Error
	if err != nil {
		return 0, err
	}
	return agentUID.UID, nil
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

//
//
//  ------------------- For API -------------------
//
//

func GetAgentByUID(uid uint64) (*Agent, error) {
	var agentCards *Agent
	if err := db.Where("uid = ?", uid).Find(&agentCards).Error; err != nil {
		return nil, err
	}
	return agentCards, nil
}

func GetAgentsByUIDs(uids []uint64) ([]*Agent, error) {
	if len(uids) == 0 {
		return []*Agent{}, nil
	}

	var agents []*Agent
	if err := db.Where("uid IN ?", uids).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func GetOASFSkillsByAgentUID(uid uint64) ([]*OASFSkill, error) {
	var oasfSkills []*OASFSkill
	if err := db.Where("agent_uid = ?", uid).Find(&oasfSkills).Error; err != nil {
		return nil, err
	}
	return oasfSkills, nil
}

func GetOASFSkillsByAgentUIDs(uids []uint64) (map[uint64][]*OASFSkill, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var oasfSkills []*OASFSkill
	if err := db.Where("agent_uid IN ?", uids).Find(&oasfSkills).Error; err != nil {
		return nil, err
	}

	var oasfSkillsMap = make(map[uint64][]*OASFSkill)
	for _, oasfSkill := range oasfSkills {
		oasfSkillsMap[oasfSkill.AgentUID] = append(oasfSkillsMap[oasfSkill.AgentUID], oasfSkill)
	}
	return oasfSkillsMap, nil
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

func GetTrustModelsByAgentUID(uid uint64) ([]*TrustModel, error) {
	var trustModels []*TrustModel
	if err := db.Where("agent_uid = ?", uid).Find(&trustModels).Error; err != nil {
		return nil, err
	}
	return trustModels, nil
}

func GetA2AProvidersByAgentUIDs(uids []uint64) (map[uint64]*A2AProvider, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var providers []*A2AProvider
	if err := db.Where("agent_uid IN ?", uids).Find(&providers).Error; err != nil {
		return nil, err
	}

	var providersMap = make(map[uint64]*A2AProvider)
	for _, provider := range providers {
		providersMap[provider.AgentUID] = provider
	}
	return providersMap, nil
}

func GetA2AProviderByAgentUID(uid uint64) (*A2AProvider, error) {
	var provider *A2AProvider
	if err := db.Where("agent_uid = ?", uid).Find(&provider).Error; err != nil {
		return nil, err
	}
	return provider, nil
}

func GetServicesByAgentUIDs(uids []uint64) (map[uint64][]*Service, error) {
	if len(uids) == 0 {
		return nil, nil
	}

	var services []*Service
	if err := db.Where("agent_uid IN ?", uids).Find(&services).Error; err != nil {
		return nil, err
	}

	var servicesMap = make(map[uint64][]*Service)
	for _, service := range services {
		servicesMap[service.AgentUID] = append(servicesMap[service.AgentUID], service)
	}
	return servicesMap, nil
}

func GetServicesByAgentUID(uid uint64) ([]*Service, error) {
	var services []*Service
	if err := db.Where("agent_uid = ?", uid).Find(&services).Error; err != nil {
		return nil, err
	}
	return services, nil
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

func GetMetadata(chainID string, identityRegistry string, agentID string) ([]Metadata, error) {
	var metadatas []Metadata
	err := db.Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", chainID, identityRegistry, agentID).Find(&metadatas).Error
	if err != nil {
		return nil, err
	}
	return metadatas, nil
}

func SearchAgentsBySkill(skill string, page, pageSize int) ([]*Agent, int, error) {
	var agentUIDs []uint64

	query := db.Model(&OASFSkill{}).
		Select("DISTINCT agent_uid").
		Where("LOWER(skill_name) LIKE LOWER(?)", "%"+skill+"%")

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
