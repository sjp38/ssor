package sor

import (
)

type Signup struct {
    GoogleId int `json:"googleId"`
}

type Collector struct {
    GoogleId int `json:"googleId"`
    Nickname string `json:"nickname"`
    Hp int `json:"hp"`
    Mp int `json:"mp"`
    Atk int `json:"atk"`
    Def int `json:"def"`
    Int int `json:"int"`
    HunterClass string `json:"hunterClass"`
}

type Result struct {
    Success string `json:"result"`
}
