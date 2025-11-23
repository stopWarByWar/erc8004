package api

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
)

func TestEncodeFeedbackAuth(t *testing.T) {
	// {
	// 	"UID": 54,
	// 	"score": 79,
	// 	"tag1": "123",
	// 	"tag2": "456",
	// 	"context": "asdasdasd",
	// 	"feedbackAuth": {
	// 		"agentId": "323",
	// 		"clientAddress": "0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345",
	// 		"indexLimit": 1,
	// 		"expiry": 1763901876,
	// 		"chainId": "97",
	// 		"identityRegistry": "0xA98A5542a1AaB336397d487e32021E0E48BEF717",
	// 		"signerAddress": "0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345"
	// 	},
	// 	"feedbackAuthSignature": "0x072eb93d9d6ee24174bae4c973f5f7d8eeb57919831e288a46e3ec1e6f11f044242efd0a9dbdd5943c3b389ac25c5d70736430bedd120166e1b4cc4b93ed91811b"
	// }
	agentID := big.NewInt(323)
	clientAddr := common.HexToAddress("0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345")
	indexLimit := uint64(1)
	expiry := big.NewInt(1763901876)
	chainID := big.NewInt(97)
	identityRegistry := common.HexToAddress("0xA98A5542a1AaB336397d487e32021E0E48BEF717")
	signerAddr := common.HexToAddress("0x5C33f9bAFcC7e1347937e0E986Ee14e84A6DF345")
	encodedFeedbackAuth, err := EncodeFeedbackAuth(agentID, clientAddr, indexLimit, expiry, chainID, identityRegistry, signerAddr)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println("encodedFeedbackAuth:", hex.EncodeToString(encodedFeedbackAuth))

	signature := "0x072eb93d9d6ee24174bae4c973f5f7d8eeb57919831e288a46e3ec1e6f11f044242efd0a9dbdd5943c3b389ac25c5d70736430bedd120166e1b4cc4b93ed91811b"
	signerAddress, err := ExtractAddressFromSignature(signature, encodedFeedbackAuth)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("signerAddress:", signerAddress.String())
}
