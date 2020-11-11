package gui

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"

	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/atotto/clipboard"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/ctl"
)

type Console struct {
	w              l.Widget
	output         []l.Widget
	outputList     *p9.List
	editor         *p9.Editor
	clearClickable *p9.Clickable
	clearButton    *p9.IconButton
	copyClickable  *p9.Clickable
	copyButton     *p9.IconButton
	pasteClickable *p9.Clickable
	pasteButton    *p9.IconButton
}

const unusableFlags = btcjson.UFWebsocketOnly | btcjson.UFNotification

var findSpaceRegexp = regexp.MustCompile(`\s+`)

func (wg *WalletGUI) ConsolePage() *Console {
	c := &Console{
		editor:         wg.th.Editor().SingleLine().Submit(true),
		clearClickable: wg.th.Clickable(),
		copyClickable:  wg.th.Clickable(),
		pasteClickable: wg.th.Clickable(),
		outputList:     wg.th.List().ScrollToEnd(),
	}
	c.editor.SetSubmit(func(txt string) {
		Debug("submit", txt)
		c.output = append(c.output,
			func(gtx l.Context) l.Dimensions {
				return wg.th.Body1(txt).Color("DocText").Font("bariol bold").Fn(gtx)
			})
		c.editor.SetText("")
		split := strings.Split(txt, " ")
		method, args := split[0], split[1:]
		var params []interface{}
		for i := range args {
			params = append(params, args[i])
		}
		Debug("method", method, "args", args)
		var err error
		var result []byte
		if result, err = ctl.Call(wg.cx, method, params...); Check(err) {
		}
		var buf bytes.Buffer
		json.Indent(&buf, result, "", "  ")
		out := buf.String()
		Debug(out)
		splitResult := strings.Split(out, "\n")
		for i := range splitResult {
			sri := splitResult[i]
			c.output = append(c.output,
				func(gtx l.Context) l.Dimensions {
					return wg.th.Caption(sri).Color("DocText").Font("go regular").Fn(gtx)
				})
		}
	})
	clearClickableFn := func() {
		c.editor.SetText("")
		c.editor.Focus()
	}
	copyClickableFn := func() {
		go clipboard.WriteAll(c.editor.Text())
		c.editor.Focus()
	}
	pasteClickableFn := func() {
		col := c.editor.Caret.Col
		go func() {
			txt := c.editor.Text()
			var err error
			var cb string
			if cb, err = clipboard.ReadAll(); Check(err) {
			}
			cb = findSpaceRegexp.ReplaceAllString(cb, " ")
			txt = txt[:col] + cb + txt[col:]
			c.editor.SetText(txt)
			c.editor.Move(col + len(cb))
		}()
		// c.editor.Focus()
	}
	c.clearButton = wg.th.IconButton(c.clearClickable.SetClick(clearClickableFn)).
		Icon(
			wg.th.Icon().
				Color("DocText").
				Src(&icons2.ContentBackspace),
		).
		Background("Transparent").
		Inset(0.25)
	c.copyButton = wg.th.IconButton(c.copyClickable.SetClick(copyClickableFn)).
		Icon(
			wg.th.Icon().
				Color("DocText").
				Src(&icons2.ContentContentCopy),
		).
		Background("Transparent").
		Inset(0.25)
	c.pasteButton = wg.th.IconButton(c.pasteClickable.SetClick(pasteClickableFn)).
		Icon(
			wg.th.Icon().
				Color("DocText").
				Src(&icons2.ContentContentPaste),
		).
		Background("Transparent").
		Inset(0.25)
	c.w = func(gtx l.Context) l.Dimensions {
		return wg.th.VFlex().
			Flexed(1,
				wg.th.Fill("DocBg",
					// p9.EmptyMaxHeight(),
					// wg.th.H6("output area").Color("DocText").Fn,
					func(gtx l.Context) l.Dimensions {
						le := func(gtx l.Context, index int) l.Dimensions {
							return c.output[index](gtx)
						}
						return func(gtx l.Context) l.Dimensions {
							return wg.Inset(0.25,
								wg.Fill("DocBg",
									c.outputList.
										Vertical().
										// Background("DocBg").Color("DocText").Active("Primary").
										Length(len(c.output)).
										ListElement(le).
										Fn,
								).Fn,
							).
								Fn(gtx)
						}(gtx)
					},
				).Fn,
			).
			Rigid(
				wg.th.Fill("PanelBg",
					wg.th.Inset(0.25,
						wg.th.Flex().
							Flexed(1,
								wg.th.TextInput(c.editor, "enter an rpc command").
									Color("DocText").

									Fn,
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
