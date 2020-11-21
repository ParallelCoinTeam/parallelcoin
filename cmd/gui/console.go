package gui

import (
	"regexp"
	"sort"
	"strconv"
	"strings"

	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	l "gioui.org/layout"

	"github.com/atotto/clipboard"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/ctl"
)

type Console struct {
	th             *p9.Theme
	output         []l.Widget
	outputList     *p9.List
	editor         *p9.Editor
	clearClickable *p9.Clickable
	clearButton    *p9.IconButton
	copyClickable  *p9.Clickable
	copyButton     *p9.IconButton
	pasteClickable *p9.Clickable
	pasteButton    *p9.IconButton
	submitFunc     func(txt string)
	clickables     []*p9.Clickable
}

var findSpaceRegexp = regexp.MustCompile(`\s+`)

func (wg *WalletGUI) ConsolePage() *Console {
	Debug("running ConsolePage")
	c := &Console{
		th:             wg.th,
		editor:         wg.th.Editor().SingleLine().Submit(true),
		clearClickable: wg.th.Clickable(),
		copyClickable:  wg.th.Clickable(),
		pasteClickable: wg.th.Clickable(),
		outputList:     wg.th.List().ScrollToEnd(),
	}
	c.submitFunc = func(txt string) {
		go func() {
			Debug("submit", txt)
			c.output = append(c.output,
				func(gtx l.Context) l.Dimensions {
					return wg.th.VFlex().
						Rigid(wg.th.Inset(0.5, p9.EmptySpace(0, 0)).Fn).
						Rigid(
							wg.th.Flex().
								Flexed(1,
									wg.th.Body1(txt).Color("DocText").Font("bariol bold").Fn,
								).
								Fn,
						).Fn(gtx)
				})
			c.editor.SetText("")
			split := strings.Split(txt, " ")
			method, args := split[0], split[1:]
			var params []interface{}
			var err error
			var result []byte
			var o string
			var errString, prev string
			for i := range args {
				params = append(params, args[i])
			}
			if method == "clear" || method == "cls" {
				// clear the list of display widgets
				c.output = c.output[:0]
				// free up the pool widgets used in the current output
				for i := range c.clickables {
					wg.th.WidgetPool.FreeClickable(c.clickables[i])
				}
				c.clickables = c.clickables[:0]
				return
			}
			if method == "help" {
				if len(args) == 0 {
					Debug("rpc called help")
					var result1, result2 []byte
					if result1, err = ctl.Call(wg.cx, false, method, params...); Check(err) {
					}
					r1 := string(result1)
					if r1, err = strconv.Unquote(r1); Check(err) {
					}
					o = r1 + "\n"
					if result2, err = ctl.Call(wg.cx, true, method, params...); Check(err) {
					}
					r2 := string(result2)
					if r2, err = strconv.Unquote(r2); Check(err) {
					}
					o += r2 + "\n"
					splitted := strings.Split(o, "\n")
					sort.Strings(splitted)
					var dedup []string
					for i := range splitted {
						if i > 0 {
							if splitted[i] != prev {
								dedup = append(dedup, splitted[i])
							}
						}
						prev = splitted[i]
					}
					o = strings.Join(dedup, "\n")
					if errString != "" {
						o += "BTCJSONError:\n"
						o += errString
					}
					splitResult := strings.Split(o, "\n")
					const maxPerWidget = 6
					for i := 0; i < len(splitResult)-maxPerWidget; i += maxPerWidget {
						sri := strings.Join(splitResult[i:i+maxPerWidget], "\n")
						c.output = append(c.output,
							func(gtx l.Context) l.Dimensions {
								return wg.th.Flex().
									Flexed(1,
										wg.th.Caption(sri).
											Color("DocText").
											Font("bariol regular").
											MaxLines(maxPerWidget).Fn,
									).
									Fn(gtx)
							})
					}
					return
				} else {
					if result, err = ctl.Call(wg.cx, false, method, params...); Check(err) {
					}
					var out string
					if out, err = strconv.Unquote(string(result)); Check(err) {
					}
					strings.ReplaceAll(out, "\t", "  ")
					Debug(out)
					splitResult := strings.Split(out, "\n")
					for i := range splitResult {
						sri := splitResult[i]
						c.output = append(c.output,
							func(gtx l.Context) l.Dimensions {
								return wg.th.Flex().Flexed(1,
									wg.th.Caption(sri).
										Color("DocText").
										Font("go regular").MaxLines(4).
										Fn,
								).Fn(gtx)
							})
					}
					return
				}
			} else {
				Debug("method", method, "args", args)
				if result, err = ctl.Call(wg.cx, false, method, params...); Check(err) {
					if result, err = ctl.Call(wg.cx, true, method, params...); Check(err) {
						c.output = append(c.output, wg.th.Caption(err.Error()).Color("Error").Fn)
						// return
					}
				}
				c.output = append(c.output, wg.console.JSONWidget("DocText", result)...)
			}
			c.outputList.JumpToEnd()
		}()
	}
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
	c.output = append(c.output, func(gtx l.Context) l.Dimensions {
		return c.th.Flex().AlignStart().Rigid(c.th.H6("Welcome to the Parallelcoin RPC console").Color("DocText").Fn).Fn(gtx)
	}, func(gtx l.Context) l.Dimensions {
		return c.th.Flex().AlignStart().Rigid(c.th.Caption("Type 'help' to get available commands and 'clear' or 'cls' to clear the screen").Color("DocText").Fn).Fn(gtx)
	})
	return c
}

func (c *Console) Fn(gtx l.Context) l.Dimensions {
	le := func(gtx l.Context, index int) l.Dimensions {
		if index >= len(c.output) || index < 0 {
			return l.Dimensions{}
		} else {
			return c.output[index](gtx)
		}
	}
	fn := c.th.VFlex().
		Flexed(0.1,
			c.th.Fill("PanelBg",
				func(gtx l.Context) l.Dimensions {
					return c.th.Inset(0.25,
						c.outputList.
							ScrollToEnd().
							End().
							Background("PanelBg").
							Color("PanelText").
							Active("Primary").
							Vertical().
							Length(len(c.output)).
							ListElement(le).
							Fn,
					).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			c.th.Fill("DocBg",
				c.th.Inset(0.25,
					c.th.Flex().
						Flexed(1,
							c.th.TextInput(c.editor.SetSubmit(c.submitFunc), "enter an rpc command").
								Color("DocText").
								Fn,
						).
						Rigid(c.copyButton.Fn).
						Rigid(c.pasteButton.Fn).
						Rigid(c.clearButton.Fn).
						Fn,
				).Fn,
			).Fn,
		).
		Fn
	return fn(gtx)
}
