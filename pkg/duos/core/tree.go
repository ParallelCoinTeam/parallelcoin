package core

import (
	bnd "github.com/p9c/pod/pkg/bundler"
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
	GuI *GuI             `json:"gui"`
	DbS db.DuOSdb        `json:"database"`
	CgG *conf.DuOSconfig `json:"configuration"`
	SrV srv.DuOSservices `json:"services"`
	BnD bnd.DuOSassets   `json:"assets"`
}
