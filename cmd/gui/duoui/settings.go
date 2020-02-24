package duoui

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/mvc/theme"
	"github.com/p9c/pod/pkg/log"
)

var (
	groupsList = &layout.List{
		Axis:      layout.Horizontal,
		Alignment: layout.Start,
	}
	fieldsList = &layout.List{
		Axis: layout.Vertical,
	}
)

func (ui *DuoUI) settingsFieldLabel(f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
			name := ui.ly.Theme.H6(fmt.Sprint(f.Field.Label))
			name.Font.Typeface = ui.ly.Theme.Font.Primary
			name.Layout(ui.ly.Context)
		})
	}
}

func (ui *DuoUI) settingsFieldDescription(f *Field) func() {
	return func() {
		layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
			desc := ui.ly.Theme.Body2(fmt.Sprint(f.Field.Description))
			desc.Font.Typeface = ui.ly.Theme.Font.Primary
			desc.Layout(ui.ly.Context)
		})
	}
}

func (ui *DuoUI) headerSettings() func() {
	return func() {
		layout.Flex{Spacing: layout.SpaceAround}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				t := ui.ly.Theme.H6(ui.rc.Settings.Tabs.Current)
				t.Font.Typeface = ui.ly.Theme.Font.Primary
				t.Color = theme.HexARGB(ui.ly.Theme.Color.Light)
				t.Alignment = text.Start
				t.Layout(ui.ly.Context)
			}),
			layout.Rigid(func() {
				groupsNumber := len(ui.rc.Settings.Daemon.Schema.Groups)
				groupsList.Layout(ui.ly.Context, groupsNumber, func(i int) {
					layout.UniformInset(unit.Dp(0)).Layout(ui.ly.Context, func() {
						i = groupsNumber - 1 - i
						t := ui.rc.Settings.Daemon.Schema.Groups[i]
						txt := fmt.Sprint(t.Legend)
						for ui.rc.Settings.Tabs.TabsList[txt].Clicked(ui.ly.Context) {
							ui.rc.Settings.Tabs.Current = txt
							log.INFO("unutra: ", txt)
						}
						ui.ly.Theme.DuoUIbutton(ui.ly.Theme.Font.Primary, txt, ui.ly.Theme.Color.Dark, "ff989898", "", ui.ly.Theme.Color.Dark, 16, 0, 80, 32, 4, 4).Layout(ui.ly.Context, ui.rc.Settings.Tabs.TabsList[txt])
					})
				})
			}))
	}
}

func (ui *DuoUI) settingsBody() func() {
	return func() {
		for _, fields := range ui.rc.Settings.Daemon.Schema.Groups {
			if fmt.Sprint(fields.Legend) == ui.rc.Settings.Tabs.Current {
				fieldsList.Layout(ui.ly.Context, len(fields.Fields), func(il int) {
					il = len(fields.Fields) - 1 - il
					tl := Field{
						Field: &fields.Fields[il],
					}
					layout.Flex{
						Axis: layout.Vertical,
					}.Layout(ui.ly.Context,
						layout.Rigid(ui.settingsItemRow(&tl)),
						layout.Rigid(func() {
							cs := ui.ly.Context.Constraints
							theme.DuoUIdrawRectangle(ui.ly.Context, cs.Width.Max, 1, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
						}))
				})
			}
		}
	}
}

func (ui *DuoUI) settingsItemRow(f *Field) func() {
	return func() {
		layout.Flex{}.Layout(ui.ly.Context,
			layout.Rigid(func() {
				theme.DuoUIdrawRectangle(ui.ly.Context, 30, 3, ui.ly.Theme.Color.Dark, [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
			}),
			layout.Flexed(0.62, func() {
				layout.Flex{
					Axis:    layout.Vertical,
					Spacing: 10,
				}.Layout(ui.ly.Context,
					layout.Rigid(ui.settingsFieldLabel(f)),
					layout.Rigid(ui.settingsFieldDescription(f)),
				)
			}),
			layout.Flexed(0.38, ui.InputField(f)),
		)
	}
}
