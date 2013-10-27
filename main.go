package sor

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "net/http"
    "time"
)

type result struct {
    success bool
    reason string
}

func init() {
    http.HandleFunc("/", welcome)
    http.HandleFunc("/collector", collectorHandler)
    http.HandleFunc("/rune", runeHandler)
    http.HandleFunc("/runes", runesHandler)
    http.HandleFunc("/fight", fightHandler)
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

func doSetCollectorInitStat(collector *Collector, maxHp int, maxMp int,
        atk int, def int, int int) {
    collector.MaxHp = maxHp
    collector.Hp = maxHp
    collector.MaxMp = maxMp
    collector.Mp = maxMp
    collector.Atk = atk
    collector.Def = def
    collector.Int = int
}

func setCollectorInitStat(collector *Collector) {
    switch collector.CollectorClass {
    case "Geek":
        doSetCollectorInitStat(collector, 100,100,10,5,10)
    case "Nerd":
        doSetCollectorInitStat(collector, 100,100,10,10,5)
    case "Dork":
        doSetCollectorInitStat(collector, 100,100,5,10,10)
    }
}

func createCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var collectorMinInfo CollectorMinInfo
    json.Unmarshal(body, &collectorMinInfo)

    var collector Collector
    collector.GoogleId = collectorMinInfo.GoogleId
    collector.Email = collectorMinInfo.Email
    collector.ProfileUrl = collectorMinInfo.ProfileUrl
    collector.Nickname = collectorMinInfo.Nickname
    collector.CollectorClass = collectorMinInfo.CollectorClass
    setCollectorInitStat(&collector)
    collector.Exp = 0
    collector.ScanCount = 5

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

func setRuneOwner(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var runeMinInfo RuneMinInfo
    json.Unmarshal(body, &runeMinInfo)

    rune, succeed := getRuneFromData(runeMinInfo.ISBN, c)
    if succeed == false {
        respFail(w, "fail to get rune from datastore")
        return
    }
    rune.OwnerGoogleId = runeMinInfo.OwnerGoogleId
    succeed = insertRune(*rune, c)
    if succeed == false {
        respFail(w, "fail to updated info to datastore")
        return
    }

    var runeResult RuneResult
    runeResult.Success = "success"
    runeResult.Rune = *rune
    respInJson(w, runeResult)
}

func registerRune(w http.ResponseWriter, r *http.Request, rune Rune) bool {
    c := appengine.NewContext(r)

    success := insertRune(rune, c)
    if success == false {
        return false
    }
    return true
}

func makeRune(c appengine.Context, isbn string) (Rune, bool) {
    var rune Rune
    searchUrl := "http://apis.daum.net/search/book"
    searchUrl += "?output=json&apikey=" + daumApiKey
    searchUrl += "&q=" + isbn

    client := urlfetch.Client(c)
    resp, err := client.Get(searchUrl)
    if err != nil {
        log.Print(err)
        return rune, false
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    var searchResult DaumBookSearchResult
    json.Unmarshal(body, &searchResult)
    if len(searchResult.Channel.Item) <= 0 {
        return rune, false
    }
    itemInfo := searchResult.Channel.Item[0]

    rune.ISBN = isbn

    imageUrl := itemInfo.Cover_l_url
    imageUrl = imageUrl[0:len("http://book.daum-img.net/")] +
            "image" + imageUrl[len("http://book.daum-img.net/R110x160"):]
    rune.ImageUrl = imageUrl

    rune.Title = itemInfo.Title
    rune.Type = "Basic"
    rune.MaxHp = 10
    rune.Hp = 10
    rune.OwnerGoogleId = ""
    return rune, true
}

// Get rune info
// If not exist in datastore yet, register new rune to datastore
func getRune(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    isbn := r.URL.Query()["ISBN"][0]

    rune, succeed := getRuneFromData(isbn, c)
    if succeed == false {
        log.Println("fail to get rune from datastore. make it!")
        newRune, succeed := makeRune(c, isbn)
        if succeed == false {
            respFail(w, "fail to make rune")
            return
        }
        log.Println("made rune. register it!!")
        registerRune(w, r, newRune)
        *rune = newRune
    }
    var runeResult RuneResult
    runeResult.Success = "success"
    runeResult.Rune = *rune
    respInJson(w, runeResult)
}

func getRunes(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    googleId := r.URL.Query()["googleId"][0]

    q := datastore.NewQuery("rune").Filter("OwnerGoogleId = ", googleId)
    var runes []Rune
    _, err := q.GetAll(c, &runes)
    if err != nil {
        log.Println(err)
        respFail(w, "fail to get rune from datastore")
    }

    respInJson(w, runes)
}

func do_fight(attacker *Collector, defender *Collector, rune *Rune) {
    rand.Seed(time.Now().UTC().UnixNano())
    attackPoint := attacker.Atk
    attackPoint += rand.Intn(int(float32(attackPoint) * 0.1) + 1)
    defencePoint := defender.Def
    defencePoint += rand.Intn(int(float32(defencePoint) * 0.1) + 1)

    damage := attackPoint - defencePoint
    if damage > 0 {
        rune.Hp -= damage
    } else {
        attacker.Hp += damage
    }
}

func fight(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    attackerId := r.URL.Query()["attacker"][0]
    defenderId := r.URL.Query()["defender"][0]
    isbn := r.URL.Query()["ISBN"][0]

    attacker, succeed := getCollectorFromData(attackerId, c)
    if succeed == false {
        respFail(w, "fail to get attacker")
        return
    }
    defender, succeed := getCollectorFromData(defenderId, c)
    if succeed == false {
        respFail(w, "fail to get defender")
        return
    }
    rune, succeed := getRuneFromData(isbn, c)
    if succeed == false {
        respFail(w, "fail to get rune")
        return
    }
    do_fight(attacker, defender, rune)

    succeed = insertCollector(*attacker, c)
    if succeed == false {
        respFail(w, "fail to update attacker")
        return
    }
    succeed = insertRune(*rune, c)
    if succeed == false {
        respFail(w, "fail to update rune")
        return
    }

    var fightResult FightResult
    fightResult.Success = "success"
    fightResult.Attacker = *attacker
    fightResult.Defender = *defender
    fightResult.Rune = *rune
    respInJson(w, fightResult)
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
        setRuneOwner(w, r)
    case "GET":
        getRune(w, r)
    case "DELETE":
        fmt.Fprintf(w, "DELETE is not supported")
    }
}

func runesHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        fmt.Fprint(w, "not supported")
    case "PUT":
        fmt.Fprint(w, "not supported")
    case "GET":
        getRunes(w, r)
    case "DELETE":
        fmt.Fprint(w, "not supported")
    default:
    }
}

func fightHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        fmt.Fprint(w, "not supported")
    case "PUT":
        fmt.Fprint(w, "not supported")
    case "GET":
        fight(w, r)
    case "DELETE":
        fmt.Fprint(w, "not supported")
    default:
        fmt.Fprint(w, "not supported")
    }
}
