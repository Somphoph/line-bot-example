package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpPostMsg(t *testing.T) {
	t.Run("it should return httpCode 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/msg", nil)
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		handler := http.HandlerFunc(msgHandler)
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("Wrong code: got %v want %v", status, http.StatusOK)
		}
	})
}
