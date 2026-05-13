package handle

import (
	"agent_identity/logger"
	"agent_identity/model"
	apiUtils "agent_identity/server/api/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

// Phase 4 — handler-level tests for GET /agent/feedback/credit
//
// The handler must:
//   - Reject missing / non-numeric uid with Invalid Request (200 envelope).
//   - On valid uid, return SuccessResp envelope containing the orchestrator
//     response under "credit".
//   - Cold-start agents must serialize with "credit_score": null.

func setupHandlerTestMode(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	lg, err := logger.New(&logger.Config{Level: logrus.ErrorLevel, ReportCaller: false, FilePath: "./log/test"})
	if err == nil {
		apiUtils.Init(lg, apiUtils.RuntimeConfig{})
	} else {
		apiUtils.Init(nil, apiUtils.RuntimeConfig{})
	}
}

func initFeedbackCreditHandleTest(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run DB-backed handler tests")
	}
	setupHandlerTestMode(t)

	data, err := os.ReadFile("../logic/config.yaml")
	if err != nil {
		t.Fatalf("read ../logic/config.yaml: %v", err)
	}
	var cfg struct {
		Dns          string `yaml:"dns"`
		OpenaiAPIKey string `yaml:"openai_api_key"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	model.InitDB(cfg.Dns, cfg.OpenaiAPIKey)
}

// ─── Validation: bad uid ───────────────────────────────────────────────────

func TestGetFeedbackCreditHandler_InvalidUID_Return400(t *testing.T) {
	setupHandlerTestMode(t)
	r := gin.New()
	r.GET("/agent/feedback/credit", GetFeedbackCreditHandler)

	for _, path := range []string{
		"/agent/feedback/credit",            // missing uid
		"/agent/feedback/credit?uid=",       // empty uid
		"/agent/feedback/credit?uid=abc",    // non-numeric
		"/agent/feedback/credit?uid=-1",     // negative
	} {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("path=%s: expected 200 envelope, got %d body=%s", path, w.Code, w.Body.String())
		}
		body := w.Body.String()
		if !containsAll(body, []string{`"code":1`, `"status":"failed"`, `"msg":"Invalid Request"`}) {
			t.Fatalf("path=%s: expected Invalid Request body, got %s", path, body)
		}
	}
}

// ─── Cold start agent: credit_score is null ─────────────────────────────

func TestGetFeedbackCreditHandler_ColdStart_NullCreditScore(t *testing.T) {
	initFeedbackCreditHandleTest(t)
	const aid uint64 = 9200001
	_ = model.DeleteFeedbacksForTest(aid)
	t.Cleanup(func() { _ = model.DeleteFeedbacksForTest(aid) })

	r := gin.New()
	r.GET("/agent/feedback/credit", GetFeedbackCreditHandler)

	req := httptest.NewRequest(http.MethodGet, "/agent/feedback/credit?uid=9200001", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}

	var env struct {
		Code   int             `json:"code"`
		Status string          `json:"status"`
		Data   json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &env); err != nil {
		t.Fatalf("unmarshal envelope: %v body=%s", err, w.Body.String())
	}
	if env.Code != 0 || env.Status != "success" {
		t.Fatalf("expected success envelope, got code=%d status=%s body=%s",
			env.Code, env.Status, w.Body.String())
	}

	var data struct {
		Credit struct {
			AgentUID    uint64   `json:"agent_uid"`
			CreditScore *float64 `json:"credit_score"`
			TagCount    int      `json:"tag_count"`
		} `json:"credit"`
	}
	if err := json.Unmarshal(env.Data, &data); err != nil {
		t.Fatalf("unmarshal data: %v raw=%s", err, env.Data)
	}
	if data.Credit.AgentUID != aid {
		t.Errorf("agent_uid=%d want %d", data.Credit.AgentUID, aid)
	}
	if data.Credit.CreditScore != nil {
		t.Errorf("credit_score=%v want null", *data.Credit.CreditScore)
	}
	if data.Credit.TagCount != 0 {
		t.Errorf("tag_count=%d want 0", data.Credit.TagCount)
	}
}
