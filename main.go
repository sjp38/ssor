package sor

import (
    "appengine"
    "appengine/datastore"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
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

func createCollector(collector Collector, c appengine.Context) bool {
    encKey := datastore.NewKey(c, "collector",
            "", int64(collector.GoogleId), nil)
    _, err := datastore.Put(c, encKey, &collector)
    if nil != err {
        log.Println(err)
        return false
    }
    return true
}

func collectorHandler(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var collector Collector
    switch r.Method {
    case "POST":
        json.Unmarshal(body, &collector)
        createCollector(collector, c)
    case "PUT":
    case "GET":
    case "DEL":
        fmt.Fprintf(w, "Implementing yet...")
    }
}

func runeHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    fmt.Fprintf(w, "Rune handler called. method: %s, body: %s",
            r.Method, body)
}
