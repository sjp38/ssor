package hello

import (
    "fmt"
    "net/http"
)

func init() {
    http.HandleFunc("/", welcome)
}

func welcome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Welcome to SOR</h1>")
}
