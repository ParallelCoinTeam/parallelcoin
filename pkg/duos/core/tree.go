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
	CtX *conte.Xt        `json:"context"`
	CrN *cron.Cron       `json:"cron"`
	GuI lorca.UI         `json:"__OLDgui"`
	DbS db.DuOSdb        `json:"database"`
	CgG *conf.DuOSconfig `json:"configuration"`
	SrV srv.DuOSservices `json:"services"`
}
