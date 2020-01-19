package duoui

import (
	"github.com/p9c/pod/cmd/gui/ico"
	"github.com/p9c/pod/pkg/gui/widget/parallel"
	"github.com/p9c/pod/pkg/log"
)

func DuoIcons() map[string]*parallel.DuoUIicon {
	//ics := make(map[string]*theme.DuoUIicon)
	//// Icons
	logo, err := parallel.NewDuoUIicon(ico.ParallelCoin)
	if err != nil {
		log.FATAL(err)
	}

	return map[string]*parallel.DuoUIicon{
		"Logo": logo,
	}
}
