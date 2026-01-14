package logic

import (
	"agent_identity/server/api/types"
	"fmt"
	"testing"
)

func TestGetAgentFeedbacksList(t *testing.T) {
	initTest()
	agentUID := uint64(1900)
	page := 1
	pageSize := 10
	feedbacks, total, err := GetAgentFeedbacksList(agentUID, page, pageSize)
	if err != nil {
		t.Errorf("GetAgentFeedbacksList error: %v", err)
	}
	fmt.Println(feedbacks, total)
}

func TestSetFeedback(t *testing.T) {
	initTest()
	tag1 := "test"
	tag2 := "test"
	skill := "test"
	context := "test"
	domain := "test"
	name := "test"
	request := types.UploadFeedbackRequest{
		UID:           uint64(1),
		ClientAddress: "0x0004AA63c570c570eBF15376c0dB199918BFe9Fb",
		Skill:         &skill,
		Context:       &context,
		Domain:        &domain,
		Name:          &name,
		IndexLimit:    100,
		Score:         1,
		Tag1:          &tag1,
		Tag2:          &tag2,
	}
	feedbackURI, feedbackHash, err := SetFeedback(request)
	if err != nil {
		t.Errorf("SetFeedback error: %v", err)
	}
	fmt.Println(feedbackURI, feedbackHash)
}
