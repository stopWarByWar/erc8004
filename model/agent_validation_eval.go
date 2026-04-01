package model

import (
	"errors"
	"time"

	"gorm.io/gorm"
)

type AgentValidationReport struct {
	ID           uint      `gorm:"column:id;primaryKey"`
	AgentUID     uint64    `gorm:"column:agent_uid;not null"`
	AgentTokenID int       `gorm:"column:agent_token_id;not null"`
	ChainID      int       `gorm:"column:chainid;not null"`
	Reporter     string    `gorm:"column:reporter;type:varchar(128);not null;default:bas-agent-eval"`
	Version      string    `gorm:"column:version;type:varchar(64);not null;default:1.0"`
	Score        float64   `gorm:"column:score;type:double precision;not null"`
	ReportURL    string    `gorm:"column:report_url;type:text"`
	Desc         string    `gorm:"column:desc;type:text"`
	ValidatedAt  time.Time `gorm:"column:validated_at;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
}

func (AgentValidationReport) TableName() string { return "agent_validation_reports" }

type AgentValidationDimension struct {
	ID           uint      `gorm:"column:id;primaryKey"`
	ReportID     uint      `gorm:"column:report_id;not null"`
	AgentUID     uint64    `gorm:"column:agent_uid;not null"`
	AgentTokenID int       `gorm:"column:agent_token_id;not null"`
	ChainID      int       `gorm:"column:chainid;not null"`
	Dimension    string    `gorm:"column:dimension;type:varchar(64);not null"`
	Score        float64   `gorm:"column:score;type:double precision;not null"`
	CreatedAt    time.Time `gorm:"column:created_at;not null"`
}

func (AgentValidationDimension) TableName() string { return "agent_validation_dimensions" }

// GetLatestAgentValidationEvalByAgentUID returns the newest report by validated_at (then id) and its dimension rows.
func GetLatestAgentValidationEvalByAgentUID(uid uint64) (*AgentValidationReport, []AgentValidationDimension, error) {
	var report AgentValidationReport
	err := db.Where("agent_uid = ?", uid).Order("validated_at DESC, id DESC").First(&report).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	var dims []AgentValidationDimension
	if err := db.Where("report_id = ?", report.ID).Order("dimension ASC").Find(&dims).Error; err != nil {
		return nil, nil, err
	}
	return &report, dims, nil
}
