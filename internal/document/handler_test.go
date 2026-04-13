package document

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

type mockIngester struct {
	doc *Document
	err error
}

func (m *mockIngester) Ingest(_ context.Context, _ string, _ io.Reader) (*Document, error) {
	return m.doc, m.err
}

func makeUploadRequest(t *testing.T, filename string, content string) *http.Request {
	t.Helper()
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Fprint(fw, content)
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "/documents", &buf)
	req.Header.Set("Content-Type", w.FormDataContentType())
	return req
}

func TestUploadHandler_Success(t *testing.T) {
	doc := &Document{ID: uuid.New(), Filename: "test.txt"}
	h := &Handler{service: &mockIngester{doc: doc}}

	req := makeUploadRequest(t, "test.txt", "hello world")
	w := httptest.NewRecorder()

	h.Upload(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected 201, got %d", w.Code)
	}
}

func TestUploadHandler_MissingFile(t *testing.T) {
	h := &Handler{service: &mockIngester{}}

	// valid multipart body but no "file" field
	var buf bytes.Buffer
	mw := multipart.NewWriter(&buf)
	mw.WriteField("other", "value")
	mw.Close()

	req := httptest.NewRequest(http.MethodPost, "/documents", &buf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	w := httptest.NewRecorder()

	h.Upload(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUploadHandler_UnsupportedExtension(t *testing.T) {
	h := &Handler{service: &mockIngester{
		err: fmt.Errorf("parse document: unsupported file extension: .csv"),
	}}

	req := makeUploadRequest(t, "data.csv", "a,b,c")
	w := httptest.NewRecorder()

	h.Upload(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestUploadHandler_ServiceError(t *testing.T) {
	h := &Handler{service: &mockIngester{err: errors.New("db down")}}

	req := makeUploadRequest(t, "test.txt", "hello")
	w := httptest.NewRecorder()

	h.Upload(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}
