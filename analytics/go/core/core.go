package core

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Data struct {
	APIKey       string `json:"api_key"`
	Hostname     string `json:"hostname"`
	Path         string `json:"path"`
	UserAgent    string `json:"user_agent"`
	Method       int    `json:"method"`
	ResponseTime int64  `json:"response_time"`
	Status       int    `json:"status"`
	Framework    int8   `json:"framework"`
}

var MethodMap = map[string]int{
	"GET":     0,
	"POST":    1,
	"PUT":     2,
	"PATCH":   3,
	"DELETE":  4,
	"OPTIONS": 5,
	"CONNECT": 6,
	"HEAD":    7,
	"TRACE":   8,
}

func LogRequest(data Data) {
	reqBody, err := json.Marshal(data)
	if err != nil {
		print(err)
	}
	http.Post("https://api-analytics-server.vercel.app/api/log-request", "application/json", bytes.NewBuffer(reqBody))
}