package core

import (
	"time"
)

type DuOSalert struct {
	Time      time.Time   `json:"time"`
	Title     string      `json:"title"`
	Message   interface{} `json:"message"`
	AlertType string      `json:"type"`
}
