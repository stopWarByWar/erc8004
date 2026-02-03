package model

import (
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/gorm"
)

func GetLatestFeedbackAndResponse(chainID string, reputationRegistry string) (uint64, uint64, error) {
	var feedback *Feedback
	err := db.Where("chain_id = ? and (reputation_registry = ? or reputation_registry = ?)", chainID, common.HexToAddress(reputationRegistry).String(), common.HexToAddress(reputationRegistry).String()).Order("block_number DESC, index DESC").First(&feedback).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	var response *Response
	err = db.Where("chain_id = ? and (reputation_registry = ? or reputation_registry = ?)", chainID, common.HexToAddress(reputationRegistry).String(), common.HexToAddress(reputationRegistry).String()).Order("block_number DESC, index DESC").First(&response).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, 0, nil
	} else if err != nil {
		return 0, 0, err
	}

	if response.BlockNumber > feedback.BlockNumber {
		return response.BlockNumber, response.Index, nil
	} else if response.BlockNumber == feedback.BlockNumber && response.Index > feedback.Index {
		return response.BlockNumber, response.Index, nil
	}
	return feedback.BlockNumber, feedback.Index, nil
}

func CreateFeedback(feedback *Feedback) error {
	var amount int64
	err := db.Model(&Feedback{}).
		Where("chain_id = ? and agent_id = ? and reputation_registry = ? and client_address = ? and feedback_index = ?", feedback.ChainID, feedback.AgentID, feedback.ReputationRegistry, feedback.ClientAddress, feedback.FeedbackIndex).
		Count(&amount).Error
	if err != nil {
		return err
	}
	if amount > 0 {
		return nil
	}
	return db.Omit("uid").Create(feedback).Error
}

func UpdateFeedbackRevoked(chainID, agentID, reputationRegistry, clientAddress string, feedbackIndex uint64) error {
	return db.
		Model(&Feedback{}).
		Where("chain_id = ? and agent_id = ? and reputation_registry =? and client_address = ? and feedback_index = ?", chainID, agentID, reputationRegistry, clientAddress, feedbackIndex).
		Update("revoked", true).
		Update("endpoint", "").Error
}

func GetFeedbackUIDAndAgentUID(chainID, agentID, reputationRegistry, clientAddress string, feedbackIndex uint64) (uint64, uint64, error) {
	var feedback *Feedback
	err := db.
		Where("chain_id = ? and agent_id = ? and reputation_registry =? and client_address = ? and feedback_index = ?", chainID, agentID, reputationRegistry, clientAddress, feedbackIndex).
		First(&feedback).Error
	if err != nil {
		return 0, 0, err
	}
	return feedback.UID, feedback.AgentUID, nil
}

func CreateResponse(chainID string, response *Response) error {
	var amount int64
	err := db.Model(&Response{}).Where("chain_id = ? and block_number = ? and index = ?", chainID, response.BlockNumber, response.Index).Count(&amount).Error
	if err != nil || amount > 0 {
		return err
	}

	return db.Omit("uid").Create(&response).Error
}

//
//  ------------------- For API -------------------
//
//
//

func GetFeedbacksByAgentUID(uid uint64, tag1 string, page, pageSize int) ([]*FeedbackResp, int64, error) {
	if page <= 0 || pageSize <= 0 {
		return nil, 0, errors.New("invalid page or pageSize")
	}
	// 基础查询（不含分页），用于统计总数和复用条件
	baseQuery := db.Model(&Feedback{}).
		Where("agent_uid = ? and revoked = ?", uid, false)
	if tag1 != "" {
		baseQuery = baseQuery.Where("tag1 = ?", tag1)
	}

	var total int64
	if err := baseQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*FeedbackResp{}, 0, nil
	}

	// 分页查询当前页数据
	var feedbacks []*Feedback
	if err := baseQuery.
		Order("timestamps DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&feedbacks).Error; err != nil {
		return nil, 0, err
	}

	var feedbackAddrs []string
	for _, feedback := range feedbacks {
		feedbackAddrs = append(feedbackAddrs, feedback.ClientAddress)
	}

	passportMap, err := checkPassport(feedbackAddrs)
	if err != nil {
		return nil, 0, err
	}

	feedbackResp := make([]*FeedbackResp, 0)
	for _, feedback := range feedbacks {
		if account, ok := passportMap[feedback.ClientAddress]; ok {
			feedbackResp = append(feedbackResp, &FeedbackResp{
				Feedback: *feedback,
				Name:     account.Name,
				Avatar:   account.Avatar,
				Passport: true,
			})
		} else {
			feedbackResp = append(feedbackResp, &FeedbackResp{
				Feedback: *feedback,
				Passport: false,
			})
		}
	}
	return feedbackResp, total, nil
}

type ScoreInfo struct {
	Tag                              string
	Score                            float64
	UniqueFeedbackClientAddressCount uint64
	FeedbackCount                    uint64
}

func GetScoreForEachTag1(agentUID uint64, offset, limit int) ([]ScoreInfo, error) {
	if offset < 0 || limit <= 0 {
		return nil, errors.New("invalid offset or limit")
	}
	var scores []*FeedbackTagScore
	if err := db.Model(&FeedbackTagScore{}).
		Where("agent_uid = ?", agentUID).
		Order("unique_feedback_client_address_count DESC, feedback_count DESC").
		Offset(offset).
		Limit(limit).
		Scan(&scores).Error; err != nil {
		return nil, err
	}
	scoresInfo := make([]ScoreInfo, len(scores))
	for _, score := range scores {
		scoresInfo = append(scoresInfo, ScoreInfo{
			Tag:                              score.Tag,
			Score:                            score.Score,
			UniqueFeedbackClientAddressCount: score.UniqueFeedbackClientAddressCount,
			FeedbackCount:                    score.FeedbackCount,
		})
	}
	return scoresInfo, nil
}
