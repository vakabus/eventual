package routes

import (
	"embed"
	"events/backend/routes/api"
	"events/backend/routes/auth"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func staticHandler(staticFiles embed.FS) http.Handler {
	if os.Getenv("DEV") == "" {
		// if not set or empty, serve static content from embeded fs
		sub, _ := fs.Sub(staticFiles, "static")
		return http.StripPrefix("/app", http.FileServer(http.FS(sub)))
	} else {
		// else serve by reverse proxying to vite dev server
		url, err := url.Parse("http://localhost:5173")
		if err != nil {
			panic(err)
		}
		proxy := httputil.NewSingleHostReverseProxy(url)
		return proxy
	}
}

func AddRoutes(mux *http.ServeMux, staticFiles embed.FS) {
	mux.Handle("/api/", http.StripPrefix("/api", api.Server()))
	mux.Handle("/auth/", http.StripPrefix("/auth", auth.Server()))
	mux.Handle("/app/", staticHandler(staticFiles))
	mux.Handle("/{$}", http.RedirectHandler("/app/", http.StatusTemporaryRedirect))
}
