package sor

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func init() {
    http.HandleFunc("/", welcome)
    http.HandleFunc("/collector", collectorHandler)
    http.HandleFunc("/rune", runeHandler)
}

func welcome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<h1>Welcome to SOR</h1>")
}

func collectorHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    if "POST" == r.Method {
        var collector Collector
        json.Unmarshal(body, &collector)
        fmt.Fprintf(w, "parsed nick: %s\n", collector.Nickname)
    }
    fmt.Fprintf(w, "Collector handler called. method: %s, body: %s\n",
            r.Method, body)
}

func runeHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    fmt.Fprintf(w, "Rune handler called. method: %s, body: %s",
            r.Method, body)
}
