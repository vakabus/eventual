package main

import (
    "encoding/json"
    "log"
    "mime"
    "net/http"
)

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
    http.Handle("/", http.FileServer(http.Dir("../frontend/build")))
    log.Println("Server listening on http://127.0.0.1:8080")
    log.Fatal(http.ListenAndServe("127.0.0.1:8080", nil))
}
