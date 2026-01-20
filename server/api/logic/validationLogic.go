package logic

import (
	"agent_identity/config"
	"agent_identity/model"
	serverTypes "agent_identity/server/api/types"
	"fmt"
	"strings"
)

func GetAgentValidationList(uid uint64, page, pageSize int, filter string) ([]*serverTypes.Validation, int64, error) {
	validationResponses, total, err := model.GetValidationListByAgent(uid, page, pageSize, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get validation responses: %v", err)
	}

	var validationResponsesList []*serverTypes.Validation
	for _, validationResponse := range validationResponses {
		status := "pending"
		if len(strings.Trim(validationResponse.ResponseTxHash, " ")) > 0 {
			status = "finished"
		}
		validationResponsesList = append(validationResponsesList, &serverTypes.Validation{
			AgentUID:           validationResponse.AgentUID,
			ChainID:            validationResponse.ChainID,
			AgentID:            validationResponse.AgentID,
			ValidationRegistry: validationResponse.ValidationRegistry,

			ValidatorAddress: validationResponse.ValidatorAddress,
			ValidatorLogo:    validationResponse.ValidatorLogo,
			RequestURI:       validationResponse.RequestURI,
			RequestHash:      strings.Trim(validationResponse.RequestHash, " "),

			Response:     validationResponse.Response,
			ResponseURI:  validationResponse.ResponseURI,
			ResponseHash: strings.Trim(validationResponse.ResponseHash, " "),
			Tag1:         validationResponse.Tag1,
			Timestamps:   validationResponse.Timestamps,
			Status:       status,
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

func GetValidatorValidationList(validatorAddress string, page, pageSize int, filter string) ([]*serverTypes.ValidatorValidation, int64, error) {
	validations, total, err := model.GetValidationListByValidatorAddress(validatorAddress, page, pageSize, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("fail to get validator requests: %v", err)
	}

	var validationsList []*serverTypes.ValidatorValidation
	for _, validation := range validations {
		chainInfo, _ := config.GetChainInfo(validation.ChainID)
		registryAddress := config.GetIdentityAddressByValidationAddress(validation.ChainID, validation.ValidationRegistry)
		deployerInfo := config.GetContractsDeployerInfo(validation.ChainID, registryAddress)
		status := "pending"
		if len(strings.Trim(validation.ResponseTxHash, " ")) > 0 {
			status = "finished"
		}

		validationsList = append(validationsList, &serverTypes.ValidatorValidation{
			ChainName:          chainInfo.ChainName,
			ChainLogo:          chainInfo.ChainLogo,
			ContractDeployer:   deployerInfo.Deployer,
			AgentUID:           validation.AgentUID,
			ChainID:            validation.ChainID,
			AgentID:            validation.AgentID,
			AgentName:          validation.AgentName,
			ValidationRegistry: validation.ValidationRegistry,

			ValidatorAddress: validation.ValidatorAddress,

			RequestURI:  validation.RequestURI,
			RequestHash: strings.Trim(validation.RequestHash, " "),

			Response:     validation.Response,
			ResponseURI:  validation.ResponseURI,
			ResponseHash: strings.Trim(validation.ResponseHash, " "),
			Tag1:         validation.Tag1,
			Timestamps:   validation.Timestamps,
			Status:       status,
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
