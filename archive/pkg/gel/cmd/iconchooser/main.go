package main

import (
	"sort"

	l "github.com/p9c/gio/layout"
	"github.com/atotto/clipboard"
	"github.com/p9c/interrupt"
	"github.com/p9c/qu"

	"github.com/p9c/gel/icons"

	"github.com/p9c/gel"
)

type State struct {
	*gel.Window
}

func NewState(quit qu.C) *State {
	return &State{
		Window: gel.NewWindowP9(quit),
	}
}

func main() {
	quit := qu.T()
	state := NewState(quit)
	var e error
	rootWidget := state.rootWidget()
	if e = state.Window.
		Size(48, 32).
		Title("icons chooser").
		Open().
		Run(rootWidget, func() {
			interrupt.Request()
			quit.Q()
		}, quit,
		); E.Chk(e) {

	}
}

type (
	iconInfo struct {
		name string
		data *[]byte
	}
	iconsInfo []iconInfo
)

func (ii iconsInfo) Len() int           { return len(ii) }
func (ii iconsInfo) Less(i, j int) bool { return ii[i].name < ii[j].name }
func (ii iconsInfo) Swap(i, j int)      { ii[i], ii[j] = ii[j], ii[i] }

func (s *State) rootWidget() (o l.Widget) {
	lis := s.WidgetPool.GetList()
	ow := make(iconsInfo, len(icons.Material))
	counter := 0
	for i, x := range icons.Material {
		ow[counter] = iconInfo{i, x}
		counter++
	}
	sort.Sort(ow)
	clicks := make([]*gel.Clickable, len(ow))
	for i := range ow {
		clicks[i] = s.WidgetPool.GetClickable()
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		clicks[index].SetClick(func() {
			var e error
			if e = clipboard.WriteAll("icons." + ow[index].name); E.Chk(e) {
			}
		},
		)
		return s.Flex().AlignStart().
			// Rigid(s.Inset(0.5, gel.EmptySpace(0, 0)).Fn).
			Rigid(
				s.ButtonLayout(clicks[index]).Embed(
					s.Flex().AlignMiddle().
						Rigid(
							s.Icon().
								Scale(gel.Scales["H3"]).
								Color("DocText").
								Src(ow[index].data).
								Fn,
						).
						Rigid(
							s.Body1(ow[index].name).Fn,
						).Fn,
				).Background("PanelBg").Fn,
			).
			Fn(gtx)
	}
	return s.VFlex().AlignStart().
		Rigid(s.Fill("DocBg", l.Center, 0, 0,
			s.Inset(0.5,
				s.Flex().AlignStart().
					Rigid(
						s.H4("material icons").Fn,
					).
					Rigid(s.Inset(0.5, gel.EmptySpace(0, 0)).Fn).
					Flexed(1,
						s.Body1("click to copy icon's variable name").Fn,
					).Fn,
			).Fn,
		).Fn,
		).
		Flexed(1,
			s.Inset(0.25,
				lis.Vertical().Length(len(ow)).ListElement(le).Fn,
			).Fn,
		).
		Fn
}
