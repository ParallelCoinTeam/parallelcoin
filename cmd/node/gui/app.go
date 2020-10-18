package gui

import "github.com/p9c/pod/pkg/gui/p9"

func (ng *NodeGUI) GetAppWidget() *p9.App {
	return ng.th.App()
}