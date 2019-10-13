package core

import (
	"github.com/p9c/pod/cmd/gui/db"
	"github.com/p9c/pod/cmd/gui/mod"
	"github.com/p9c/pod/pkg/conte"

	"github.com/robfig/cron"
	"github.com/zserge/webview"
	"time"
)

// Core
type DuOS struct {
	Cx     *conte.Xt
	Cr     *cron.Cron
	Wv     webview.WebView
	db     db.DuOSdb    `json:"db"`
	Config DuOSconfig   `json:"conf"`
	Data   mod.DuOSdata `json:"d"`
	Alert  DuOSalert    `json:"alert"`
}

type DuOSalert struct {
	Time      time.Time   `json:"time"`
	Title     string      `json:"title"`
	Message   interface{} `json:"message"`
	AlertType string      `json:"type"`
}
