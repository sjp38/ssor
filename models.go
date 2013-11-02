package sor

import (
)

type CollectorMinInfo struct {
    GoogleId string `json:"googleId"`
    Email string `json:"email"`
    ProfileUrl string `json:"profileUrl"`
    Nickname string `json:"nickname"`
    CollectorClass string `json:"collectorClass"`
}

type Collector struct {
    CollectorMinInfo
    MaxHp int `json:"maxHp"`
    Hp int `json:"hp"`
    MaxMp int `json:"maxMp"`
    Mp int `json:"mp"`
    Atk int `json:"atk"`
    Def int `json:"def"`
    Int int `json:"int"`
    Exp int `json:"exp"`
    ScanCount int `json:"scanCount"`
}

type CollectorResult struct {
    Success string `json:"result"`
    Collector Collector `json:"collector"`
}

// Model for internal - for GAE datastore
type CollectorInternal struct {
    Collector
    CreatedTime int64
    LastScannedTime int64
    LastMpConsumedTime int64
    TotalScanCount int
}

type RuneMinInfo struct {
    ISBN string `json:"ISBN"`
    OwnerGoogleId string `json:"ownerGoogleId"`
}

type Rune struct {
    RuneMinInfo
    ImageUrl string `json:"imageUrl"`
    Title string `json:"title"`
    Type string `json:"type"`
    MaxHp int `json:"maxHp"`
    Hp int `json:"hp"`
}

type RuneResult struct {
    Success string `json:"result"`
    Rune Rune `json:"rune"`
}

type RunesResult struct {
    Success string `json:"results"`
    Runes []Rune `json:"runes"`
}

type FightResult struct {
    Success string `json:"result"`
    Attacker Collector `json:"attacker"`
    Defender Collector `json:"defender"`
    Rune Rune `json:"rune"`
}

type HealRequest struct {
    Target string `json:"target"`
    Type string `json:"type"`
    GoogleId string `json:"googleId"`
    ISBN string `json:"ISBN"`
}

type FailResult struct {
    Success string `json:"result"`
    Reason string `json:"reason"`
}

type ItemType struct {
    Cover_l_url string `json:"cover_l_url"`
    Title string `json:"title"`
    InfoUrl string `json:"link"`
}

type ChannelType struct {
    Item []ItemType
}

type DaumBookSearchResult struct {
    Channel ChannelType
}
