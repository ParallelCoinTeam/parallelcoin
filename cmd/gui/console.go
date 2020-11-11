package gui

import (
	"bytes"
	"encoding/json"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/atotto/clipboard"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/ctl"
)

type Console struct {
	// w              l.Widget
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
}

const unusableFlags = btcjson.UFWebsocketOnly | btcjson.UFNotification

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
					return wg.th.Body1(txt).Color("DocText").Font("bariol bold").Fn(gtx)
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
				c.output = c.output[:0]
			}
			if method == "help" {
				if len(args) == 0 {
					Debug("rpc called help")
					// var subcommand string
					// if len(args) > 0 {
					// 	subcommand = args[0]
					// }
					// cmd := &btcjson.HelpCmd{Command: &subcommand}

					*wg.cx.Config.Wallet = true
					var result1, result2 []byte
					if result1, err = ctl.Call(wg.cx, method, params...); Check(err) {
					}
					r1 := string(result1)
					// Debug(r1)
					if r1, err = strconv.Unquote(r1); Check(err) {
					}
					o = r1 + "\n"

					*wg.cx.Config.Wallet = false
					if result2, err = ctl.Call(wg.cx, method, params...); Check(err) {
					}
					r2 := string(result2)
					// Debug(r2)
					// r2 = strings.ReplaceAll(r2, "`", "")
					if r2, err = strconv.Unquote(r2); Check(err) {
					}
					o += r2 + "\n"

					// Debug(o)
					// if res, err = chainrpc.RPCHandlers["help"].Fn(rpcSrv, cmd, nil); Check(err) {
					// 	errString += fmt.Sprintln(err)
					// }
					// o += fmt.Sprintln(res)
					// if res, err = lrpcHnd["help"].Handler(cmd, ws, cc); Check(err) {
					// 	errString += fmt.Sprintln(err)
					// }
					// o += fmt.Sprintln(res)
					// o = strings.ReplaceAll(o, "\\\"", "'")

					// strings.ReplaceAll(o, "\t", "  ")
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
					// Debug(dedup)
					o = strings.Join(dedup, "\n")
					if errString != "" {
						o += "BTCJSONError:\n"
						o += errString
					}
					splitResult := strings.Split(o, "\n")
					const maxPerWidget = 16
					for i := 0; i < len(splitResult)-maxPerWidget; i += maxPerWidget {
						sri := strings.Join(splitResult[i:i+maxPerWidget], "\n")
						c.output = append(c.output,
							func(gtx l.Context) l.Dimensions {
								return wg.th.Caption(sri).
									Color("DocText").
									Font("go regular").MaxLines(maxPerWidget).
									Fn(gtx)
							})
					}
					return
				} else {
					if result, err = ctl.Call(wg.cx, method, params...); Check(err) {
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
								return wg.th.Flex().Rigid(
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
							return wg.th.Flex().Rigid(
								wg.th.Caption(sri).Color("DocText").Font("go regular").MaxLines(4).Fn,
							).Fn(gtx)
						})
				}
			}
			var ifc interface{}
			if err = json.Unmarshal(result, &ifc); Check(err) {
			}
			Debugs(ifc)
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

	return c
}

func (c *Console) Fn(gtx l.Context) l.Dimensions {
	tn := time.Now()
	defer func() {
		Debugf("console render time %d", time.Now().Sub(tn))
	}()
	le := func(gtx l.Context, index int) l.Dimensions {
		if index >= len(c.output) {
			return l.Dimensions{}
		} else {
			return c.output[index](gtx)
		}
	}
	// gtx.Constraints.Min = gtx.Constraints.Max
	fn := c.th.VFlex().
		Flexed(0.1,
			c.th.Fill("DocBg",
				// p9.EmptyMaxHeight(),
				// wg.th.H6("output area").Color("DocText").Fn,
				func(gtx l.Context) l.Dimensions {
					// gtx.Constraints.Max =
					// 	gtx.Constraints.Min
					return c.th.Inset(0.25,
						c.outputList.
							Vertical().
							// Background("DocBg").Color("DocText").Active("Primary").
							Length(len(c.output)).
							ListElement(le).
							Fn,
					).
						// Fn,
						// ).
						Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			c.th.Fill("PanelBg",
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

// CommandUsage display the usage for a specific command.
func CommandUsage(method string) (usage string) {
	var err error
	usage, err = btcjson.MethodUsageText(method)
	if err != nil {
		Error(err)
		// This should never happen since the method was already checked before calling this function, but be safe.
		usage = "Failed to obtain command usage: " + err.Error()
		return
	}
	return
}
