package logic

import (
	"agent_identity/server/api/types"
	"fmt"
	"testing"
)

func TestGetAgentFeedbacksList(t *testing.T) {
	initTest(t)
	agentUID := uint64(19838)
	page := 1
	pageSize := 10
	feedbacks, total, err := GetAgentFeedbacksList(agentUID, []string{"nice", "Good"}, page, pageSize)
	if err != nil {
		t.Errorf("GetAgentFeedbacksList error: %v", err)
	}
	fmt.Println(feedbacks, total)
}

func TestSetFeedback(t *testing.T) {
	initTest(t)
	tag1 := "tag1"
	tag2 := "tag2"
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

func TestGetAgentScoreForEachTag1(t *testing.T) {
	initTest(t)
	agentUID := uint64(19838)
	offset := 0
	limit := 10
	scores, err := GetAgentScoreForEachTag1(agentUID, offset, limit)
	if err != nil {
		t.Errorf("GetAgentScoreForEachTag1 error: %v", err)
	}
	fmt.Println(scores)
}
