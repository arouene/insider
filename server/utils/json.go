package utils

import (
	"encoding/json"
	"net/http"
)

type Struct map[string]interface{}

func JSON(w http.ResponseWriter, s Struct, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(s)
}
