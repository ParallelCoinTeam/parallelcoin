package gui

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	l "gioui.org/layout"
	icons2 "golang.org/x/exp/shiny/materialdesign/icons"

	"github.com/atotto/clipboard"

	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
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
			wg.th.Body1(txt).Color("DocText").Fn)
		c.editor.SetText("")
		split := strings.Split(txt, " ")
		method, args := split[0], split[1:]
		var params []interface{}
		for i := range args {
			params = append(params, args[i])
		}
		Debug("method", method, "args", args)
		if wg.WalletClient != nil {
			// Ensure the specified method identifies a valid registered command and is one of the usable types.
			usageFlags, err := btcjson.MethodUsageFlags(method)
			if err != nil {
				Error(err)
				fmt.Fprintf(os.Stderr, "Unrecognized command '%s'\n", method)
				// HelpPrint()
				// os.Exit(1)
			}
			if usageFlags&unusableFlags != 0 {
				fmt.Fprintf(os.Stderr, "The '%s' command can only be used via websockets\n", method)
				// HelpPrint()
				// os.Exit(1)
			}
			// // Attempt to create the appropriate command using the arguments provided by the user.
			// cmd, err := btcjson.NewCmd(method, params...)
			// if err != nil {
			// 	Error(err)
			// 	// Show the error along with its error code when it's a json. BTCJSONError as it realistically will always be
			// 	// since the NewCmd function is only supposed to return errors of that type.
			// 	if jerr, ok := err.(btcjson.BTCJSONError); ok {
			// 		fmt.Fprintf(os.Stderr, "%s command: %v (code: %s)\n",
			// 			method, err, jerr.ErrorCode)
			// 		// CommandUsage(method)
			// 		// os.Exit(1)
			// 	}
			// 	// The error is not a json.BTCJSONError and this really should not happen. Nevertheless fall back to just
			// 	// showing the error if it should happen due to a bug in the package.
			// 	fmt.Fprintf(os.Stderr, "%s command: %v\n", method, err)
			// CommandUsage(method)
			// os.Exit(1)
			// }
			// Marshal the command into a JSON-RPC byte slice in preparation for sending it to the RPC server.
			// marshalledJSON, err = btcjson.MarshalCmd(1, cmd)

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
