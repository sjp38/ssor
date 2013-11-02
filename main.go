package sor

import (
    "appengine"
    "appengine/datastore"
    "appengine/urlfetch"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "math"
    "math/rand"
    "net/http"
    "time"
)

type result struct {
    success bool
    reason string
}

const (
    HEAL_COLLECTOR_UNIT = 10
    HEAL_RUNE_UNIT = 1
    HEAL_MP_UNIT = 5
    PLAYSTORE_URL = "https://play.google.com/store/apps/details?id="
    APP_DOWN_URL = PLAYSTORE_URL + APP_PACKAGE_NAME
)

func init() {
    http.HandleFunc("/", welcome)
    http.HandleFunc("/collector", collectorHandler)
    http.HandleFunc("/rune", runeHandler)
    http.HandleFunc("/runes", runesHandler)
    http.HandleFunc("/fight", fightHandler)
    http.HandleFunc("/heal", healHandler)
    http.HandleFunc("/changestat", changeStatHandler)
}

func welcome(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "<center><h1>Welcome to SOR</h1><br>\n")
    fmt.Fprintf(w, "<h1><a href=\"%s\">Download from Google play" +
            "</a></h1>\n</center>", APP_DOWN_URL)
}

func insertCollector(collector CollectorInternal, c appengine.Context) bool {
    encKey := datastore.NewKey(c, "collector",
            collector.GoogleId, 0, nil)
    _, err := datastore.Put(c, encKey, &collector)
    if nil != err {
        log.Println(err)
        return false
    }
    return true
}

func getCollectorFromData(id string, c appengine.Context) (
        *CollectorInternal, bool) {
    encKey := datastore.NewKey(c, "collector", id, 0, nil)
    collector := &CollectorInternal{}

    err := datastore.Get(c, encKey, collector)
    if err != nil {
        log.Println(err)
        return collector, false
    }

    now := time.Now().UTC().Unix()
    if collector.LastMpConsumedTime != 0 {
        mpHeal := (now - collector.LastMpConsumedTime) / 60
        collector.Mp += int(mpHeal)
        if collector.Mp > collector.MaxMp {
            collector.Mp = collector.MaxMp
        }
        collector.LastMpConsumedTime = now

        success := insertCollector(*collector, c)
        if success == false {
            return collector, false
        }
    }
    if collector.LastScannedTime != 0 {
        scanHeal := (now - collector.LastScannedTime) / 120
        collector.ScanCount += int(scanHeal)
        if collector.ScanCount > 5 {
            collector.ScanCount = 5
        }
        collector.LastScannedTime = now

        success := insertCollector(*collector, c)
        if success == false {
            return collector, false
        }
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
    log.Println("failed: ", reason)
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

    var collectorInternal CollectorInternal
    collectorInternal.CreatedTime = time.Now().UTC().Unix()

    collector := &collectorInternal.Collector
    collector.GoogleId = collectorMinInfo.GoogleId
    collector.Email = collectorMinInfo.Email
    collector.ProfileUrl = collectorMinInfo.ProfileUrl
    collector.Nickname = collectorMinInfo.Nickname
    collector.CollectorClass = collectorMinInfo.CollectorClass
    setCollectorInitStat(collector)
    collector.Level = 1
    collector.ExpToNext = collector.Level * 100 +
            int(math.Pow(8, float64(collector.Level - 1)))
    collector.Exp = 0
    collector.ScanCount = 5
    collector.BonusPoint = 0

    success := insertCollector(collectorInternal, c)

    respCollector(w, result{success, "Unknown"}, collector)
}

func updateCollector(w http.ResponseWriter, r *http.Request) {
    createCollector(w, r);
}

func getCollector(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    id := r.URL.Query()["googleId"][0]
    collector, succeed := getCollectorFromData(id, c)

    respCollector(w, result{succeed, "Unknown"}, &collector.Collector)
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

    respCollector(w, res, &collector.Collector)
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

    rune.ThumbnailUrl = itemInfo.Cover_l_url

    rune.Title = itemInfo.Title
    rune.Type = "Basic"
    rune.MaxHp = 10
    rune.Hp = 10
    rune.OwnerGoogleId = ""
    return rune, true
}

func increaseExp(collector *Collector, exp int) {
    collector.Exp += exp
    if collector.Exp >= collector.ExpToNext {
        collector.Level++
        collector.Exp -= collector.ExpToNext
        collector.ExpToNext = collector.Level * 100 +
                int(math.Pow(8, float64(collector.Level - 1)))
        collector.BonusPoint += 8
    }
}

// Get rune info
// If not exist in datastore yet, register new rune to datastore
func getRune(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    isbn := r.URL.Query()["ISBN"][0]

    googleIdExist := len(r.URL.Query()["googleId"]) > 0
    if googleIdExist {
        googleId := r.URL.Query()["googleId"][0]
        collector, succeed := getCollectorFromData(googleId, c)
        if collector.ScanCount <= 0 {
            respFail(w, "scan count is not enough")
            return
        }
        collector.LastScannedTime = time.Now().UTC().Unix()
        collector.TotalScanCount++
        collector.ScanCount--
        increaseExp(&collector.Collector, 1)
        succeed = insertCollector(*collector, c)
        if succeed == false {
            respFail(w, "fail to update collector")
            return
        }
    }

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
        return
    }

    var runesResult RunesResult
    runesResult.Success = "success"
    runesResult.Runes = runes
    respInJson(w, runesResult)
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
        if rune.Hp < 0 {
            rune.Hp = 0
        }
        increaseExp(attacker, damage)
    } else {
        attacker.Hp += damage
        increaseExp(attacker, 1)
        rune.Hp -= 1
        if attacker.Hp < 0 {
            attacker.Hp = 0
        }
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
    do_fight(&attacker.Collector, &defender.Collector, rune)

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
    fightResult.Attacker = attacker.Collector
    fightResult.Defender = defender.Collector
    fightResult.Rune = *rune
    respInJson(w, fightResult)
}

func healCollector(healRequest HealRequest,
        w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    collector, succeed := getCollectorFromData(healRequest.GoogleId, c)
    if succeed == false {
        respFail(w, "fail to get collector from datastore")
        return
    }

    if collector.Mp < HEAL_MP_UNIT {
        respFail(w, "mp is not enough")
        return
    }
    collector.Mp -= HEAL_MP_UNIT
    collector.LastMpConsumedTime = time.Now().UTC().Unix()
    collector.Hp += HEAL_COLLECTOR_UNIT
    increaseExp(&collector.Collector, HEAL_COLLECTOR_UNIT)
    if collector.Hp > collector.MaxHp {
        collector.Hp = collector.MaxHp
    }
    succeed = insertCollector(*collector, c)
    if succeed == false {
        respFail(w, "fail to update healed collector")
        return
    }
    var collectorResult CollectorResult
    collectorResult.Success = "success"
    collectorResult.Collector = collector.Collector
    respInJson(w, collectorResult)
}

func healRune(request HealRequest,
        w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    rune, succeed := getRuneFromData(request.ISBN, c)
    if succeed == false {
        respFail(w, "fail to get rune from datastore")
        return
    }
    collector, succeed := getCollectorFromData(rune.OwnerGoogleId, c)
    if succeed == false {
        respFail(w, "fail to get collector from datastore")
        return
    }

    if request.Type == "mp" {
        if collector.Mp < HEAL_MP_UNIT {
            respFail(w, "mp is not enough")
            return
        }
        collector.Mp -= HEAL_MP_UNIT
        collector.LastMpConsumedTime = time.Now().UTC().Unix()
    } else {
        if collector.ScanCount < 1 {
            respFail(w, "scan count is not enough")
            return
        }
        collector.ScanCount--
    }
    rune.Hp += HEAL_RUNE_UNIT
    increaseExp(&collector.Collector, HEAL_RUNE_UNIT)
    if rune.Hp > rune.MaxHp {
        rune.Hp = rune.MaxHp
    }
    succeed = insertCollector(*collector, c)
    if succeed == false {
        respFail(w, "fail to update consumed collector")
        return
    }
    succeed = insertRune(*rune, c)
    if succeed == false {
        respFail(w, "fail to update healed rune")
        return
    }
    var runeResult RuneResult
    runeResult.Success = "success"
    runeResult.Rune = *rune
    respInJson(w, runeResult)
}

func heal(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var healRequest HealRequest
    json.Unmarshal(body, &healRequest)

    switch healRequest.Target {
    case "collector":
        healCollector(healRequest, w, r)
    case "rune":
        healRune(healRequest, w, r)
    }
}

func changeStat(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    defer r.Body.Close()
    body, _ := ioutil.ReadAll(r.Body)
    var change ChangeStat
    json.Unmarshal(body, &change)

    collector, succeed := getCollectorFromData(change.GoogleId, c)
    if succeed == false {
        respFail(w, "fail to get collector with id " + change.GoogleId)
        return
    }
    totalConsume := change.Atk + change.Def + change.Int
    if collector.BonusPoint < totalConsume {
        respFail(w, "bonus point is not sufficient")
        return
    }
    collector.BonusPoint -= totalConsume
    collector.Atk += change.Atk
    collector.Def += change.Def
    collector.Int += change.Int
    succeed = insertCollector(*collector, c)
    if succeed == false {
        respFail(w, "fail to update collector")
        return
    }

    var res CollectorResult
    res.Success = "success"
    res.Collector = collector.Collector
    respInJson(w, res)
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
    case "PUT":
        setRuneOwner(w, r)
    case "GET":
        getRune(w, r)
    default:
        fmt.Fprintf(w, "not supported")
    }
}

func runesHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        getRunes(w, r)
    default:
        fmt.Fprint(w, "not supported")
    }
}

func fightHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "GET":
        fight(w, r)
    default:
        fmt.Fprint(w, "not supported")
    }
}

func healHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        heal(w, r)
    default:
        fmt.Fprint(w, "not supported")
    }
}

func changeStatHandler(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case "POST":
        changeStat(w, r)
    default:
        fmt.Fprint(w, "not supported")
    }
}
