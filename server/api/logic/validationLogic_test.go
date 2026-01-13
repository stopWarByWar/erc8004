package logic

import (
	"fmt"
	"testing"
)

func TestGetAgentValidationList(t *testing.T) {
	initTest()
	chainID := "11155111"
	validationRegistry := "0x8004AA63c570c570eBF15376c0dB199918BFe9Fb"
	agentID := "1"
	page := 1
	pageSize := 10
	filter := "all"
	validationResponses, total, err := GetAgentValidationList(chainID, validationRegistry, agentID, page, pageSize, filter)
	if err != nil {
		t.Errorf("GetAgentValidationList error: %v", err)
	}
	fmt.Println(validationResponses, total)
}

func TestGetValidatorList(t *testing.T) {
	initTest()
	page := 1
	pageSize := 10
	validatorList, total, err := GetValidatorList(page, pageSize)
	if err != nil {
		t.Errorf("GetValidatorList error: %v", err)
	}
	fmt.Println(validatorList, total)
}

func TestGetValidatorValidationList(t *testing.T) {
	initTest()
	validatorAddress := "0x0004AA63c570c570eBF15376c0dB199918BFe9Fb"
	page := 1
	pageSize := 10
	filter := "all"
	validations, total, err := GetValidatorValidationList(validatorAddress, page, pageSize, filter)
	if err != nil {
		t.Errorf("GetValidatorValidationList error: %v", err)
	}
	fmt.Println(validations, total)
}

func TestGetValidatorByAddress(t *testing.T) {
	initTest()
	validatorAddress := "0x0004AA63c570c570eBF15376c0dB199918BFe9Fb"
	validator, err := GetValidatorByAddress(validatorAddress)
	if err != nil {
		t.Errorf("GetValidatorByAddress error: %v", err)
	}
	fmt.Println(validator)
}
