package processor

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
)

func Test_DecodeCommentEvent(t *testing.T) {
	// Keep a legacy fixture around for quick regression checks (not used for decode).
	_, _ = hex.DecodeString("00000000000000000000000010aa9ce20a1b03b18b4e2fd7b5d747add9a0203000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000c7465737420636f6d6d656e740000000000000000000000000000000000000000")

	typ := abi.MustNewType(CommentSchema)
	addr := common.HexToAddress("0x10aa9ce20a1b03b18b4e2fd7b5d747add9a02030")
	payload, err := typ.Encode(map[string]any{
		"agentId":    big.NewInt(3),
		"commenter":  ethgo.Address(addr),
		"score":      uint16(100),
		"comment":    "test comment",
	})
	if err != nil {
		t.Fatal(err)
	}

	comment, err := DecodeCommentEvent(addr.String(), payload)
	if err != nil {
		t.Fatal(err)
	}
	if comment.Commenter != addr {
		t.Fatalf("commenter mismatch: %s != %s", comment.Commenter, addr)
	}
	if comment.AgentID.Cmp(big.NewInt(3)) != 0 {
		t.Fatalf("agentId mismatch: %s", comment.AgentID.String())
	}
	if comment.Score != 100 {
		t.Fatalf("score mismatch: %d", comment.Score)
	}
	if comment.CommentText != "test comment" {
		t.Fatalf("comment mismatch: %q", comment.CommentText)
	}

}

func Test_CalculateScore(t *testing.T) {
	value := big.NewInt(10000000)
	valueDecimals := uint8(2)
	score, err := calculateScore(value, valueDecimals)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(score)
}
