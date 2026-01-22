package logic

import (
	apiUtils "agent_identity/server/api/utils"
	"fmt"
	"testing"
)

func TestGetNetworkList(t *testing.T) {
	initTest()
	networkList, err := GetNetworkList()
	if err != nil {
		t.Errorf("GetNetworkList error: %v", err)
	}
	fmt.Println(networkList)
}

func TestSetChainAgentAmount(t *testing.T) {
	initTest()
	apiUtils.UpdateGeneralInfo()
	info := apiUtils.GetGeneralInfo()
	fmt.Println(info)
}
