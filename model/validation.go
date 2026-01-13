package model

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// GetLatestValidationRequestAndResponse 获取最新的验证请求和响应
// @param chainID string 链ID
// @param validationRegistry string 验证注册表
// @return uint64 请求块号
// @return uint64 请求索引
// @return error 错误
func GetLatestValidationRequestAndResponse(chainID string, validationRegistry string) (uint64, uint64, error) {
	var validation *Validation
	err := db.Where("chain_id = ? and validation_registry = ?", chainID, validationRegistry).Order("block_number DESC, index DESC").First(&validation).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, 0, err
	}
	if err == gorm.ErrRecordNotFound {
		return 0, 0, nil
	}
	if validation != nil {
		return validation.BlockNumber, validation.Index, nil
	} else {
		return 0, 0, nil
	}
}

func InsertValidation(validation *Validation) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "request_tx_hash"}},
		DoNothing: true,
	}).Create(validation).Error
}

func UpdateValidation(validation *Validation) error {
	params := map[string]interface{}{
		"response_uri":     validation.ResponseURI,
		"response_hash":    validation.ResponseHash,
		"tag1":             validation.Tag1,
		"response_tx_hash": validation.ResponseTxHash,
		"timestamps":       validation.Timestamps,
		"block_number":     validation.BlockNumber,
		"index":            validation.Index,
		"response":         validation.Response,
	}
	return db.Model(&Validation{}).Where("chain_id = ? and validation_registry = ? and request_hash = ?", validation.ChainID, validation.ValidationRegistry, validation.RequestHash).Updates(params).Error
}

// GetValidatorList 获取验证者列表
// @param page int 页码
// @param pageSize int 每页数量
// @return []*Validator 验证者列表
// @return int64 总数
// @return error 错误
func GetValidatorList(page int, pageSize int) ([]*Validator, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}
	var validatorList []*Validator
	err := db.Order("finished_amount DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&validatorList).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return validatorList, total, nil
}

func GetValidationListByAgent(chainID string, validationRegistry string, agentID string, page int, pageSize int, filter string) ([]*Validation, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}

	// 过滤逻辑：
	// - 如果 filter == "all"，返回所有验证记录。
	// - 如果 filter == "pending"，返回 request_tx_hash 非空且 response_tx_hash 为空的验证记录。
	// - 如果 filter == "finished"，返回 response_tx_hash 非空的验证记录。

	// 构建基础查询条件
	query := db.Where("chain_id = ? and validation_registry = ? and agent_id = ?", chainID, validationRegistry, agentID)

	// 根据 filter 参数添加过滤条件
	switch filter {
	case "pending":
		// response_tx_hash 为空或 NULL（request_tx_hash 总是非空，因为是主键）
		query = query.Where("response_tx_hash = '' OR response_tx_hash IS NULL")
	case "finished":
		// response_tx_hash 非空
		query = query.Where("response_tx_hash != '' AND response_tx_hash IS NOT NULL")
	default:
		// 如果 filter 不是上述值，默认返回所有记录
	}

	var validationList []*Validation
	err := query.Order("block_number DESC, index DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&validationList).Error
	if err != nil {
		return nil, 0, err
	}

	var total int64
	// 计算总数时使用相同的过滤条件
	countQuery := db.Model(&Validation{}).Where("chain_id = ? and validation_registry = ? and agent_id = ?", chainID, validationRegistry, agentID)
	switch filter {
	case "pending":
		countQuery = countQuery.Where("response_tx_hash = '' OR response_tx_hash IS NULL")
	case "finished":
		countQuery = countQuery.Where("response_tx_hash != '' AND response_tx_hash IS NOT NULL")
	}
	err = countQuery.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return validationList, total, nil
}

func GetValidationListByValidatorAddress(validatorAddress string, page int, pageSize int, filter string) ([]*Validation, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}

	// 过滤逻辑：
	// - 如果 filter == "all"，返回所有验证记录。
	// - 如果 filter == "pending"，返回 response_tx_hash 为空的验证记录。
	// - 如果 filter == "finished"，返回 response_tx_hash 非空的验证记录。

	// 构建基础查询条件（按验证者地址）
	query := db.Where("validator_address = ?", validatorAddress)

	// 根据 filter 参数添加过滤条件
	switch filter {
	case "pending":
		// response_tx_hash 为空或 NULL
		query = query.Where("response_tx_hash = '' OR response_tx_hash IS NULL")
	case "finished":
		// response_tx_hash 非空
		query = query.Where("response_tx_hash != '' AND response_tx_hash IS NOT NULL")
	case "all":
		// 不添加额外过滤条件，返回所有记录
	default:
		// 如果 filter 不是上述值，默认返回所有记录
	}

	var validationList []*Validation
	if err := query.Order("block_number DESC, index DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&validationList).Error; err != nil {
		return nil, 0, err
	}

	// 计算总数时使用相同的过滤条件
	var total int64
	countQuery := db.Model(&Validation{}).Where("validator_address = ?", validatorAddress)
	switch filter {
	case "pending":
		countQuery = countQuery.Where("response_tx_hash = '' OR response_tx_hash IS NULL")
	case "finished":
		countQuery = countQuery.Where("response_tx_hash != '' AND response_tx_hash IS NOT NULL")
	}
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	return validationList, total, nil
}
