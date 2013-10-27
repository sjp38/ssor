package sor

import (
)

type Collector struct {
    GoogleId string `json:"googleId"`
    Email string `json:"email"`
    ProfileUrl string `json:"profileUrl"`
    Nickname string `json:"nickname"`
    MaxHp int `json:"maxHp"`
    Hp int `json:"hp"`
    MaxMp int `json:"maxMp"`
    Mp int `json:"mp"`
    Atk int `json:"atk"`
    Def int `json:"def"`
    Int int `json:"int"`
    Exp int `json:"exp"`
    ScanCount int `json:"scanCount"`
    CollectorClass string `json:"collectorClass"`
}

type CollectorResult struct {
    Success string `json:"result"`
    Collector Collector `json:"collector"`
}

type Rune struct {
    ISBN string `json:"ISBN"`
    ImageUrl string `json:"imageUrl"`
    Title string `json:"title"`
    Type string `json:"type"`
    MaxHp int `json:"maxHp"`
    Hp int `json:"hp"`
    OwnerGoogleId string `json:"ownerGoogleId"`
}

type RuneResult struct {
    Success string `json:"result"`
    Rune Rune `json:"rune"`
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
