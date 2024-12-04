package main

import (
	"embed"
	"encoding/json"
	"io/fs"
	"log"
	"mime"
	"net/http"
)

// holds our static web server content.
//
//go:embed all:static
var staticFiles embed.FS

func dirWithStaticFiles() fs.FS {
	sub, _ := fs.Sub(staticFiles, "static")
	return sub
}

func databases(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // for CORS
	w.WriteHeader(http.StatusOK)
	test := []string{}
	test = append(test, "Hello")
	test = append(test, "World")
	err := json.NewEncoder(w).Encode(test)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Windows may be missing this
	err := mime.AddExtensionType(".js", "application/javascript")
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/test", http.HandlerFunc(databases))
	http.Handle("/", http.FileServer(http.FS(dirWithStaticFiles())))
	log.Println("Server listening on http://127.0.0.1:8080")
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
