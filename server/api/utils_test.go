package api

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEncodeFeedbackAuth(t *testing.T) {
	// 对应的 JSON 示例：
	// {
	//   "agentId": 1,
	//   "clientAddress": "0x10aa9ce20a1b03b18b4e2fd7b5d747add9a02030",
	//   "indexLimit": 1,
	//   "expiry": 1731196800,
	//   "chainId": 1,
	//   "identityRegistry": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7",
	//   "signerAddress": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7"
	// }
	agentID := big.NewInt(1)
	clientAddr := common.HexToAddress("0x10aa9ce20a1b03b18b4e2fd7b5d747add9a02030")
	indexLimit := uint64(1)
	expiry := big.NewInt(1731196800)
	chainID := big.NewInt(1)
	identityRegistry := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7")
	signerAddr := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb7")
	encodedFeedbackAuth, err := EncodeFeedbackAuth(agentID, clientAddr, indexLimit, expiry, chainID, identityRegistry, signerAddr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("encodedFeedbackAuth:", hex.EncodeToString(encodedFeedbackAuth))
}
