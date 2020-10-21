package main

import (
	"encoding/json"
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

func TestReplyMsg(t *testing.T) {
	t.Run("it should found ReplyMsg", func(t *testing.T) {
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

		var replyMsg ReplyMsg
		err = json.Unmarshal(resp.Body.Bytes(), &replyMsg)
		if err != nil {
			t.Errorf("Can't unmarshal to type ReplyMsg found error %v", err)
		}

	})
}
