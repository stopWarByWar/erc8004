package handle

import (
	"agent_identity/logger"
	apiUtils "agent_identity/server/api/utils"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func TestGetCommerceActionsHandler_InvalidEnums_Return400(t *testing.T) {
	gin.SetMode(gin.TestMode)
	_, _ = logger.New(&logger.Config{Level: logrus.ErrorLevel, ReportCaller: false, FilePath: "./log/test"})
	apiUtils.Init(nil)

	r := gin.New()
	r.GET("/agent/commerce/actions", GetCommerceActionsHandler)

	cases := []string{
		"/agent/commerce/actions?uid=1&role=bad",
		"/agent/commerce/actions?uid=1&certainty=bad",
		"/agent/commerce/actions?uid=1&sort_by=bad",
		"/agent/commerce/actions?uid=1&action=bad_action",
		"/agent/commerce/actions?uid=1&polarity=bad_polarity",
		"/agent/commerce/actions?uid=1&action=job_completed,bad_action",
		"/agent/commerce/actions?uid=1&polarity=positive,bad_polarity",
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

