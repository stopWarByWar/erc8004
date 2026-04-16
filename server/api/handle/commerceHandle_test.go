package handle

import (
	agentLogger "agent_identity/logger"
	"agent_identity/model"
	"agent_identity/logger"
	apiUtils "agent_identity/server/api/utils"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func TestGetCommerceActionsHandler_InvalidEnums_Return400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, _ = logger.New(&logger.Config{Level: logrus.ErrorLevel, ReportCaller: false, FilePath: "./log/test"})
	apiUtils.Init(nil, apiUtils.RuntimeConfig{})

	r := gin.New()
	r.GET("/agent/identity/commerce/actions", GetCommerceActionsHandler)

	cases := []string{
		"/agent/identity/commerce/actions?uid=1&role=bad",
		"/agent/identity/commerce/actions?uid=1&certainty=bad",
		"/agent/identity/commerce/actions?uid=1&sort_by=bad",
		"/agent/identity/commerce/actions?uid=1&action=bad_action",
		"/agent/identity/commerce/actions?uid=1&polarity=bad_polarity",
		"/agent/identity/commerce/actions?uid=1&action=job_completed,bad_action",
		"/agent/identity/commerce/actions?uid=1&polarity=positive,bad_polarity",
	}

	for _, path := range cases {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		// Project's ErrResp always returns HTTP 200 with code/status fields.
		if w.Code != http.StatusOK {
			t.Fatalf("path=%s: expected 200, got %d body=%s", path, w.Code, w.Body.String())
		}
		body := w.Body.String()
		if body == "" || !containsAll(body, []string{`"code":1`, `"status":"failed"`, `"msg":"Invalid Request"`}) {
			t.Fatalf("path=%s: expected Invalid Request body, got %s", path, body)
		}
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !strings.Contains(s, sub) {
			return false
		}
	}
	return true
}

type testLogicConfig struct {
	Dns          string `yaml:"dns"`
	OpenaiAPIKey string `yaml:"openai_api_key"`
}

func initHandleTest(t *testing.T) {
	t.Helper()
	if os.Getenv("RUN_DB_TESTS") != "1" {
		t.Skip("set RUN_DB_TESTS=1 to run DB-backed tests")
	}
	gin.SetMode(gin.TestMode)

	// Init logger for ErrResp paths that include errorInfos (would call _logger).
	lg, err := agentLogger.New(&agentLogger.Config{Level: logrus.ErrorLevel, ReportCaller: false, FilePath: "./log/test"})
	if err == nil {
		apiUtils.Init(lg, apiUtils.RuntimeConfig{})
	} else {
		apiUtils.Init(nil, apiUtils.RuntimeConfig{})
	}

	// Load DB dns from server/api/logic/config.yaml (allowed by user).
	data, err := os.ReadFile("../logic/config.yaml")
	if err != nil {
		t.Fatalf("read ../logic/config.yaml: %v", err)
	}
	var cfg testLogicConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		t.Fatalf("unmarshal config.yaml: %v", err)
	}
	if cfg.Dns == "" {
		t.Fatalf("empty dns in ../logic/config.yaml")
	}

	model.InitDB(cfg.Dns, cfg.OpenaiAPIKey)
}

type apiResp struct {
	Code   int             `json:"code"`
	Status string          `json:"status"`
	Msg    string          `json:"msg"`
	Data   json.RawMessage `json:"data"`
}

func decodeResp(t *testing.T, body string) apiResp {
	t.Helper()
	var r apiResp
	if err := json.Unmarshal([]byte(body), &r); err != nil {
		t.Fatalf("decode resp: %v body=%s", err, body)
	}
	return r
}

func initMockHandleTest(t *testing.T) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	lg, _ := agentLogger.New(&agentLogger.Config{Level: logrus.ErrorLevel, ReportCaller: false, FilePath: "./log/test"})
	apiUtils.Init(lg, apiUtils.RuntimeConfig{Mock: true})
}

func assertMockHeader(t *testing.T, w *httptest.ResponseRecorder) {
	t.Helper()
	if w.Header().Get("X-Mock") != "1" {
		t.Fatalf("expected X-Mock=1 header, got %q body=%s", w.Header().Get("X-Mock"), w.Body.String())
	}
}

func TestCommerceHandlers_MockMode_HeaderAndShape(t *testing.T) {
	initMockHandleTest(t)

	r := gin.New()
	r.GET("/agent/identity/commerce/scores", GetCommerceScoresHandler)
	r.GET("/agent/identity/commerce/actions", GetCommerceActionsHandler)
	r.GET("/agent/identity/commerce/stats", GetCommerceStatsHandler)
	r.GET("/agent/commerce/jobs", GetCommerceJobsHandler)
	r.GET("/agent/commerce/jobs/detail", GetCommerceJobDetailHandler)
	r.GET("/agent/commerce/jobs/general", GetCommerceJobsGeneralHandler)
	r.GET("/agent/commerce/jobs/charts", GetCommerceJobsChartsHandler)
	r.GET("/agent/commerce/jobs/actions", GetCommerceJobActionsHandler)
	r.GET("/agent/commerce/jobs/filters", GetCommerceJobsFiltersHandler)

	cases := []struct {
		path         string
		expectSubstr []string
	}{
		{"/agent/identity/commerce/scores?uid=1", []string{`"scores"`}},
		{"/agent/identity/commerce/actions?uid=1&page=1&page_size=3", []string{`"actions"`, `"total"`}},
		{"/agent/identity/commerce/stats?uid=1", []string{`"stats"`}},
		{"/agent/commerce/jobs?page=1&page_size=5", []string{`"jobs"`, `"total"`}},
		{"/agent/commerce/jobs/detail?chain_id=1&commerce_contract=0x1111111111111111111111111111111111111111&job_id=10001", []string{`"job"`, `"evidence"`}},
		{"/agent/commerce/jobs/general", []string{`"summary"`, `"distributions"`, `"chain_contracts"`, `"erc8183_contracts"`}},
		{"/agent/commerce/jobs/charts", []string{`"charts"`, `"activity"`, `"volume"`, `"paid_volume_usd_over_time"`, `"value_usd"`, `"drilldown"`, `"chain_contracts"`, `"erc8183_contracts"`}},
		{"/agent/commerce/jobs/actions?chain_id=1&commerce_contract=0x1111111111111111111111111111111111111111&job_id=10001&page=1&page_size=5", []string{`"actions"`, `"total"`}},
		{"/agent/commerce/jobs/filters", []string{`"filters"`, `"chains"`, `"commerce_contracts"`, `"payment_tokens"`, `"last_updated"`}},
	}

	for _, tc := range cases {
		req := httptest.NewRequest(http.MethodGet, tc.path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("path=%s: expected 200, got %d body=%s", tc.path, w.Code, w.Body.String())
		}
		assertMockHeader(t, w)
		resp := decodeResp(t, w.Body.String())
		if resp.Code != 0 || resp.Status != "success" {
			t.Fatalf("path=%s: expected success, got code=%d status=%s msg=%s body=%s", tc.path, resp.Code, resp.Status, resp.Msg, w.Body.String())
		}
		for _, sub := range tc.expectSubstr {
			if !strings.Contains(w.Body.String(), sub) {
				t.Fatalf("path=%s: expected substring %q, got body=%s", tc.path, sub, w.Body.String())
			}
		}
	}
}

func TestCommerceJobs_List_Success(t *testing.T) {
	initHandleTest(t)

	r := gin.New()
	r.GET("/agent/commerce/jobs", GetCommerceJobsHandler)

	req := httptest.NewRequest(http.MethodGet, "/agent/commerce/jobs?page=1&page_size=5", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	resp := decodeResp(t, w.Body.String())
	if resp.Code != 0 || resp.Status != "success" {
		t.Fatalf("expected success, got code=%d status=%s msg=%s body=%s", resp.Code, resp.Status, resp.Msg, w.Body.String())
	}
	// Ensure data contains jobs/total keys (structure check, no assumption about non-empty).
	if !strings.Contains(w.Body.String(), `"jobs"`) || !strings.Contains(w.Body.String(), `"total"`) {
		t.Fatalf("expected jobs+total in data, got %s", w.Body.String())
	}
}

func TestCommerceJobs_Detail_InvalidRequest(t *testing.T) {
	initHandleTest(t)

	r := gin.New()
	r.GET("/agent/commerce/jobs/detail", GetCommerceJobDetailHandler)

	cases := []string{
		"/agent/commerce/jobs/detail?job_id=1",                                   // missing chain/contract
		"/agent/commerce/jobs/detail?chain_id=1&commerce_contract=0x&job_id=bad", // invalid job_id
	}
	for _, path := range cases {
		req := httptest.NewRequest(http.MethodGet, path, nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if w.Code != http.StatusOK {
			t.Fatalf("path=%s: expected 200, got %d body=%s", path, w.Code, w.Body.String())
		}
		resp := decodeResp(t, w.Body.String())
		if resp.Code != 1 || resp.Status != "failed" || resp.Msg != "Invalid Request" {
			t.Fatalf("path=%s: expected Invalid Request, got code=%d status=%s msg=%s body=%s", path, resp.Code, resp.Status, resp.Msg, w.Body.String())
		}
	}
}

func TestCommerceJobs_Detail_FromList_IfAny(t *testing.T) {
	initHandleTest(t)

	// Fetch a job via model, then call detail handler.
	jobs, total, err := model.GetCommerceJobs(model.CommerceJobsQuery{Page: 1, PageSize: 1})
	if err != nil {
		t.Fatalf("seed read jobs: %v", err)
	}
	if total == 0 || len(jobs) == 0 {
		t.Skip("no commerce_jobs data; skip detail success test")
	}
	j := jobs[0]

	r := gin.New()
	r.GET("/agent/commerce/jobs/detail", GetCommerceJobDetailHandler)

	path := "/agent/commerce/jobs/detail?chain_id=" + j.ChainID + "&commerce_contract=" + j.CommerceContract + "&job_id=1"
	// use the actual job_id
	path = "/agent/commerce/jobs/detail?chain_id=" + j.ChainID + "&commerce_contract=" + j.CommerceContract + "&job_id=" + func() string {
		b, _ := json.Marshal(j.JobID)
		// json marshals uint64 as number; convert without quotes by trimming spaces
		return strings.TrimSpace(string(b))
	}()

	req := httptest.NewRequest(http.MethodGet, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	resp := decodeResp(t, w.Body.String())
	if resp.Code != 0 || resp.Status != "success" {
		t.Fatalf("expected success, got code=%d status=%s msg=%s body=%s", resp.Code, resp.Status, resp.Msg, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"job"`) {
		t.Fatalf("expected job in data, got %s", w.Body.String())
	}
}

func TestCommerceJobs_General_Success(t *testing.T) {
	initHandleTest(t)

	r := gin.New()
	r.GET("/agent/commerce/jobs/general", GetCommerceJobsGeneralHandler)

	req := httptest.NewRequest(http.MethodGet, "/agent/commerce/jobs/general", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	resp := decodeResp(t, w.Body.String())
	if resp.Code != 0 || resp.Status != "success" {
		t.Fatalf("expected success, got code=%d status=%s msg=%s body=%s", resp.Code, resp.Status, resp.Msg, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"summary"`) || !strings.Contains(w.Body.String(), `"distributions"`) {
		t.Fatalf("expected summary+distributions in data, got %s", w.Body.String())
	}
}

func TestCommerceJobs_Charts_Success(t *testing.T) {
	initHandleTest(t)

	r := gin.New()
	r.GET("/agent/commerce/jobs/charts", GetCommerceJobsChartsHandler)

	req := httptest.NewRequest(http.MethodGet, "/agent/commerce/jobs/charts", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d body=%s", w.Code, w.Body.String())
	}
	resp := decodeResp(t, w.Body.String())
	if resp.Code != 0 || resp.Status != "success" {
		t.Fatalf("expected success, got code=%d status=%s msg=%s body=%s", resp.Code, resp.Status, resp.Msg, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), `"charts"`) {
		t.Fatalf("expected charts in data, got %s", w.Body.String())
	}
}
