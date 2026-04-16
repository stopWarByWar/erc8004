package logic

import (
	"fmt"
	"testing"
)

func TestGetNetworkList(t *testing.T) {
	initTest(t)
	networkList, err := GetNetworkList()
	if err != nil {
		t.Errorf("GetNetworkList error: %v", err)
	}
	fmt.Println(networkList)
}
