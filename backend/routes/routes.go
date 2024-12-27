package routes

import (
	"events/backend/routes/api"
	"events/backend/routes/auth"
	"net/http"
)

func AddRoutes(mux *http.ServeMux) {
	mux.Handle("/api/", http.StripPrefix("/api", api.Server()))
	mux.Handle("/auth/", http.StripPrefix("/auth", auth.Server()))
}
