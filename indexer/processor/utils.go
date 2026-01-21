package processor

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
)

type Comment struct {
	IdentityRegistry string
	Commenter        common.Address
	AgentID          *big.Int
	Score            uint16
	CommentText      string
}

var CommentSchema = `tuple(uint256 agentId,address commenter, uint16 score,string comment)`

func DecodeCommentEvent(reputationRegistry string, data []byte) (*Comment, error) {
	typ := abi.MustNewType(CommentSchema)

	res, err := typ.Decode(data)
	if err != nil {
		return nil, fmt.Errorf("fail to decode Comment with error %v", err)
	}

	result, ok := res.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert decode result to map")
	}

	comment := &Comment{}

	comment.IdentityRegistry = reputationRegistry

	comment.AgentID, ok = result["agentId"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("failed to convert agentClientId to *big.Int")
	}

	commenter, ok := result["commenter"].(ethgo.Address)
	if !ok {
		return nil, fmt.Errorf("failed to convert commenter to ethgo.Address")
	}
	comment.Commenter = common.Address(commenter)

	comment.Score, ok = result["score"].(uint16)
	if !ok {
		return nil, fmt.Errorf("failed to convert score to uint16")
	}

	comment.CommentText, ok = result["comment"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to convert comment to string")
	}

	return comment, nil
}

func calculateScore(value *big.Int, valueDecimals uint8) uint8 {
	if value == nil {
		return 0
	}

	// 步骤1：计算 realValue = value / 10^d （即value右移d位）
	denominator := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(valueDecimals)), nil)
	realValue := new(big.Int).Div(value, denominator) // 整数除法（舍去小数，如需四舍五入则保留之前的halfDenominator逻辑）

	// 步骤2：限制结果范围在0-100之间
	switch {
	case realValue.Sign() < 0:
		return 0
	case realValue.Cmp(big.NewInt(100)) > 0:
		return 100
	default:
		return uint8(realValue.Uint64())
	}
}
