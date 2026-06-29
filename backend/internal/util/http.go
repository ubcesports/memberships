package util

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func WriteApiResponse(w http.ResponseWriter, status int, code, message string) {
	WriteJson(w, status, map[string]string{"code": code, "message": message})
}
