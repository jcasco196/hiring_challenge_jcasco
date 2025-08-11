package api

import (
    "encoding/json"
    "net/http"
)

func OKResponse(w http.ResponseWriter, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    if data != nil {
        json.NewEncoder(w).Encode(data)
    }
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    resp := map[string]string{"error": message}
    json.NewEncoder(w).Encode(resp)
}