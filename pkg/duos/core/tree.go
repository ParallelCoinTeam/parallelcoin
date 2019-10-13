package core

import (
	"github.com/p9c/lorca"
	"github.com/p9c/pod/pkg/conf"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/duos/db"
	"github.com/p9c/pod/pkg/duos/srv"
	"github.com/robfig/cron"
)

// Core
type DuOS struct {
	Cx       *conte.Xt        `json:"context"`
	Cr       *cron.Cron       `json:"cron"`
	Ui       *lorca.UI        `json:"ui"`
	DB       db.DuOSdb        `json:"database"`
	Config   *conf.DuOSconfig `json:"configuration"`
	Services srv.DuOSservices `json:"services"`
}
