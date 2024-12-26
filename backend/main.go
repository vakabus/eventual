package main

import (
	"embed"
	"events/backend/api"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
)

//go:generate go run github.com/gzuidhof/tygo@latest generate
//go:generate sh -c "cd ../frontend && npm run build"

// holds our static web server content.
//
//go:embed all:static
var staticFiles embed.FS

type htmlDir struct {
	d http.Dir
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
		return http.FileServer(htmlDir{http.FS(sub).(http.Dir)})
	} else {
		// else serve from the static directory
		return http.FileServer(htmlDir{http.Dir("static")})
	}
}

func main() {
	// Windows may be missing this
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", api.Server()))
	router.Handle("/static/", http.StripPrefix("/static", staticHandler()))
	router.Handle("/{$}", http.RedirectHandler("/static/", http.StatusTemporaryRedirect))

	log.Println("Server listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
