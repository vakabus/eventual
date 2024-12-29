package api

import (
	"encoding/json"
	"events/backend/routes/api/types"
	"log"
	"net/http"
)

func errorJson(w http.ResponseWriter, err string, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	e := json.NewEncoder(w).Encode(types.ErrorResponse{ErrorMessage: err})
	if e != nil {
		log.Println("ERROR: encoding error response:", e)
	}
}
