package model

import "time"

type DuoUIconsoleHistory struct {
	Commands       []DuoUIconsoleCommand `json:"coms"`
	CommandsNumber int                   `json:"comnumber"`
}
type DuoUIconsoleCommand struct {
	Com      interface{}
	ComID    string
	Category string
	Out      string
	Time     time.Time
}

type DuoUIconsoleCommandsNumber struct {
	CommandsNumber int `json:"comnumber"`
}
