package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type lineRequestMock struct {
}

func (l lineRequestMock) validateXLineSignature(xLineSignature string, bytes []byte) bool {
	return true
}

func TestHttpPostMsg(t *testing.T) {
	t.Run("it should return httpCode 200", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/msg", strings.NewReader("{}"))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		lr = lineRequestMock{}
		handler := http.HandlerFunc(msgHandler)
		handler.ServeHTTP(resp, req)
		if status := resp.Code; status != http.StatusOK {
			t.Errorf("Wrong code: got %v want %v", status, http.StatusOK)
		}
	})
}

func TestReplyMsg(t *testing.T) {
	t.Run("it should found ReplyMsg", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/msg", strings.NewReader(`{"events":[{"type":"message","replyToken":"8f87f36dce1e4ffca7185b802b7d3538","source":{"userId":"Ucf5a09a475816b57dd1cb2f15214d791","type":"user"},"timestamp":1613371067332,"mode":"active","message":{"type":"text","id":"13560694145870","text":"Lol"}}],"destination":"Uc07389493ea119a5283227c08e2f49f0"}`))
		if err != nil {
			t.Error(err)
		}
		resp := httptest.NewRecorder()
		lr = lineRequestMock{}
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

func TestValidateXLineSignatureWhenNotFoundInHeader(t *testing.T) {
	t.Run("it should return false when Not found XLineSignature in header", func(t *testing.T) {
		_, err := http.NewRequest(http.MethodPost, "/msg", strings.NewReader(`{"timestamp":9}`))
		if err != nil {
			t.Error(err)
		}
		lr = lineRequest{}
		passed := lr.validateXLineSignature("", nil)
		if passed {
			t.Error("Not found XLineSignature in header should be return false.")
		}
	})
}
func TestValidateXLineSignatureWhenFoundInHeader(t *testing.T) {
	t.Run("it should return false when Not found XLineSignature in header", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/msg", strings.NewReader("{}"))
		req.Header.Add("X-Line-Signature", "PKmpEpmCa6DSXf0Bsc3Dtpl55lq2Jrj4o6vk3GZ/kvE=")
		if err != nil {
			t.Error(err)
		}
		lr = lineRequest{}
		channelSecret = "Test"
		passed := lr.validateXLineSignature("PKmpEpmCa6DSXf0Bsc3Dtpl55lq2Jrj4o6vk3GZ/kvE=", nil)
		if passed {
			t.Error("Not found XLineSignature in header should be return false.")
		}
	})
}
