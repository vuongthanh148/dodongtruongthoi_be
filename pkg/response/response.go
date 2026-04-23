package response

import (
"encoding/json"
"net/http"
)

type envelope map[string]any

func Success(w http.ResponseWriter, status int, data any) {
writeJSON(w, status, envelope{
"success": true,
"data":    data,
})
}

func Error(w http.ResponseWriter, status int, message string) {
writeJSON(w, status, envelope{
"success": false,
"message": message,
})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
w.Header().Set("Content-Type", "application/json")
w.WriteHeader(status)
if err := json.NewEncoder(w).Encode(payload); err != nil {
http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
}
