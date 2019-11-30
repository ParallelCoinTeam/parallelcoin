package gui

import (
	"github.com/p9c/pod/pkg/conte"
	"github.com/p9c/pod/pkg/log"
)

type DuOSnav struct {
	cx     *conte.Xt
	Screen string `json:"screen"`
	Config DuOSsettings `json:"config"`
}

func
(nav *DuOSnav) GetScreen(s string) {
	nav.Screen = s
	log.INFO("NAV:VAR->", s)
	switch s {
	case "PageOverview":
		log.INFO("NAV:VAR->", s)
	case "PageHistory":
		log.INFO("NAV:VAR->", s)
	case "PageAddressBook":
		log.INFO("NAV:VAR->", s)
	case "PageExplorer":
		log.INFO("NAV:VAR->", s)
	case "PageSettings":
		log.INFO("NAV:VAR->", s)
		//nav.Config.cx = nav.cx
		//nav.Config.GetCoreSettings()
		log.INFO("NAV:SETTINGS:VAR->", nav.Config.Daemon.Config)
	}
}
