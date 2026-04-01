package logic

import (
	"agent_identity/model"
	serverTypes "agent_identity/server/api/types"
)

func GetLatestAgentValidationEval(agentUID uint64) (*serverTypes.AgentValidationEvalLatestResponse, error) {
	report, dims, err := model.GetLatestAgentValidationEvalByAgentUID(agentUID)
	if err != nil {
		return nil, err
	}
	out := &serverTypes.AgentValidationEvalLatestResponse{
		Dimensions: make([]serverTypes.AgentValidationEvalDimensionResponse, 0),
	}
	if report == nil {
		return out, nil
	}
	out.Report = &serverTypes.AgentValidationEvalReportResponse{
		ID:           report.ID,
		AgentUID:     report.AgentUID,
		AgentTokenID: report.AgentTokenID,
		ChainID:      report.ChainID,
		Reporter:     report.Reporter,
		Version:      report.Version,
		Score:        report.Score,
		ReportURL:    report.ReportURL,
		Desc:         report.Desc,
		ValidatedAt:  report.ValidatedAt,
		CreatedAt:    report.CreatedAt,
	}
	for _, d := range dims {
		out.Dimensions = append(out.Dimensions, serverTypes.AgentValidationEvalDimensionResponse{
			Dimension: d.Dimension,
			Score:     d.Score,
			CreatedAt: d.CreatedAt,
		})
	}
	return out, nil
}
