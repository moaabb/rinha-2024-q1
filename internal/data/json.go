package data

import (
	"encoding/json"
	"net/http"
)

type H map[string]string

func Response(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(data)
}
