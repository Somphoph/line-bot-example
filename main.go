package main

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type lineRequest struct {
}
type lineRequestUtil interface {
	validateXLineSignature(r *http.Request) bool
}

var lr lineRequestUtil

func init() {
	lr = lineRequest{}
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
		return
	}
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}

	if !lr.validateXLineSignature(r) {
		http.NotFound(w, r)
		return
	}
	replyMsg := ReplyMsg{
		ReplyToken: "",
		Type:       "message",
		Mode:       "active",
		Timestamp:  getMillisecondTime(),
		Source: Source{
			Type:   "user",
			UserId: "U4af4980629...",
		},
		Message: Message{
			Id:   "325708",
			Type: "text",
			Text: "Hello, world! (love)",
		},
	}

	js, err := json.Marshal(replyMsg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
func (lr lineRequest) validateXLineSignature(r *http.Request) bool {
	decoded, err := base64.StdEncoding.DecodeString(r.Header.Get("X-Line-Signature"))
	if err != nil {
		return false
	}
	return hmac.Equal(decoded, []byte("6952513badb650d5cf9c14a3c79cd8c8"))
}
func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/msg", msgHandler)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func getMillisecondTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

type ReplyMsg struct {
	ReplyToken string  `json:"replayToken"`
	Type       string  `json:"type"`
	Mode       string  `json:"mode"`
	Timestamp  int64   `json:"timestamp"`
	Source     Source  `json:"source"`
	Message    Message `json:"message"`
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

type WebHookEvent struct {
	Destination string
	Events      []Event
	ReplyToken  string `json:"replyToken"`
	Type        string `json:"type"`
	Mode        string `json:"mode"`
	Timestamp   int64  `json:"timestamp"`
	Source      Source `json:"source"`
}
type Event struct {
	ReplyToken string  `json:"replyToken"`
	Type       string  `json:"type"`
	Mode       string  `json:"mode"`
	Timestamp  int64   `json:"timestamp"`
	Source     Source  `json:"source"`
	Message    Message `json:"message"`
}
