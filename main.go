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

func insertCollector(collector Collector, c appengine.Context) bool {
    encKey := datastore.NewKey(c, "collector",
            collector.GoogleId, 0, nil)
    _, err := datastore.Put(c, encKey, &collector)
    if nil != err {
        log.Println(err)
        return false
    }
    return true
}

func getCollectorFromData(id string, c appengine.Context) (*Collector, bool) {
    encKey := datastore.NewKey(c, "collector", id, 0, nil)
    collector := &Collector{}

    err := datastore.Get(c, encKey, collector)
    if err != nil {
        log.Println(err)
        return collector, false
    }
    return collector, true
}

func strSuccess(success bool) string {
    if success {
        return "success"
    } else {
        return "fail"
    }
}

func responseSuccess(w http.ResponseWriter, success bool) {
    var res Result
    res.Success = strSuccess(success)

    dat, err := json.Marshal(res)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Fprint(w, string(dat))
}

func createCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var collector Collector
    json.Unmarshal(body, &collector)
    success := insertCollector(collector, c)

    var res CollectorWriteResult
    res.Success = strSuccess(success)
    res.GoogleId = collector.GoogleId
    dat, err := json.Marshal(res)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Fprint(w, string(dat))
}

func updateCollector(w http.ResponseWriter, r *http.Request) {
    createCollector(w, r);
}

func getCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    id := r.URL.Query()["googleId"][0]
    collector, succeed := getCollectorFromData(id, c)
    collector.GoogleId = id

    var resp CollectorReadResult
    resp.Success = strSuccess(succeed)
    resp.Collector = *collector
    dat, err := json.Marshal(resp)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Fprint(w, string(dat))
}

func delCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    id := r.URL.Query()["googleId"][0]
    encKey := datastore.NewKey(c, "collector", id, 0, nil)
    err := datastore.Delete(c, encKey)
    responseSuccess(w, err == nil)
}

func collectorHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        createCollector(w, r)
    case "PUT":
        updateCollector(w, r)
    case "GET":
        getCollector(w, r)
    case "DELETE":
        delCollector(w, r)
    }
}

func runeHandler(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    fmt.Fprintf(w, "Rune handler called. method: %s, body: %s",
            r.Method, body)
}
