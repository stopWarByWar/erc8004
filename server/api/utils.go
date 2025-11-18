package api

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// EncodeFeedbackAuth 将 FeedbackAuth 各字段按 solidity ABI 顺序编码成 []byte。
func EncodeFeedbackAuth(
	agentID *big.Int,
	clientAddr common.Address,
	indexLimit uint64,
	expiry *big.Int,
	chainID *big.Int,
	identityRegistry common.Address,
	signerAddr common.Address,
) ([]byte, error) {
	uint256Type, err := abi.NewType("uint256", "", nil)
	if err != nil {
		return nil, fmt.Errorf("init uint256 type failed: %w", err)
	}
	addressType, err := abi.NewType("address", "", nil)
	if err != nil {
		return nil, fmt.Errorf("init address type failed: %w", err)
	}
	uint64Type, err := abi.NewType("uint64", "", nil)
	if err != nil {
		return nil, fmt.Errorf("init uint64 type failed: %w", err)
	}

	arguments := abi.Arguments{
		{Type: uint256Type}, // uint256 agentId
		{Type: addressType}, // address clientAddress
		{Type: uint64Type},  // uint64 indexLimit
		{Type: uint256Type}, // uint256 expiry
		{Type: uint256Type}, // uint256 chainId
		{Type: addressType}, // address identityRegistry
		{Type: addressType}, // address signerAddress
	}

	values := []interface{}{
		agentID,
		clientAddr,
		indexLimit,
		expiry,
		chainID,
		identityRegistry,
		signerAddr,
	}

	data, err := arguments.Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("encode feedback auth failed: %w", err)
	}
	if len(data) != 224 {
		return nil, fmt.Errorf("unexpected length %d, expect 224", len(data))
	}
	return data, nil
}

func ExtractAddressFromSignature(feedbackAuthSignature string, encodedFeedbackAuth []byte) (common.Address, error) {
	if strings.HasPrefix(feedbackAuthSignature, "0x") {
		feedbackAuthSignature = feedbackAuthSignature[2:]
	}
	sigBytes, err := hex.DecodeString(feedbackAuthSignature)
	if err != nil || len(sigBytes) != 65 {
		return common.Address{}, fmt.Errorf("invalid feedbackAuthSignature: %w", err)
	}

	// get ethereum hash (keccak256)
	msgHash := crypto.Keccak256(encodedFeedbackAuth)

	// EIP-191:\x19Ethereum Signed Message:\n长度+内容
	ethMsgHash := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		msgHash,
	)

	// v值修正
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	recoveredPubKey, err := crypto.SigToPub(ethMsgHash, sigBytes)
	if err != nil {
		return common.Address{}, fmt.Errorf("fail to recover EOA: %w", err)
	}
	recoveredAddr := crypto.PubkeyToAddress(*recoveredPubKey).String()
	if recoveredAddr != feedbackAuthSignature {
		return common.Address{}, fmt.Errorf("invalid signer address: %s", recoveredAddr)
	}
	return crypto.PubkeyToAddress(*recoveredPubKey), nil
}
