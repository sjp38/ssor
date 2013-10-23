package sor

import (
)

type Collector struct {
    GoogleId string `json:"googleId"`
    Email string `json:"email"`
    ProfileUrl string `json:"profileUrl"`
    Nickname string `json:"nickname"`
    Hp int `json:"hp"`
    Mp int `json:"mp"`
    Atk int `json:"atk"`
    Def int `json:"def"`
    Int int `json:"int"`
    Exp int `json:"exp"`
    ScanCount int `json:"scanCount"`
    CollectorClass string `json:"collectorClass"`
}

type CollectorWriteResult struct {
    Success string `json:"result"`
    GoogleId string `json:"googleId"`
}

type CollectorReadResult struct {
    Success string `json:"result"`
    Collector Collector `json:"collector"`
}

type Result struct {
    Success string `json:"result"`
}
