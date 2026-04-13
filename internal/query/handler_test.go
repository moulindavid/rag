package query

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type mockQueryService struct {
	resp *Response
	err  error
}

func (m *mockQueryService) Query(_ context.Context, _ *Request) (*Response, error) {
	return m.resp, m.err
}

func TestQueryHandler_Success(t *testing.T) {
	h := &Handler{service: &mockQueryService{resp: &Response{Answer: "42"}}}

	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(`{"question":"what is the answer?"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Query(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Code)
	}
	var resp Response
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Answer != "42" {
		t.Errorf("expected answer %q, got %q", "42", resp.Answer)
	}
}

func TestQueryHandler_InvalidBody(t *testing.T) {
	h := &Handler{service: &mockQueryService{}}

	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader("not json"))
	w := httptest.NewRecorder()

	h.Query(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestQueryHandler_EmptyQuestion(t *testing.T) {
	h := &Handler{service: &mockQueryService{}}

	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(`{"question":""}`))
	w := httptest.NewRecorder()

	h.Query(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestQueryHandler_ServiceError(t *testing.T) {
	h := &Handler{service: &mockQueryService{err: errors.New("boom")}}

	req := httptest.NewRequest(http.MethodPost, "/query", strings.NewReader(`{"question":"what?"}`))
	w := httptest.NewRecorder()

	h.Query(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
