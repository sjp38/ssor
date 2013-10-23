package sor

import (
)

type Collector struct {
    GoogleId int `json:"googleId"`
    Nickname string `json:"nickname"`
    Hp int `json:"hp"`
    Mp int `json:"mp"`
    Atk int `json:"atk"`
    Def int `json:"def"`
    Int int `json:"int"`
    CollectorClass string `json:"collectorClass"`
}

type Result struct {
    Success string `json:"result"`
}
