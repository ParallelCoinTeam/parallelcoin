package gui

import (
	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/pod/pkg/gui/p9"
)

type Console struct {
	w              l.Widget
	editor         *p9.Editor
	clearClickable *p9.Clickable
	clearButton    *p9.IconButton
	copyClickable  *p9.Clickable
	copyButton     *p9.IconButton
	pasteClickable *p9.Clickable
	pasteButton    *p9.IconButton
}

func (wg *WalletGUI) ConsolePage() *Console {
	c := &Console{
		editor:         wg.th.Editor().SingleLine(),
		clearClickable: wg.th.Clickable(),
		copyClickable:  wg.th.Clickable(),
		pasteClickable: wg.th.Clickable(),
	}
	c.clearButton = wg.th.IconButton(c.clearClickable).
		Icon(
			wg.th.Icon().
				Color("DocText").
				Src(&icons2.ContentBackspace),
		).Background("Transparent").Inset(0.25)
	c.copyButton = wg.th.IconButton(c.copyClickable).
		Icon(
			wg.th.Icon().
				Color("DocText").
				Src(&icons2.ContentContentCopy),
		).Background("Transparent").Inset(0.25)
	c.pasteButton = wg.th.IconButton(c.pasteClickable).
		Icon(
			wg.th.Icon().
				Color("DocText").
				Src(&icons2.ContentContentPaste),
		).Background("Transparent").Inset(0.25)
	c.w = func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Flexed(1,
				wg.th.Fill("DocBg",
					p9.EmptyMaxHeight(),
				).Fn,
			).
			Rigid(
				wg.th.Fill("PanelBg",
					wg.th.Inset(0.25,
						wg.th.Flex().
							Flexed(1,
								wg.th.TextInput(c.editor, "enter an rpc command").Color("DocText").Fn,
							).
							Rigid(c.copyButton.Fn).
							Rigid(c.pasteButton.Fn).
							Rigid(c.clearButton.Fn).
							Fn,
					).Fn,
				).Fn,
			).Fn(gtx)
	}
	return c
}

func (c *Console) Fn(gtx l.Context) l.Dimensions {
	return c.w(gtx)
}
