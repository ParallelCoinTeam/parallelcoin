package duoui

import (
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/cmd/gui/theme"
	"github.com/p9c/pod/pkg/log"
)

func DuoIcons() map[string]*theme.DuoUIicon {
	//ics := make(map[string]*theme.DuoUIicon)
	//// Icons
	logo, err := theme.NewDuoUIicon(ico.ParallelCoin)
	if err != nil {
		log.FATAL(err)
	}

	return map[string]*theme.DuoUIicon{
		"Logo": logo,
	}
}
