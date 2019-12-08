package mod

import (
	"gioui.org/layout"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type DuoUIlayouts struct {
	View *layout.Flex
	Main *layout.Flex
	Menu *layout.Flex
}

type DuoUIicons struct {
	Logo     *material.Icon
	Overview *material.Icon
	History  *material.Icon
	Network  *material.Icon
	Settings *material.Icon
}

type DuoUIbuttons struct {
	Logo *widget.Button
}

//  Vue component model
type DuoUIcom struct {
	Name        string `json:"name"`
	ID          string `json:"id"`
	Version     string `json:"ver"`
	Description string `json:"desc"`
	State       string `json:"state"`
	Image       string `json:"img"`
	URL         string `json:"url"`
	CompType    string `json:"comtype"`
	SubType     string `json:"subtype"`
	Js          string `json:"js"`
	Html        string `json:"html"`
	Css         string `json:"css"`
}
