package processor

import (
	"errors"
	"testing"
)

func TestRetryableError(t *testing.T) {
	orig := errors.New("boom")
	err := retryable(orig)
	if err == nil {
		t.Fatal("expected non-nil error")
	}
	if !isRetryable(err) {
		t.Fatal("expected retryable error")
	}
	if !errors.Is(err, orig) {
		t.Fatal("expected wrapped error to match with errors.Is")
	}
}

