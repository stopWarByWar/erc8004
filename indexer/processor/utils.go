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
