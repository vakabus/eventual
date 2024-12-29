package main

import (
	"embed"
	"events/backend/config"
	"events/backend/routes"
	"events/backend/routes/auth"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

//go:generate go run github.com/gzuidhof/tygo@latest generate
//go:generate sh -c "cd ../frontend && npm run build"

// holds our static web server content.
//
//go:embed all:static
var staticFiles embed.FS

type htmlDir struct {
	d http.FileSystem
}

func (d htmlDir) Open(name string) (http.File, error) {
	// Try name as supplied
	f, err := d.d.Open(name)
	if os.IsNotExist(err) {
		// Not found, try with .html
		if f, err := d.d.Open(name + ".html"); err == nil {
			return f, nil
		}
	}
	return f, err
}

func staticHandler() http.Handler {
	if os.Getenv("DEV") == "" {
		// if not set or empty, serve static content from embeded fs
		sub, _ := fs.Sub(staticFiles, "static")
		return http.StripPrefix("/static", http.FileServer(htmlDir{http.FS(sub)}))
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

func main() {
	// Initialize configuration
	config.InitializeViper()

	// Windows may be missing this
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	routes.AddRoutes(router)
	router.Handle("/static/", staticHandler())
	router.HandleFunc("/{$}", func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.GetUserFromCookies(r)
		if err != nil {
			http.Redirect(w, r, "/static/", http.StatusTemporaryRedirect)
		} else {
			http.Redirect(w, r, "/static/dashboard", http.StatusTemporaryRedirect)
		}
	})

	log.Println("Server listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
