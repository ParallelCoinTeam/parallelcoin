package monitor

import (
	"fmt"

	"gioui.org/layout"
	"gioui.org/unit"

	"github.com/p9c/pod/pkg/gelook"
)

func (st *State) DuoUIheader() layout.FlexChild {
	return Rigid(func() {
		st.FlexH(Rigid(func() {
			cs := st.Gtx.Constraints
			st.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
			var (
				textSize, iconSize       = 64, 64
				width, height            = 72, 72
				paddingV, paddingH       = 8, 8
				insetSize, textInsetSize = 16, 24
				closeInsetSize           = 4
			)
			if st.WindowWidth < 1024 || st.WindowHeight < 1280 {
				textSize, iconSize = 24, 24
				width, height = 32, 32
				paddingV, paddingH = 8, 8
				insetSize = 10
				textInsetSize = 16
				closeInsetSize = 4
			}
			st.FlexH(Rigid(func() {
				st.Inset(insetSize,
					func() {
						var logoMeniItem gelook.DuoUIbutton
						logoMeniItem = st.Theme.DuoUIbutton(
							"", "",
							"", st.Theme.Colors["PanelBg"],
							"", "",
							"logo", st.Theme.Colors["PanelText"],
							textSize, iconSize,
							width, height,
							paddingV, paddingH)
						for st.LogoButton.Clicked(st.Gtx) {
							st.FlipTheme()
							st.SaveConfig()
						}
						logoMeniItem.IconLayout(st.Gtx, st.LogoButton)
					},
				)
			}), Rigid(func() {
				st.Inset(textInsetSize, func() {
					t := st.Theme.DuoUIlabel(unit.Dp(float32(
						textSize)),
						"monitor")
					t.Color = st.Theme.Colors["PanelText"]
					t.Layout(st.Gtx)
				},
				)
			}), Spacer(), Rigid(func() {
				st.Inset(closeInsetSize*2, func() {
					t := st.Theme.DuoUIlabel(unit.Dp(float32(24)),
						fmt.Sprintf("%dx%d",
							st.WindowWidth,
							st.WindowHeight))
					t.Color = st.Theme.Colors["PanelText"]
					t.Font.Typeface = st.Theme.Fonts["Primary"]
					t.Layout(st.Gtx)
				})
			}), Rigid(func() {
				st.Inset(closeInsetSize, func() {
					st.IconButton("closeIcon", "PanelText",
						"PanelBg", st.CloseButton)
					for st.CloseButton.Clicked(st.Gtx) {
						L.Debug("close button clicked")
						st.SaveConfig()
						close(st.Ctx.KillAll)
					}
				})
			}),
			)
		}),
		)
	})
}
