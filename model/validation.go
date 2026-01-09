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
	var request *ValidationRequest
	var response *ValidationResponse
	err := db.Where("chain_id = ? and validation_registry = ?", chainID, validationRegistry).Order("block_number DESC, index DESC").First(&request).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, 0, err
	}
	err = db.Where("chain_id = ? and validation_registry = ?", chainID, validationRegistry).Order("block_number DESC, index DESC").First(&response).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 0, 0, err
	}

	if request != nil && response != nil {
		if request.BlockNumber > response.BlockNumber || (request.BlockNumber == response.BlockNumber && request.Index > response.Index) {
			return request.BlockNumber, request.Index, nil
		} else {
			return response.BlockNumber, response.Index, nil
		}
	}
	if request != nil {
		return request.BlockNumber, request.Index, nil
	}
	if response != nil {
		return response.BlockNumber, response.Index, nil
	}
	return 0, 0, nil
}

// InsertValidationRequest 插入验证请求
// @param validationRequest *ValidationRequest 验证请求
// @dev 首先检查 tx_hash 是否存在，如果存在跳过，如果不存在插入
// @return error 错误
func InsertValidationRequest(validationRequest *ValidationRequest) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tx_hash"}},
		DoNothing: true,
	}).Create(validationRequest).Error
}

// InsertValidationResponse 插入验证响应
// @param validationResponse *ValidationResponse 验证响应
// @dev 首先检查 txHash 是否存在，如果存在跳过，如果不存在插入
// @return error 错误
func InsertValidationResponse(validationResponse *ValidationResponse) error {
	return db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "tx_hash"}},
		DoNothing: true,
	}).Create(validationResponse).Error
}

// GetValidationRespList 获取验证响应列表
// @param chainID string 链ID
// @param validationRegistry string 验证注册表
// @param agentID string 代理ID
// @return []*ValidationResponse 验证响应列表
// @return int64 总数
// @return error 错误
func GetValidationRespList(chainID string, validationRegistry string, agentID string, page int, pageSize int) ([]*ValidationResponse, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}
	var respList []*ValidationResponse
	err := db.Where("chain_id = ? and validation_registry = ? and agent_id = ?", chainID, validationRegistry, agentID).Order("block_number DESC, index DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&respList).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = db.Where("chain_id = ? and validation_registry = ? and agent_id = ?", chainID, validationRegistry, agentID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return respList, total, nil
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
	err := db.Order("response_count DESC, request_count DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&validatorList).Error
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

// GetValidationResoListByValidatorAddress 获取验证者响应列表
// @param validatorAddress string 验证者地址
// @param page int 页码
// @param pageSize int 每页数量
// @return []*ValidationResponse 验证者响应列表
// @return int64 总数
// @return error 错误
func GetValidationRespListByValidatorAddress(validatorAddress string, page int, pageSize int) ([]*ValidationResponse, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}
	var respList []*ValidationResponse
	err := db.Where("validator_address = ?", validatorAddress).Order("block_number DESC, index DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&respList).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = db.Model(&ValidationResponse{}).Where("validator_address = ?", validatorAddress).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return respList, total, nil
}

// GetValidationReqListByValidatorAddress 获取验证者请求列表
// @param validatorAddress string 验证者地址
// @param page int 页码
// @param pageSize int 每页数量
// @return []*ValidationRequest 验证者请求列表
// @return int64 总数
// @return error 错误
func GetValidationReqListByValidatorAddress(validatorAddress string, page int, pageSize int) ([]*ValidationRequest, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}
	var reqList []*ValidationRequest
	err := db.Where("validator_address = ?", validatorAddress).Order("block_number DESC, index DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&reqList).Error
	if err != nil {
		return nil, 0, err
	}
	var total int64
	err = db.Model(&ValidationRequest{}).Where("validator_address = ?", validatorAddress).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}
	return reqList, total, nil
}
