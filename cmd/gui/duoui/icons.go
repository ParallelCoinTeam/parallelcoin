package duoui

import (
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/cmd/gui/components"
	"github.com/p9c/pod/pkg/log"
)

func DuoIcons() map[string]*components.DuoUIicon {
	//ics := make(map[string]*theme.DuoUIicon)
	//// Icons
	logo, err := components.NewDuoUIicon(ico.ParallelCoin)
	if err != nil {
		log.FATAL(err)
	}

	return map[string]*components.DuoUIicon{
		"Logo": logo,
	}
}
