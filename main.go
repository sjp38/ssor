package hello

import (
    "fmt"
    "net/http"
)

func init() {
    http.HandleFunc("/", welcome)
    http.HandleFunc("/Signup", signupHandler)
    http.HandleFunc("/Collector", collectorHandler)
    http.HandleFunc("/Rune", runeHandler)
}

func welcome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Welcome to SOR</h1>")
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "signup handler")
}

func collectorHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "collector handler")
}

func runeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "rune handler")
}
