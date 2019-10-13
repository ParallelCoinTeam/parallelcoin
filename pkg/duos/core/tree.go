package core

import (
	"github.com/p9c/gui"
	"github.com/p9c/pod/pkg/conf"
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/duos/db"
	"github.com/p9c/pod/pkg/duos/srv"
	"github.com/robfig/cron"
	"sync"
)

// Core
type DuOS struct {
	sync.Mutex
	CtX *conte.Xt        `json:"context"`
	CrN *cron.Cron       `json:"cron"`
	GuI gui.UI           `json:"gui"`
	DbS db.DuOSdb        `json:"database"`
	CgG *conf.DuOSconfig `json:"configuration"`
	SrV srv.DuOSservices `json:"services"`
}
