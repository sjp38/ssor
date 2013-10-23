package sor

import (
    "appengine"
    "appengine/datastore"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
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
            "", int64(collector.GoogleId), nil)
    _, err := datastore.Put(c, encKey, &collector)
    if nil != err {
        log.Println(err)
        return false
    }
    return true
}

func getCollectorFromData(id int, c appengine.Context) (*Collector, bool) {
    encKey := datastore.NewKey(c, "collector", "", int64(id), nil)
    collector := &Collector{}

    err := datastore.Get(c, encKey, collector)
    if err != nil {
        log.Println(err)
        return collector, false
    }
    return collector, true
}

func createCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var collector Collector
    json.Unmarshal(body, &collector)
    result := insertCollector(collector, c)
    var resp Result
    if result {
        resp.Success = "success"
    } else {
        resp.Success = "fail"
    }
    dat, err := json.Marshal(resp)
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
    id, _ := strconv.Atoi(r.URL.Query()["googleId"][0])
    collector, succeed := getCollectorFromData(id, c)
    if false == succeed {
        var resp Result
        resp.Success = "fail"
        dat, err := json.Marshal(resp)
        if err != nil {
            log.Println(err)
            return
        }
        fmt.Fprint(w, string(dat))
    } else {
        dat, err := json.Marshal(collector)
        if err != nil {
            log.Println(err)
            return
        }
        fmt.Fprint(w, string(dat))
    }
}

func delCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    id, err := strconv.Atoi(r.URL.Query()["googleId"][0])
    log.Println("WRYYY")
    log.Println("will delete collector with id %s", id)
    if err != nil {
        log.Println(err)
        return
    }
    encKey := datastore.NewKey(c, "collector", "", int64(id), nil)
    err = datastore.Delete(c, encKey)
    var resp Result
    if err != nil {
        resp.Success = "fail"
        dat, err := json.Marshal(resp)
        if err != nil {
            log.Println(err)
            return
        }
        fmt.Fprint(w, string(dat))
    } else {
        resp.Success = "success"
        dat, err := json.Marshal(resp)
        if err != nil {
            log.Println(err)
            return
        }
        fmt.Fprint(w, string(dat))
    }
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
