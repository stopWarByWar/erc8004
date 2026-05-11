package model

import (
	"fmt"
	"strings"
	"time"

	agentcard "agent_identity/agentCard"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func ListAgentsBatch(lastUID uint64, limit int) ([]Agent, error) {
	var agents []Agent
	query := db.Order("uid ASC").Limit(limit)
	if lastUID > 0 {
		query = query.Where("uid > ?", lastUID)
	}
	if err := query.Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

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
	if agent == nil {
		return nil
	}

	var existing Agent
	err := db.Where("chain_id = ? AND identity_registry = ? AND agent_id = ?",
		agent.ChainID, agent.IdentityRegistry, agent.AgentID).First(&existing).Error

	if err == gorm.ErrRecordNotFound {
		return db.Omit("uid").Create(agent).Error
	}
	if err != nil {
		return err
	}

	updates := map[string]interface{}{
		"owner":        agent.Owner,
		"agent_uri":    agent.AgentURI,
		"block_number": agent.BlockNumber,
		"index":        agent.Index,
		"tx_hash":      agent.TxHash,
		"timestamps":   agent.Timestamps,
	}
	if existing.AgentURI != agent.AgentURI {
		updates["inserted"] = false
	}
	return db.Model(&Agent{}).Where("uid = ?", existing.UID).Updates(updates).Error
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
	// inserted IS NOT TRUE 会匹配 NULL 和 false（SQL 中 NULL <> true 为 NULL，不会匹配）
	err = db.Where("chain_id = ? and identity_registry = ? and inserted IS NOT TRUE and length(agent_uri) <> 0", chainID, identityRegistry).Limit(limit).Find(&agents).Error
	return agents, err
}

func UpdateAgent(chainID, identityRegistry, agentID string, agentProfile *agentcard.TokenURLResponse) (agent *Agent, err error) {
	// 检查 agentProfile 是否为 nil
	if agentProfile == nil {
		return nil, gorm.ErrRecordNotFound
	}

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
			serviceName := strings.TrimSpace(service.Name)
			if serviceName == "" {
				continue
			}
			var version string
			if service.Version != nil {
				version = *service.Version
			}
			services = append(services, Service{
				AgentUID:    existingAgent.UID,
				ServiceName: serviceName,
				Endpoint:    service.Endpoint,
				Version:     version,
			})

			if strings.ToLower(serviceName) == "oasf" {
				seenSkill := make(map[string]struct{})
				for _, skill := range service.Skills {
					skill = strings.TrimSpace(skill)
					if skill == "" {
						continue
					}
					skillKey := skill + "\x00" + version
					if _, ok := seenSkill[skillKey]; ok {
						continue
					}
					seenSkill[skillKey] = struct{}{}
					newOasfSkill := OASFSkill{
						AgentUID:  existingAgent.UID,
						SkillName: skill,
					}
					if service.Version != nil {
						newOasfSkill.Version = *service.Version
					}
					oasfSkills = append(oasfSkills, newOasfSkill)
				}
				seenDomain := make(map[string]struct{})
				for _, domain := range service.Domains {
					domain = strings.TrimSpace(domain)
					if domain == "" {
						continue
					}
					domainKey := domain + "\x00" + version
					if _, ok := seenDomain[domainKey]; ok {
						continue
					}
					seenDomain[domainKey] = struct{}{}
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
	if len(agentUIDs) == 0 {
		return nil
	}
	return db.Model(&Agent{}).Where("uid IN (?)", agentUIDs).Update("inserted", true).Error
}

func UpdateAgentWallet(chainID string, identityRegistry string, agentID string, agentWallet string, blockNumber uint64, index uint64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&Agent{}).
			Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", chainID, identityRegistry, agentID).
			Updates(map[string]interface{}{
				"agent_wallet": agentWallet,
				"block_number": blockNumber,
				"index":        index,
			}).Error; err != nil {
			return err
		}

		var fullAgent Agent
		if err := tx.Where("chain_id = ? AND identity_registry = ? AND agent_id = ?", chainID, identityRegistry, agentID).
			First(&fullAgent).Error; err != nil {
			return err
		}

		// 查找同一 chain_id + agent_wallet 的 stub（agent_id 为空说明是 commerce 创建的占位记录）
		var stub Agent
		err := tx.Where("chain_id = ? AND agent_wallet = ? AND agent_id = '' AND uid != ?", chainID, agentWallet, fullAgent.UID).
			First(&stub).Error
		if err == gorm.ErrRecordNotFound {
			return nil
		}
		if err != nil {
			return err
		}

		// 将 stub 的 commerce_actions 迁移到正式 agent
		if err := tx.Model(&CommerceAction{}).
			Where("agent_uid = ?", stub.UID).
			Update("agent_uid", fullAgent.UID).Error; err != nil {
			return err
		}

		// commerce_scores / commerce_scores_global 的 agent_uid 是联合主键，无法直接改，删掉让下次重算
		if err := tx.Where("agent_uid = ?", stub.UID).Delete(&CommerceScore{}).Error; err != nil {
			return err
		}
		if err := tx.Where("agent_uid = ?", stub.UID).Delete(&CommerceScoreGlobal{}).Error; err != nil {
			return err
		}

		return tx.Delete(&Agent{}, stub.UID).Error
	})
}

func GetAgentUID(chainID string, identityRegistry string, agentID string) (uint64, error) {
	var agentUID Agent
	err := db.Model(&Agent{}).Where("chain_id = ? and identity_registry = ? and agent_id = ?", chainID, identityRegistry, agentID).First(&agentUID).Error
	if err != nil {
		return 0, err
	}
	return agentUID.UID, nil
}

// FindOrCreateAgentByWallet looks up an agent by wallet address and chain.
// If no agent exists, it creates a minimal stub record (active=false, inserted=false)
// and returns the new UID. This supports the ERC-8183 Commerce Reputation system
// where on-chain addresses may not yet be registered in the identity registry.
func FindOrCreateAgentByWallet(chainID, wallet string) (uint64, error) {
	var agent Agent
	err := db.Where("agent_wallet = ? AND chain_id = ?", wallet, chainID).First(&agent).Error
	if err == nil {
		return agent.UID, nil
	}
	if err != gorm.ErrRecordNotFound {
		return 0, err
	}

	stub := Agent{
		ChainID:     chainID,
		AgentWallet: wallet,
		Active:      false,
		Inserted:    false,
	}

	// Use upsert to handle the race where the unique constraint
	// (chain_id, identity_registry, agent_id) is violated because stub
	// records share empty identity_registry + agent_id on the same chain.
	// We use agent_wallet as the conflict target since it's the actual
	// unique identifier for stub lookups.
	err = db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "chain_id"}, {Name: "identity_registry"}, {Name: "agent_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"agent_wallet", "active", "inserted"}),
	}).Create(&stub).Error
	if err != nil {
		// Fallback: if upsert still fails (e.g. unique constraint on a different column),
		// fall back to select-by-wallet to handle any other edge case.
		var existing Agent
		if err2 := db.Where("agent_wallet = ? AND chain_id = ?", wallet, chainID).First(&existing).Error; err2 == nil {
			return existing.UID, nil
		}
		return 0, err
	}
	return stub.UID, nil
}

func CreateMetadata(metadata *Metadata) error {
	if metadata == nil {
		return nil
	}
	// Use DB upsert to avoid insert+select+update round-trips.
	return db.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "chain_id"},
			{Name: "identity_registry"},
			{Name: "agent_id"},
			{Name: "key"},
		},
		DoUpdates: clause.AssignmentColumns([]string{"value", "block", "index", "tx_hash"}),
	}).Create(metadata).Error
}

//
//
//  ------------------- For API -------------------
//
//

func GetAgentByUID(uid uint64) (*Agent, error) {
	var agent *Agent
	if err := db.Where("uid = ?", uid).Find(&agent).Error; err != nil {
		return nil, err
	}
	return agent, nil
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

func GetAgentsByFilter(name *string, page, pageSize int, trustModelIDs, chainIDs, skills *[]string, x402Support, active, haveFeedback *bool) ([]*Agent, int64, error) {
	// 构建基础查询的辅助函数（不加 DISTINCT，由具体查询决定是否去重）
	buildBaseQuery := func() *gorm.DB {
		query := db.Model(&Agent{}).Where("length(name) <> 0")
		if name != nil && *name != "" {
			query = query.Where("LOWER(agents.name) LIKE LOWER(?)", "%"+*name+"%")
		}
		if trustModelIDs != nil && len(*trustModelIDs) > 0 {
			query = query.Joins("INNER JOIN trust_models ON agents.uid = trust_models.agent_uid").
				Where("trust_models.trust_model IN (?)", *trustModelIDs)
		}
		if chainIDs != nil && len(*chainIDs) > 0 {
			query = query.Where("agents.chain_id IN (?)", *chainIDs)
		}
		if skills != nil && len(*skills) > 0 {
			query = query.Joins("INNER JOIN oasf_skills ON agents.uid = oasf_skills.agent_uid").
				Where("oasf_skills.skill_name IN (?)", *skills)
		}
		if x402Support != nil && *x402Support {
			query = query.Where("agents.x402_support = ?", *x402Support)
		}
		if active != nil && *active {
			query = query.Where("agents.active = ?", *active)
		}
		if haveFeedback != nil && *haveFeedback {
			query = query.Where("agents.feedback_count > 0")
		}
		return query
	}

	// 查询总数：按 uid 去重后 Count
	var total int64
	countQuery := buildBaseQuery().Distinct("agents.uid")
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 如果总数为 0，直接返回
	if total == 0 {
		return []*Agent{}, 0, nil
	}

	// 分页查询 agent（创建新的查询对象），对 agents.* 去重，避免 join 导致重复
	var agents []*Agent
	dataQuery := buildBaseQuery().
		Distinct().
		Order("agents.uid DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize)

	if err := dataQuery.Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
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

type FilterInfoAmount struct {
	Name   string `json:"name" yaml:"name"`
	Amount int64  `json:"amount" yaml:"amount"`
}

func GetAgentAmountForEachTrustModel() ([]FilterInfoAmount, error) {
	var result []FilterInfoAmount
	if err := db.Model(&TrustModel{}).
		Select("trust_model AS name, COUNT(*) AS amount").
		Where("trust_model IN (?)", []string{
			agentcard.TrustModelReputation,
			agentcard.TrustModelCryptoEconomicValidation,
			agentcard.TrustModelTeeAttestation,
		}).
		Group("trust_model").
		Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func GetActiveAgentAmount() (int64, error) {
	var result int64
	if err := db.Model(&Agent{}).
		Select("COUNT(*) AS amount").
		Where("active = ?", true).
		Scan(&result).Error; err != nil {
		return 0, err
	}
	return result, nil
}

func GetX402SupportAgentAmount() (int64, error) {
	var result int64
	if err := db.Model(&Agent{}).
		Select("COUNT(*) AS amount").
		Where("x402_support = ?", true).
		Scan(&result).Error; err != nil {
		return 0, err
	}
	return result, nil
}

func GetAgentAmountForEachSkill(limit int) ([]FilterInfoAmount, error) {
	var result []FilterInfoAmount
	query := db.Model(&OASFSkill{}).
		Select("skill_name AS name, COUNT(DISTINCT agent_uid) AS amount").
		Group("skill_name").
		Order("amount DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}

	if err := query.Scan(&result).Error; err != nil {
		return nil, err
	}
	return result, nil
}

func GetAgentAmountWithFeedback() (int64, error) {
	var result int64
	if err := db.Model(&Agent{}).
		Select("COUNT(*) AS amount").
		Where("feedback_count > 0").
		Scan(&result).Error; err != nil {
		return 0, err
	}
	return result, nil
}

func SearchSkills(skill string, offset, limit int) ([]FilterInfoAmount, int, error) {
	if limit <= 0 {
		limit = 50
	}

	if offset < 0 {
		offset = 0
	}

	// 基础查询：按名称模糊匹配
	baseQuery := db.Model(&OASFSkill{}).
		Where("LOWER(skill_name) LIKE LOWER(?)", "%"+skill+"%")

	// 统计去重后的 skill 数量，用于总数
	var total int64
	if err := baseQuery.
		Select("COUNT(DISTINCT skill_name)").
		Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	// 按 skill 分组，按拥有该 skill 的 agent 数量倒序排序，并分页
	var result []FilterInfoAmount
	query := baseQuery.
		Select("skill_name AS name, COUNT(DISTINCT agent_uid) AS amount").
		Group("skill_name").
		Order("amount DESC").
		Offset(offset).
		Limit(limit)

	if err := query.Scan(&result).Error; err != nil {
		return nil, 0, err
	}

	return result, int(total), nil
}

func GetAgentAmountForEachChain(chainIds []string) (map[string]int64, error) {
	type chainAmount struct {
		ChainID string
		Amount  int64
	}

	var rows []chainAmount
	if err := db.Model(&Agent{}).
		Select("chain_id, COUNT(*) as amount").
		Where("chain_id IN (?)", chainIds).
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

func GetAgentAmountWithIn7Days() (int64, error) {
	var result int64
	if err := db.Model(&Agent{}).
		Select("COUNT(*) AS amount").
		Where("timestamps > ?", time.Now().AddDate(0, 0, -7).Unix()).
		Scan(&result).Error; err != nil {
		return 0, err
	}
	return result, nil
}

// GetAgentAgentListBy 按 filter 分页获取 Agent 列表并返回总数。
// filter[0]: 时间/uid — -1 按 uid 倒序，1 按 uid 正序，0 不按此排序
// filter[1]: 反馈数量 — -1 按 feedback_count 倒序，1 正序，0 不按此排序
// filter[2]: 最新 feedback 时间 — -1 按最新 feedback 时间倒序，1 正序，0 不按此排序
func GetAgentListBy(offset, limit int, filter []int8) ([]*Agent, error) {
	var agents []*Agent
	query := db.Model(&Agent{})

	if len(filter) > 2 && filter[2] != 0 {
		query = query.Joins("LEFT JOIN (SELECT agent_uid, MAX(timestamps) AS latest_feedback_ts FROM feedbacks WHERE revoked = false GROUP BY agent_uid) AS agent_latest_feedback ON agents.uid = agent_latest_feedback.agent_uid")
	}

	var orders []string
	if len(filter) > 0 && filter[0] != 0 {
		if filter[0] == -1 {
			orders = append(orders, "agents.uid DESC")
		} else {
			orders = append(orders, "agents.uid ASC")
		}
	}
	if len(filter) > 1 && filter[1] != 0 {
		if filter[1] == -1 {
			orders = append(orders, "agents.feedback_count DESC")
		} else {
			orders = append(orders, "agents.feedback_count ASC")
		}
	}
	if len(filter) > 2 && filter[2] != 0 {
		if filter[2] == -1 {
			orders = append(orders, "agent_latest_feedback.latest_feedback_ts DESC NULLS LAST")
		} else {
			orders = append(orders, "agent_latest_feedback.latest_feedback_ts ASC NULLS LAST")
		}
	}
	if len(orders) == 0 {
		orders = append(orders, "agents.uid DESC")
	}

	if err := query.Where("length(name) <> 0").Order(strings.Join(orders, ", ")).Offset(offset).Limit(limit).Find(&agents).Error; err != nil {
		return nil, err
	}
	return agents, nil
}

func GetAttestationCountByRecipient(agentUID uint64) (int, error) {
	var count int64
	recipient := fmt.Sprintf("0x%x", agentUID)
	err := db.Model(&Attestation{}).
		Where("recipient = ? AND revoked = false", recipient).
		Count(&count).Error
	return int(count), err
}
