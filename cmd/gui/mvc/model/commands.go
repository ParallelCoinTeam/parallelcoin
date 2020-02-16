package model

import "time"

type
	DuoUIcommandsHistory struct {
		Commands       []DuoUIcommand `json:"coms"`
		CommandsNumber int            `json:"comnumber"`
	}
type
	DuoUIcommand struct {
		Com      interface{}
		ComID    string
		Category string
		Out      func()
		Time     time.Time
	}

type
	DuoUIcommandsNumber struct {
		CommandsNumber int `json:"comnumber"`
	}
