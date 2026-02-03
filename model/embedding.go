package model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	openai "github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

var openAIClient *openai.Client
var ctx = context.Background()

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
	needJoinAgents := false
	needJoinSkills := false

	if filters != nil {
		// 如果提供了 TrustModel 过滤条件，需要通过 JOIN trust_models 表
		if filters.TrustModel != nil && len(*filters.TrustModel) > 0 {
			needJoinTrustModels = true
			whereConditions = append(whereConditions, "tm.trust_model IN (?)")
			args = append(args, filters.TrustModel)
		}
		// 使用 AgentVector 表中的字段进行过滤
		if filters.IdentityRegistry != nil && len(*filters.IdentityRegistry) > 0 {
			whereConditions = append(whereConditions, "av.identity_registry IN (?)")
			args = append(args, filters.IdentityRegistry)
		}
		if filters.ChainID != nil && len(*filters.ChainID) > 0 {
			whereConditions = append(whereConditions, "av.chain_id IN (?)")
			args = append(args, filters.ChainID)
		}
		// 根据技能过滤，需要 JOIN oasf_skills 表
		if filters.Skills != nil && len(*filters.Skills) > 0 {
			needJoinSkills = true
			whereConditions = append(whereConditions, "os.skill_name IN (?)")
			args = append(args, filters.Skills)
		}
		// 根据 Agent 状态过滤，需要 JOIN agents 表
		if filters.X402Support != nil && *filters.X402Support {
			needJoinAgents = true
			whereConditions = append(whereConditions, "a.x402_support = ?")
			args = append(args, true)
		}
		if filters.Active != nil && *filters.Active {
			needJoinAgents = true
			whereConditions = append(whereConditions, "a.active = ?")
			args = append(args, true)
		}
		if filters.HaveFeedback != nil && *filters.HaveFeedback {
			needJoinAgents = true
			whereConditions = append(whereConditions, "a.feedback_count > 0")
		}
	}

	// 构建 SQL 查询
	// 根据是否需要 JOIN 其他表决定是否使用 DISTINCT
	needDistinct := needJoinTrustModels || needJoinSkills

	sql := `
		SELECT `
	if needDistinct {
		sql += `DISTINCT `
	}
	sql += `
			av.agent_uid,
			1 - (av.embedding <=> ?::vector) as similarity
		FROM agent_vectors av
	`

	if needJoinTrustModels {
		sql += `
			INNER JOIN trust_models tm ON av.agent_uid = tm.agent_uid
		`
	}
	if needJoinAgents {
		sql += `
			INNER JOIN agents a ON av.agent_uid = a.uid
		`
	}
	if needJoinSkills {
		sql += `
			INNER JOIN oasf_skills os ON av.agent_uid = os.agent_uid
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
	TrustModel       *[]string
	IdentityRegistry *[]string
	ChainID          *[]string
	Skills           *[]string
	X402Support      *bool
	Active           *bool
	HaveFeedback     *bool
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

func textToEmbedding(text string) ([]float32, error) {
	if text == "" {
		return nil, fmt.Errorf("text is empty")
	}

	resp, err := openAIClient.CreateEmbeddings(ctx, openai.EmbeddingRequest{
		Model: openai.AdaEmbeddingV2,
		Input: text,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to call embedding API: %w", err)
	}

	if len(resp.Data) == 0 {
		return nil, fmt.Errorf("embedding result is empty")
	}

	return resp.Data[0].Embedding, nil
}
