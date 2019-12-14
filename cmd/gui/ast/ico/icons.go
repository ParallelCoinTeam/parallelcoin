package ico

import (
	"gioui.org/widget/material"
	"github.com/p9c/pod/pkg/log"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

func (i *DuoUIicons) DuoUIicons() {
	var err error
	i.Logo, err = material.NewIcon(ParalleCoin)
	if err != nil {
		log.FATAL(err)
	}
	i.Overview, err = material.NewIcon(icons.ActionHome)
	if err != nil {
		log.FATAL(err)
	}
	i.History, err = material.NewIcon(icons.ActionHistory)
	if err != nil {
		log.FATAL(err)
	}
	i.Network, err = material.NewIcon(icons.DeviceNetworkCell)
	if err != nil {
		log.FATAL(err)
	}
	i.Settings, err = material.NewIcon(icons.ActionSettings)
	if err != nil {
		log.FATAL(err)
	}
	return
}