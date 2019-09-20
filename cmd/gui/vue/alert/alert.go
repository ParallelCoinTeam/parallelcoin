// +build !headless

package alert

import (
	"fmt"
	"time"
)

type DuoVUEalert struct {
	Time      time.Time   `json:"time"`
	Alert     interface{} `json:"msg"`
	AlertType string      `json:"msgtype"`
}

var Alert DuoVUEalert

// GetMsg loads the message variable
func (m *DuoVUEalert) GetAlert() {
	m.Time = Alert.Time
	m.Alert = Alert.Alert
	m.AlertType = Alert.AlertType
	fmt.Println("Bozidar:")
	fmt.Println("Alert :")
	fmt.Println("Alert :", m)
	fmt.Println("Alert :")
	fmt.Println("Alert :")
	fmt.Println("Bozidar:")
}

func PushAlert(m interface{}, t string) {
	Alert.Time = time.Now()
	Alert.Alert = m
	Alert.AlertType = t
}
