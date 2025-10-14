package processor

import (
	"fmt"
	"math/big"
	"sort"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/umbracle/ethgo/abi"
)

type Comment struct {
	Commenter     *common.Address
	AgentClientId *big.Int
	AgentServerId *big.Int
	Score         uint16
	Comment       string
	IsAuthorized  bool
}

// CommentDataABI 定义 CommentData 对应的 ABI
var CommentSchema = `tuple(address commenter,uint256 agentClientId,uint256 agentServerId,uint16 score,string comment,bool isAuthorized)`

func DecodeCommentEvent(data []byte) (*Comment, error) {
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
	comment.Commenter, ok = result["commenter"].(*common.Address)
	if !ok {
		return nil, fmt.Errorf("failed to convert commenter to common.Address")
	}
	comment.AgentClientId, ok = result["agentClientId"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("failed to convert agentClientId to *big.Int")
	}
	comment.AgentServerId, ok = result["agentServerId"].(*big.Int)
	if !ok {
		return nil, fmt.Errorf("failed to convert agentServerId to *big.Int")
	}
	comment.Score, ok = result["score"].(uint16)
	if !ok {
		return nil, fmt.Errorf("failed to convert score to uint16")
	}
	comment.Comment, ok = result["comment"].(string)
	if !ok {
		return nil, fmt.Errorf("failed to convert comment to string")
	}

	comment.IsAuthorized, ok = result["isAuthorized"].(bool)
	if !ok {
		return nil, fmt.Errorf("failed to convert isAuthorized to bool")
	}

	return comment, nil
}

func sortEvents(eventA []types.Log, eventB []types.Log) []types.Log {
	events := append(eventA, eventB...)
	sort.Slice(events, func(i, j int) bool {
		return events[i].BlockNumber < events[j].BlockNumber || (events[i].BlockNumber == events[j].BlockNumber && events[i].Index < events[j].Index)
	})
	return events
}
