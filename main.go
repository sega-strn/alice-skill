package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
)

// Request структура для входящего запроса от Алисы
type Request struct {
	Meta  Meta    `json:"meta"`
	State State   `json:"state"`
	Request AliceRequest `json:"request"`
	Session Session `json:"session"`
	Version string  `json:"version"`
}

// Meta содержит информацию о навыке
type Meta struct {
	Locale     string `json:"locale"`
	Timezone   string `json:"timezone"`
	ClientID   string `json:"client_id"`
	Interfaces struct {
		Screen interface{} `json:"screen"`
	} `json:"interfaces"`
}

// State содержит состояние
type State struct {
	Session map[string]interface{} `json:"session"`
	User    map[string]interface{} `json:"user"`
}

// AliceRequest содержит данные запроса
type AliceRequest struct {
	Command           string `json:"command"`
	OriginalUtterance string `json:"original_utterance"`
	Type              string `json:"type"`
}

// Session содержит информацию о сессии
type Session struct {
	New       bool   `json:"new"`
	MessageID int    `json:"message_id"`
	SessionID string `json:"session_id"`
	SkillID   string `json:"skill_id"`
	UserID    string `json:"user_id"`
}

// Response структура для ответа Алисе
type Response struct {
	Response struct {
		Text string `json:"text"`
		TTS  string `json:"tts,omitempty"`
		End  bool   `json:"end_session"`
	} `json:"response"`
	Session Session `json:"session"`
	Version string  `json:"version"`
}

var flagRunAddr string

func parseFlags() {
	flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
	flag.Parse()
}

func handleAlice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Создаем ответ
	resp := Response{
		Version: req.Version,
		Session: req.Session,
	}

	// Простая логика ответов
	switch {
	case req.Session.New:
		resp.Response.Text = "Привет! Я новый навык. Чем могу помочь?"
	case req.Request.Command == "помощь" || req.Request.Command == "что ты умеешь":
		resp.Response.Text = "Я могу отвечать на простые вопросы. Например, спросите меня о погоде."
	default:
		resp.Response.Text = "Извините, я пока не знаю, как ответить на это."
	}

	resp.Response.End = false

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func main() {
	parseFlags()
	http.HandleFunc("/", handleAlice)
	log.Printf("Starting server on %s", flagRunAddr)
	if err := http.ListenAndServe(flagRunAddr, nil); err != nil {
		log.Fatal(err)
	}
}
