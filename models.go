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

type FailResult struct {
    Success string `json:"result"`
    Reason string `json:"reason"`
}
