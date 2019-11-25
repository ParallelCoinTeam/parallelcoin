package gui

type DuOSnav struct {
	rc *rcvar
}

func
(nav *DuOSnav)GetScreen(screen string)  {
	nav.rc.screen = screen
}
