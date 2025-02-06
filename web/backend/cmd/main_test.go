package main

import (
	"bytes"
	"net/http/httptest"
	"testing"
)

func TestCreateParecer(t *testing.T) {
	req := httptest.NewRequest("POST", "/parecer", bytes.NewBufferString(`{"user":"user","creci":"creci","content":"content"}`))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	CreateParecerHandler(w, req)

	if w.Code != 201 {
		t.Errorf("expected status code 201, got %d", w.Code)
		return
	}
}

func TestGetParecer(t *testing.T) {
	req := httptest.NewRequest("GET", "/parecer?id=parecer-1-01-0001", nil)
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	GetParecerHandler(w, req)

	if w.Code != 200 {
		t.Errorf("expected status code 200, got %d", w.Code)
		return
	}
}
