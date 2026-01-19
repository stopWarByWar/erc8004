package logic

import (
	"agent_identity/model"
	serverTypes "agent_identity/server/api/types"
	"fmt"
)

func GetAgentValidationList(chainID, validationRegistry, agentID string, page, pageSize int, filter string) ([]*serverTypes.Validation, int64, error) {
	validationResponses, total, err := model.GetValidationListByAgent(chainID, validationRegistry, agentID, page, pageSize, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get validation responses: %v", err)
	}

	var validationResponsesList []*serverTypes.Validation
	for _, validationResponse := range validationResponses {
		validationResponsesList = append(validationResponsesList, &serverTypes.Validation{
			AgentUID:           validationResponse.AgentUID,
			ChainID:            validationResponse.ChainID,
			AgentID:            validationResponse.AgentID,
			ValidationRegistry: validationResponse.ValidationRegistry,

			ValidatorAddress: validationResponse.ValidatorAddress,

			RequestURI:  validationResponse.RequestURI,
			RequestHash: validationResponse.RequestHash,

			Response:     validationResponse.Response,
			ResponseURI:  validationResponse.ResponseURI,
			ResponseHash: validationResponse.ResponseHash,
			Tag1:         validationResponse.Tag1,
			Timestamps:   validationResponse.Timestamps,
		})
	}
	return validationResponsesList, total, nil
}

func GetValidatorList(page, pageSize int) ([]*serverTypes.ValidatorListInfo, int64, error) {
	validatorList, total, err := model.GetValidatorList(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get validator list: %v", err)
	}
	var validatorListInfoList []*serverTypes.ValidatorListInfo
	for i, validator := range validatorList {
		validatorListInfoList = append(validatorListInfoList, &serverTypes.ValidatorListInfo{
			Rank:           uint64(i + (page-1)*pageSize + 1),
			Address:        validator.Address,
			PendingAmount:  validator.PendingAmount,
			FinishedAmount: validator.FinishedAmount,
		})
	}
	return validatorListInfoList, total, nil
}

func GetValidatorValidationList(validatorAddress string, page, pageSize int, filter string) ([]*serverTypes.Validation, int64, error) {
	validations, total, err := model.GetValidationListByValidatorAddress(validatorAddress, page, pageSize, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get validator requests: %v", err)
	}

	var validationsList []*serverTypes.Validation
	for _, validation := range validations {
		validationsList = append(validationsList, &serverTypes.Validation{
			AgentUID:           validation.AgentUID,
			ChainID:            validation.ChainID,
			AgentID:            validation.AgentID,
			ValidationRegistry: validation.ValidationRegistry,

			ValidatorAddress: validation.ValidatorAddress,

			RequestURI:  validation.RequestURI,
			RequestHash: validation.RequestHash,

			Response:     validation.Response,
			ResponseURI:  validation.ResponseURI,
			ResponseHash: validation.ResponseHash,
			Tag1:         validation.Tag1,
			Timestamps:   validation.Timestamps,
		})
	}
	return validationsList, total, nil
}

func GetValidatorByAddress(address string) (any, error) {
	validatorInfo, err := model.GetValidatorByAddress(address)
	if err != nil {
		return nil, fmt.Errorf("fail to get validator: %v", err)
	}
	return validatorInfo, nil
}
