package model

import "github.com/p9c/pod/cmd/gui/mvc/controller"

type
	DuoUIcommandsHistory struct {
		Commands       []controller.Command `json:"coms"`
		CommandsNumber int            `json:"comnumber"`
	}
//type
//	DuoUIcommand struct {
//		Com      interface{}
//		ComID    string
//		Category string
//		Out      func()
//		Time     time.Time
//	}

type
	DuoUIcommandsNumber struct {
		CommandsNumber int `json:"comnumber"`
	}
