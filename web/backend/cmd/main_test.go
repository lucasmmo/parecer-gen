package main

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"parecer-gen/pkg/parecer"
	"testing"
)

func TestHandleParecer(t *testing.T) {
	t.Run("should create parecer without date", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/parecer", bytes.NewBufferString(`{"user":"user","creci":"creci","content":"content"}`))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		CreateParecer(w, req)

		if w.Code != 201 {
			t.Errorf("expected status code 201, got %d", w.Code)
			return
		}
	})

	t.Run("should create parecer with date", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/parecer", bytes.NewBufferString(`{"user":"user","creci":"creci","content":"content","date":"2025-02-07"}`))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		CreateParecer(w, req)

		if w.Code != 201 {
			t.Errorf("expected status code 201, got %d", w.Code)
			return
		}
	})

	t.Run("should not create parecer with missing data", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/parecer", bytes.NewBufferString(`{"user":"user","creci":"creci","content":""}`))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()

		CreateParecer(w, req)

		if w.Code != 400 {
			t.Errorf("expected status code 400, got %d", w.Code)
			return
		}
	})
}

func cleanUp(t *testing.T) {
	req := httptest.NewRequest("GET", "/parecer", nil)
	req.Header.Set("Accept", "application/json")

	w := httptest.NewRecorder()

	ReadParecer(w, req)

	var res []parecer.Data
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Errorf("error unmarshaling response body: %s", err)
	}

	for _, parecer := range res {
		req := httptest.NewRequest("DELETE", "/parecer?id="+parecer.ID, nil)
		w := httptest.NewRecorder()
		DeleteParecer(w, req)

		if w.Code != 200 {
			t.Errorf("expected status code 200, got %d", w.Code)
			return
		}
	}
}
