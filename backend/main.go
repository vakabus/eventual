package main

import (
	"embed"
	"events/backend/config"
	"events/backend/routes"
	"log"
	"mime"
	"net/http"
)

//go:generate go run github.com/gzuidhof/tygo@latest generate
//go:generate sh -c "cd ../frontend && npm run build"

// holds our static web server content.
//
//go:embed all:static
var staticFiles embed.FS

func main() {
	// Initialize configuration
	config.InitializeViper()

	// Windows may be missing this
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()
	routes.AddRoutes(router, staticFiles)

	log.Println("Server listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", router))
}
