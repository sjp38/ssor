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
    Level int `json:"level"`
    MaxHp int `json:"maxHp"`
    Hp int `json:"hp"`
    MaxMp int `json:"maxMp"`
    Mp int `json:"mp"`
    Atk int `json:"atk"`
    Def int `json:"def"`
    Int int `json:"int"`
    Exp int `json:"exp"`
    ExpToNext int `json:"expToNextLevel"`
    ScanCount int `json:"scanCount"`
    BonusPoint int `json:"bonusPoint"`
    CreatedTime int64 `json:"createdTime"`
    LastScannedTime int64 `json:"lastScannedTime"`
    LastMpConsumedTime int64 `json:"lastMpConsunedTime"`
    TotalScanCount int `json:"totalScanCount"`
    GcmIds []string `json:"gcmIds"`
}

type CollectorResult struct {
    Success string `json:"result"`
    Collector Collector `json:"collector"`
}

// Model for internal - for GAE datastore
type CollectorInternal struct {
    Collector
}

type RuneMinInfo struct {
    ISBN string `json:"ISBN"`
    OwnerGoogleId string `json:"ownerGoogleId"`
}

type Rune struct {
    RuneMinInfo
    ImageUrl string `json:"imageUrl"`
    ThumbnailUrl string `json:"thumbnailUrl"`
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

type ChangeStat struct {
    GoogleId string `json:"googleId"`
    Atk int `json:"atk"`
    Def int `json:"def"`
    Int int `json:"int"`
}

type GcmId struct {
    GoogleId string `json:"googleId"`
    GcmId string `json:"gcmId"`
}

type GcmPushData struct {
    Type string `json:"type"`
    ISBN string `json:"ISBN"`
    PeerGoogleId string `json:"peerGoogleId"`
}

type GcmPush struct {
    RegistrationIds []string `json:"registration_ids"`
    Data GcmPushData `json:"data"`
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
