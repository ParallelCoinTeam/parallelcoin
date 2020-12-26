package gui

import (
	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) buttonText(b *p9.Clickable, label string, click func()) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = int(wg.th.TextSize.Scale(10).V)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		return wg.th.ButtonLayout(b).Embed(
			func(gtx l.Context) l.Dimensions {
				background := "DocText"
				color := "DocBg"
				var inPad, outPad float32 = 0.5, 0
				return wg.th.Inset(outPad,
					wg.th.Fill(background,
						wg.th.Flex().
							Flexed(1,
								wg.th.Inset(inPad,
									wg.th.Caption(label).
										Color(color).
										Fn,
								).Fn,
							).Fn,
					).Fn,
				).Fn(gtx)
			},
		).
			Background("Transparent").
			SetClick(click).
			Fn(gtx)
	}
}

func (wg *WalletGUI) buttonIcon(b *p9.Clickable, page string, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := wg.App.TitleBarBackgroundGet()
		color := wg.App.MenuColorGet()
		if wg.App.ActivePageGet() == page {
			color = "PanelText"
			background = "PanelBg"
		}
		ic := wg.th.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return wg.th.Flex().Rigid(
			// wg.Inset(0.25,
			wg.th.ButtonLayout(b).
				CornerRadius(0).
				Embed(
					wg.th.Inset(0.375,
						ic,
					).Fn,
				).
				Background(background).
				SetClick(
					func() {
						if wg.App.MenuOpen {
							wg.App.MenuOpen = false
						}
						wg.App.ActivePage(page)
					}).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) buttonIconText(b *p9.Clickable, label string, ico *[]byte, onClick func()) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		ic := wg.th.Icon().
			Scale(1).
			Color("DocText").
			Src(ico).
			Fn
		return wg.th.Flex().Rigid(
			// wg.Inset(0.25,
			wg.th.ButtonLayout(b).
				CornerRadius(0).
				Embed(
					wg.th.Flex().AlignMiddle().
						Rigid(wg.th.Inset(0, ic).Fn).
						Rigid(
							wg.th.Caption(label).Color("DocText").Fn,
						).
						Fn,
				).
				Background("DocBg").
				SetClick(onClick).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}
