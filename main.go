package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type lineRequest struct {
}
type lineRequestUtil interface {
	validateXLineSignature(xLineSignature string, bytes []byte) bool
}

var lr lineRequestUtil
var channelSecret string

func init() {
	lr = lineRequest{}
	channelSecret = "<Your Channel Secret>"
}
func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

func msgHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/msg" {
		http.NotFound(w, r)
		log.Printf("URL path not found.")
		return
	}
	if r.Method != "POST" {
		http.NotFound(w, r)
		log.Printf("Http method not support.")
		return
	}
	bodyBytes, err := ioutil.ReadAll(r.Body)

	if !lr.validateXLineSignature(r.Header.Get("X-Line-Signature"), bodyBytes) {
		http.Error(w, "XLineSignature validate fail.", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Can't read request body.", http.StatusBadRequest)
	}
	log.Printf("Request :" + string(bodyBytes))
	var webHookEvent WebHookEvent
	err = json.Unmarshal(bodyBytes, &webHookEvent)
	if err != nil {
		http.Error(w, "Can't read request body.", http.StatusBadRequest)
	}
	replyMsg := ReplyMsg{
		ReplyToken: webHookEvent.Events[0].ReplyToken,
		Type:       "message",
		Mode:       "active",
		Timestamp:  getMillisecondTime(),
		Source:     webHookEvent.Events[0].Source,
		Message: Message{
			Id:   "325708",
			Type: "text",
			Text: "Hello, world! (love)",
		},
	}
	webHookStatus := WebHookStatus{
		Success:    true,
		Timestamp:  time.Now().Format(time.RFC3339),
		StatusCode: 200,
		Reason:     "OK",
		Detail:     "200",
	}
	var js []byte
	if len(webHookEvent.Events) > 0 {
		log.Printf("Reply message response.")
		js, err = json.Marshal(replyMsg)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		log.Printf("WebHookEvent response.")
		js, err = json.Marshal(webHookStatus)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf(string(js))
	w.Write(js)
}
func (lr lineRequest) validateXLineSignature(xLineSignature string, body []byte) bool {
	decoded, err := base64.StdEncoding.DecodeString(xLineSignature)
	if err != nil {
		return false
	}
	hash := hmac.New(sha256.New, []byte(channelSecret))

	_, err = hash.Write(body)
	if err != nil {
		return false
	}

	return hmac.Equal(decoded, hash.Sum(nil))
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/msg", msgHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf(`{"message":"Listening on port %s"}`, port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getMillisecondTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type ReplyMsg struct {
	ReplyToken string  `json:"replyToken"`
	Type       string  `json:"type"`
	Mode       string  `json:"mode"`
	Timestamp  int64   `json:"timestamp"`
	Source     Source  `json:"source"`
	Message    Message `json:"message"`
}
type WebHookStatus struct {
	Success    bool   `json:"success"`
	Timestamp  string `json:"timestamp"`
	StatusCode int    `json:"statusCode"`
	Reason     string `json:"reason"`
	Detail     string `json:"detail"`
}

type Source struct {
	Type   string `json:"type"`
	UserId string `json:"userId"`
}

type Message struct {
	Id   string `json:"id"`
	Type string `json:"type"`
	Text string `json:"text"`
}
type Emoji struct {
	Index     int    `json:"index"`
	Length    int    `json:"length"`
	ProductId string `json:"productId"`
	EmojiId   string `json:emojiId`
}

type WebHookEvent struct {
	Destination string  `json:"destination"`
	Events      []Event `json:"events"`
}
type Event struct {
	ReplyToken string  `json:"replyToken"`
	Type       string  `json:"type"`
	Mode       string  `json:"mode"`
	Timestamp  int64   `json:"timestamp"`
	Source     Source  `json:"source"`
	Message    Message `json:"message"`
}
