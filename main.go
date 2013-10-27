package sor

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
)

type result struct {
    success bool
    reason string
}

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

// Response in json form
func respInJson(w http.ResponseWriter, data interface{}) {
    dat, err := json.Marshal(data)
    if err != nil {
        log.Println(err)
        return
    }
    fmt.Fprint(w, string(dat))
}

func respFail(w http.ResponseWriter, reason string) {
    var res FailResult
    res.Success = "fail"
    res.Reason = reason
    respInJson(w, res)
}

func respCollector(w http.ResponseWriter, r result, collector *Collector) {
    if r.success == false {
        respFail(w, r.reason)
        return
    }
    var res CollectorResult
    res.Success = "success"
    res.Collector = *collector
    respInJson(w, res)
}

func createCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var collector Collector
    json.Unmarshal(body, &collector)
    success := insertCollector(collector, c)

    respCollector(w, result{success, "Unknown"}, &collector)
}

func updateCollector(w http.ResponseWriter, r *http.Request) {
    createCollector(w, r);
}

func getCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    id := r.URL.Query()["googleId"][0]
    collector, succeed := getCollectorFromData(id, c)

    respCollector(w, result{succeed, "Unknown"}, collector)
}

func delCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    id := r.URL.Query()["googleId"][0]
    encKey := datastore.NewKey(c, "collector", id, 0, nil)
    collector, exist := getCollectorFromData(id, c)
    var res result
    if exist {
        err := datastore.Delete(c, encKey)
        res = result{err == nil, "datastore error"}
    } else {
        res = result{false, "not exist"}
    }

    respCollector(w, res, collector)
}

func insertRune(rune Rune, c appengine.Context) bool {
    encKey := datastore.NewKey(c, "rune",
            rune.ISBN, 0, nil)
    _, err := datastore.Put(c, encKey, &rune)
    if nil != err {
        log.Println(err)
        return false
    }
    return true
}

func getRuneFromData(isbn string, c appengine.Context) (*Rune, bool) {
    encKey := datastore.NewKey(c, "rune", isbn, 0, nil)
    rune := &Rune{}

    err := datastore.Get(c, encKey, rune)
    if err != nil {
        log.Println(err)
        return rune, false
    }
    return rune, true
}

func registerRune(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    isbn := r.URL.Query()["ISBN"][0]

    searchUrl := "http://apis.daum.net/search/book"
    searchUrl += "?output=json&apikey=" + daumApiKey
    searchUrl += "&q=" + isbn

    client := urlfetch.Client(c)
    resp, err := client.Get(searchUrl)
    if err != nil {
        log.Print(err)
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    var searchResult DaumBookSearchResult
    json.Unmarshal(body, &searchResult)
    itemInfo := searchResult.Channel.Item[0]

    var rune Rune
    rune.ISBN = isbn
    rune.ImageUrl = itemInfo.Cover_l_url
    rune.Title = itemInfo.Title
    rune.Type = "Basic"
    rune.MaxHp = 10
    rune.Hp = 5

    success := insertRune(rune, c)
    if success == false {
        respFail(w, "datastore insert fail")
        return
    }
    var runeResult RuneResult
    runeResult.Success = "success"
    runeResult.Rune = rune
    respInJson(w, runeResult)
}

func getRune(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    isbn := r.URL.Query()["ISBN"][0]

    rune, succeed := getRuneFromData(isbn, c)
    if succeed == false {
        respFail(w, "datastore failure")
        return
    }
    var runeResult RuneResult
    runeResult.Success = "success"
    runeResult.Rune = *rune
    respInJson(w, runeResult)
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
    switch r.Method {
    case "POST":
        fmt.Fprint(w, "POST is not supported")
    case "PUT":
        registerRune(w, r)
    case "GET":
        getRune(w, r)
    case "DELETE":
        fmt.Fprintf(w, "DELETE is not supported")
    }
}
