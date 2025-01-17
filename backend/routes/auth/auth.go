package auth

import "net/http"

func Server() http.Handler {
	initializeGoogleOAuth()

	mux := http.NewServeMux()
	mux.HandleFunc("/google", handleGoogleLogin)
	mux.HandleFunc("/google-callback", handleGoogleCallback)

	mux.HandleFunc("/logout", handleLogout)
	return mux
}
