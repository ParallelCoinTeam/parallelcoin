package bind

import "github.com/p9c/pod/pkg/log"

type DuOSnav struct {
	Screen string `json:"screen"`
}

func
(nav *DuOSnav) GetScreen(s string) {
	nav.Screen = s
	log.INFO("NAV:VAR->", s)
}