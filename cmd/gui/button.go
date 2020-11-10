package gui

import (
	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/gui/p9"
)

func (wg *WalletGUI) buttonText(b *p9.Clickable, label string, click func()) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		gtx.Constraints.Max.X = int(wg.TextSize.Scale(10).V)
		gtx.Constraints.Min.X = gtx.Constraints.Max.X

		return wg.ButtonLayout(b).Embed(
			func(gtx l.Context) l.Dimensions {
				background := "DocText"
				color := "DocBg"
				var inPad, outPad float32 = 0.5, 0
				return wg.Inset(outPad,
					wg.Fill(background,
						wg.Flex().
							Flexed(1,
								wg.Inset(inPad,
									wg.Caption(label).
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

func (wg *WalletGUI) buttonIcon(b *p9.Clickable, label string, ico *[]byte) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		background := wg.TitleBarBackgroundGet()
		color := wg.MenuColorGet()
		if wg.ActivePageGet() == label {
			color = "PanelText"
			background = "PanelBg"
		}
		ic := wg.Icon().
			Scale(p9.Scales["H5"]).
			Color(color).
			Src(ico).
			Fn
		return wg.Flex().Rigid(
			// wg.Inset(0.25,
			wg.ButtonLayout(b).
				CornerRadius(0).
				Embed(
					wg.Inset(0.375,
						ic,
					).Fn,
				).
				Background(background).
				SetClick(
					func() {
						if wg.MenuOpen {
							wg.MenuOpen = false
						}
						wg.ActivePage(label)
					}).
				Fn,
			// ).Fn,
		).Fn(gtx)
	}
}

func (wg *WalletGUI) buttonIconText(b *p9.Clickable, label string, ico *[]byte, onClick func()) func(gtx l.Context) l.Dimensions {
	return func(gtx l.Context) l.Dimensions {
		ic := wg.Icon().
			Scale(1).
			Color("DocText").
			Src(ico).
			Fn
		return wg.Flex().Rigid(
			// wg.Inset(0.25,
			wg.ButtonLayout(b).
				CornerRadius(0).
				Embed(
					wg.th.Flex().AlignMiddle().
						Rigid(wg.Inset(0, ic).Fn).
						Rigid(
							wg.Caption(label).Color("DocText").Fn,
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
