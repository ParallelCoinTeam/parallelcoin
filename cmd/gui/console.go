package gui

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/atotto/clipboard"
	"golang.org/x/exp/shiny/materialdesign/icons"

	l "github.com/p9c/gio/layout"
	ctl2 "github.com/p9c/matrjoska/cmd/ctl"

	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/p9c/gel"
)

type Console struct {
	*gel.Window
	output         []l.Widget
	outputList     *gel.List
	editor         *gel.Editor
	clearClickable *gel.Clickable
	clearButton    *gel.IconButton
	copyClickable  *gel.Clickable
	copyButton     *gel.IconButton
	pasteClickable *gel.Clickable
	pasteButton    *gel.IconButton
	submitFunc     func(txt string)
	clickables     []*gel.Clickable
}

var findSpaceRegexp = regexp.MustCompile(`\s+`)

func (wg *WalletGUI) ConsolePage() *Console {
	D.Ln("running ConsolePage")
	c := &Console{
		Window:         wg.Window,
		editor:         wg.Editor().SingleLine().Submit(true),
		clearClickable: wg.Clickable(),
		copyClickable:  wg.Clickable(),
		pasteClickable: wg.Clickable(),
		outputList:     wg.List().ScrollToEnd(),
	}
	c.submitFunc = func(txt string) {
		go func() {
			D.Ln("submit", txt)
			c.output = append(
				c.output,
				func(gtx l.Context) l.Dimensions {
					return wg.VFlex().
						Rigid(wg.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
						Rigid(
							wg.Flex().
								Flexed(
									1,
									wg.Body2(txt).Color("DocText").Font("bariol bold").Fn,
								).
								Fn,
						).Fn(gtx)
				},
			)
			c.editor.SetText("")
			split := strings.Split(txt, " ")
			method, args := split[0], split[1:]
			var params []interface{}
			var e error
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
					wg.WidgetPool.FreeClickable(c.clickables[i])
				}
				c.clickables = c.clickables[:0]
				return
			}
			if method == "help" {
				if len(args) == 0 {
					D.Ln("rpc called help")
					var result1, result2 []byte
					if result1, e = ctl2.Call(wg.cx.Config, false, method, params...); E.Chk(e) {
					}
					r1 := string(result1)
					if r1, e = strconv.Unquote(r1); E.Chk(e) {
					}
					o = r1 + "\n"
					if result2, e = ctl2.Call(wg.cx.Config, true, method, params...); E.Chk(e) {
					}
					r2 := string(result2)
					if r2, e = strconv.Unquote(r2); E.Chk(e) {
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
						c.output = append(
							c.output,
							func(gtx l.Context) l.Dimensions {
								return wg.Flex().
									Flexed(
										1,
										wg.Caption(sri).
											Color("DocText").
											Font("bariol regular").
											MaxLines(maxPerWidget).Fn,
									).
									Fn(gtx)
							},
						)
					}
					return
				} else {
					var out string
					var isErr bool
					if result, e = ctl2.Call(wg.cx.Config, false, method, params...); E.Chk(e) {
						isErr = true
						out = e.Error()
						I.Ln(out)
						// if out, e = strconv.Unquote(); E.Chk(e) {
						// }
					} else {
						if out, e = strconv.Unquote(string(result)); E.Chk(e) {
						}
					}
					strings.ReplaceAll(out, "\t", "  ")
					D.Ln(out)
					splitResult := strings.Split(out, "\n")
					outputColor := "DocText"
					if isErr {
						outputColor = "Danger"
					}
					for i := range splitResult {
						sri := splitResult[i]
						c.output = append(
							c.output,
							func(gtx l.Context) l.Dimensions {
								return c.Theme.Flex().AlignStart().
									Rigid(
										wg.Body2(sri).
											Color(outputColor).
											Font("go regular").MaxLines(4).
											Fn,
									).
									Fn(gtx)
							},
						)
					}
					return
				}
			} else {
				D.Ln("method", method, "args", args)
				if result, e = ctl2.Call(wg.cx.Config, false, method, params...); E.Chk(e) {
					var errR string
					if result, e = ctl2.Call(wg.cx.Config, true, method, params...); E.Chk(e) {
						if e != nil {
							errR = e.Error()
						}
						c.output = append(
							c.output, c.Theme.Flex().AlignStart().
								Rigid(wg.Body2(errR).Color("Danger").Fn).Fn,
						)
						return
					}
					if e != nil {
						errR = e.Error()
					}
					c.output = append(
						c.output, c.Theme.Flex().AlignStart().
							Rigid(
								wg.Body2(errR).Color("Danger").Fn,
							).Fn,
					)
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
		go func() {
			if e := clipboard.WriteAll(c.editor.Text()); E.Chk(e) {
			}
		}()
		c.editor.Focus()
	}
	pasteClickableFn := func() {
		// // col := c.editor.Caret.Col
		// go func() {
		// 	txt := c.editor.Text()
		// 	var e error
		// 	var cb string
		// 	if cb, e = clipboard.ReadAll(); E.Chk(e) {
		// 	}
		// 	cb = findSpaceRegexp.ReplaceAllString(cb, " ")
		// 	txt = txt[:col] + cb + txt[col:]
		// 	c.editor.SetText(txt)
		// 	c.editor.Move(col + len(cb))
		// }()
	}
	c.clearButton = wg.IconButton(c.clearClickable.SetClick(clearClickableFn)).
		Icon(
			wg.Icon().
				Color("DocText").
				Src(&icons2.ContentBackspace),
		).
		Background("").
		ButtonInset(0.25)
	c.copyButton = wg.IconButton(c.copyClickable.SetClick(copyClickableFn)).
		Icon(
			wg.Icon().
				Color("DocText").
				Src(&icons2.ContentContentCopy),
		).
		Background("").
		ButtonInset(0.25)
	c.pasteButton = wg.IconButton(c.pasteClickable.SetClick(pasteClickableFn)).
		Icon(
			wg.Icon().
				Color("DocText").
				Src(&icons2.ContentContentPaste),
		).
		Background("").
		ButtonInset(0.25)
	c.output = append(
		c.output, func(gtx l.Context) l.Dimensions {
			return c.Theme.Flex().AlignStart().Rigid(c.H6("Welcome to the Parallelcoin RPC console").Color("DocText").Fn).Fn(gtx)
		}, func(gtx l.Context) l.Dimensions {
			return c.Theme.Flex().AlignStart().Rigid(c.Caption("Type 'help' to get available commands and 'clear' or 'cls' to clear the screen").Color("DocText").Fn).Fn(gtx)
		},
	)
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
	fn := c.Theme.VFlex().
		Flexed(
			0.1,
			c.Fill(
				"PanelBg", l.Center, c.TextSize.V, 0, func(gtx l.Context) l.Dimensions {
					return c.Inset(
						0.25,
						c.outputList.
							ScrollToEnd().
							End().
							Background("PanelBg").
							Color("DocBg").
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
			c.Fill(
				"DocBg", l.Center, c.TextSize.V, 0, c.Inset(
					0.25,
					c.Theme.Flex().
						Flexed(
							1,
							c.TextInput(c.editor.SetSubmit(c.submitFunc), "enter an rpc command").
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

type JSONElement struct {
	key   string
	value interface{}
}

type JSONElements []JSONElement

func (je JSONElements) Len() int {
	return len(je)
}

func (je JSONElements) Less(i, j int) bool {
	return je[i].key < je[j].key
}

func (je JSONElements) Swap(i, j int) {
	je[i], je[j] = je[j], je[i]
}

func GetJSONElements(in map[string]interface{}) (je JSONElements) {
	for i := range in {
		je = append(
			je, JSONElement{
				key:   i,
				value: in[i],
			},
		)
	}
	sort.Sort(je)
	return
}

func (c *Console) getIndent(n int, size float32, widget l.Widget) (out l.Widget) {
	o := c.Theme.Flex()
	for i := 0; i < n; i++ {
		o.Rigid(c.Inset(size/2, gel.EmptySpace(0, 0)).Fn)
	}
	o.Rigid(widget)
	out = o.Fn
	return
}

func (c *Console) JSONWidget(color string, j []byte) (out []l.Widget) {
	var ifc interface{}
	var e error
	if e = json.Unmarshal(j, &ifc); E.Chk(e) {
	}
	return c.jsonWidget(color, 0, "", ifc)
}

func (c *Console) jsonWidget(color string, depth int, key string, in interface{}) (out []l.Widget) {
	switch in.(type) {
	case []interface{}:
		if key != "" {
			out = append(
				out, c.getIndent(
					depth, 1,
					func(gtx l.Context) l.Dimensions {
						return c.Body2(key).Font("bariol bold").Color(color).Fn(gtx)
					},
				),
			)
		}
		D.Ln("got type []interface{}")
		res := in.([]interface{})
		if len(res) == 0 {
			out = append(
				out, c.getIndent(
					depth+1, 1,
					func(gtx l.Context) l.Dimensions {
						return c.Body2("[]").Color(color).Fn(gtx)
					},
				),
			)
		} else {
			for i := range res {
				// D.S(res[i])
				out = append(out, c.jsonWidget(color, depth+1, fmt.Sprint(i), res[i])...)
			}
		}
	case map[string]interface{}:
		if key != "" {
			out = append(
				out, c.getIndent(
					depth, 1,
					func(gtx l.Context) l.Dimensions {
						return c.Body2(key).Font("bariol bold").Color(color).Fn(gtx)
					},
				),
			)
		}
		D.Ln("got type map[string]interface{}")
		res := in.(map[string]interface{})
		je := GetJSONElements(res)
		// D.S(je)
		if len(res) == 0 {
			out = append(
				out, c.getIndent(
					depth+1, 1,
					func(gtx l.Context) l.Dimensions {
						return c.Body2("{}").Color(color).Fn(gtx)
					},
				),
			)
		} else {
			for i := range je {
				D.S(je[i])
				out = append(out, c.jsonWidget(color, depth+1, je[i].key, je[i].value)...)
			}
		}
	case JSONElement:
		res := in.(JSONElement)
		key = res.key
		switch res.value.(type) {
		case string:
			D.Ln("got type string")
			res := res.value.(string)
			clk := c.Theme.WidgetPool.GetClickable()
			out = append(
				out,
				c.jsonElement(
					key, color, depth, func(gtx l.Context) l.Dimensions {
						return c.Theme.Flex().
							Rigid(c.Body2("\"" + res + "\"").Color(color).Fn).
							Rigid(c.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
							Rigid(
								c.IconButton(clk).
									Background("").
									ButtonInset(0).
									Color(color).
									Icon(c.Icon().Color("DocBg").Scale(1).Src(&icons.ContentContentCopy)).
									SetClick(
										func() {
											go func() {
												if e := clipboard.WriteAll(res); E.Chk(e) {
												}
											}()
										},
									).Fn,
							).Fn(gtx)
					},
				),
			)
		case float64:
			D.Ln("got type float64")
			res := res.value.(float64)
			clk := c.Theme.WidgetPool.GetClickable()
			out = append(
				out,
				c.jsonElement(
					key, color, depth, func(gtx l.Context) l.Dimensions {
						return c.Theme.Flex().
							Rigid(c.Body2(fmt.Sprint(res)).Color(color).Fn).
							Rigid(c.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
							Rigid(
								c.IconButton(clk).
									Background("").
									ButtonInset(0).
									Color(color).
									Icon(c.Icon().Color("DocBg").Scale(1).Src(&icons.ContentContentCopy)).
									SetClick(
										func() {
											go func() {
												if e := clipboard.WriteAll(fmt.Sprint(res)); E.Chk(e) {
												}
											}()
										},
									).Fn,
							).Fn(gtx)
						// return c.th.ButtonLayout(clk).Embed(c.th.Body2().Color(color).Fn).Fn(gtx)
					},
				),
			)
		case bool:
			D.Ln("got type bool")
			res := res.value.(bool)
			out = append(
				out,
				c.jsonElement(
					key, color, depth, func(gtx l.Context) l.Dimensions {
						return c.Body2(fmt.Sprint(res)).Color(color).Fn(gtx)
					},
				),
			)
		}
	case string:
		D.Ln("got type string")
		res := in.(string)
		clk := c.Theme.WidgetPool.GetClickable()
		out = append(
			out,
			c.jsonElement(
				key, color, depth, func(gtx l.Context) l.Dimensions {
					return c.Theme.Flex().
						Rigid(c.Body2("\"" + res + "\"").Color(color).Fn).
						Rigid(c.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
						Rigid(
							c.IconButton(clk).
								Background("").
								ButtonInset(0).
								Color(color).
								Icon(c.Icon().Color("DocBg").Scale(1).Src(&icons.ContentContentCopy)).
								SetClick(
									func() {
										go func() {
											if e := clipboard.WriteAll(res); E.Chk(e) {
											}
										}()
									},
								).Fn,
						).Fn(gtx)
				},
			),
		)
	case float64:
		D.Ln("got type float64")
		res := in.(float64)
		clk := c.Theme.WidgetPool.GetClickable()
		out = append(
			out,
			c.jsonElement(
				key, color, depth, func(gtx l.Context) l.Dimensions {
					return c.Theme.Flex().
						Rigid(c.Body2(fmt.Sprint(res)).Color(color).Fn).
						Rigid(c.Inset(0.25, gel.EmptySpace(0, 0)).Fn).
						Rigid(
							c.IconButton(clk).
								Background("").
								ButtonInset(0).
								Color(color).
								Icon(c.Icon().Color("DocBg").Scale(1).Src(&icons.ContentContentCopy)).
								SetClick(
									func() {
										go func() {
											if e := clipboard.WriteAll(fmt.Sprint(res)); E.Chk(e) {
											}
										}()
									},
								).Fn,
						).Fn(gtx)
					// return c.th.ButtonLayout(clk).Embed(c.th.Body2(fmt.Sprint(res)).Color(color).Fn).Fn(gtx)
				},
			),
		)
	case bool:
		D.Ln("got type bool")
		res := in.(bool)
		out = append(
			out,
			c.jsonElement(
				key, color, depth, func(gtx l.Context) l.Dimensions {
					return c.Body2(fmt.Sprint(res)).Color(color).Fn(gtx)
				},
			),
		)
	default:
		D.S(in)
	}
	return
}

func (c *Console) jsonElement(key, color string, depth int, w l.Widget) l.Widget {
	return func(gtx l.Context) l.Dimensions {
		return c.Theme.Flex().
			Rigid(
				c.getIndent(
					depth, 1,
					c.Body2(key).Font("bariol bold").Color(color).Fn,
				),
			).
			Rigid(c.Inset(0.125, gel.EmptySpace(0, 0)).Fn).
			Rigid(w).
			Fn(gtx)
	}
}
