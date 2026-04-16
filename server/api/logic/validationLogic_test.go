package logic

import (
	"fmt"
	"testing"
)

func TestGetAgentValidationList(t *testing.T) {
	initTest(t)
	uid := uint64(1)
	page := 1
	pageSize := 10
	filter := "all"
	validationResponses, total, err := GetAgentValidationList(uid, page, pageSize, filter)
	if err != nil {
		t.Errorf("GetAgentValidationList error: %v", err)
	}
	fmt.Println(validationResponses, total)
}

func TestGetValidatorList(t *testing.T) {
	initTest(t)
	page := 1
	pageSize := 10
	validatorList, total, err := GetValidatorList(page, pageSize)
	if err != nil {
		t.Errorf("GetValidatorList error: %v", err)
	}
	fmt.Println(validatorList, total)
}

func TestGetValidatorValidationList(t *testing.T) {
	initTest(t)
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
	initTest(t)
	validatorAddress := "0x0004AA63c570c570eBF15376c0dB199918BFe9Fb"
	validator, err := GetValidatorByAddress(validatorAddress)
	if err != nil {
		t.Errorf("GetValidatorByAddress error: %v", err)
	}
	fmt.Println(validator)
}
