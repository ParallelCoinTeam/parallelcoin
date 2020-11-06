package gui

import (
	"encoding/json"
	"fmt"
	l "gioui.org/layout"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
	"github.com/p9c/pod/pkg/rpc/legacy"
	"sort"
	"strings"
	"time"
)

type Console struct {
	Commands       []ConsoleCommand
	CommandsNumber int
}
type ConsoleCommand struct {
	Com      interface{}
	ComID    string
	Category string
	Out      string
	Time     time.Time
}

type ConsoleCommandsNumber struct {
	CommandsNumber int
}

//var (
//	consoleInputField = wg.th.Input("", "Amount", "Primary", "DocText", 25, func(pass string) {})
//	consoleOutputList = &layout.List{
//		Axis:        layout.Vertical,
//		ScrollToEnd: true,
//	}
//)

func (wg *WalletGUI) ConsolePage() l.Widget {
	return wg.th.VFlex().
		Flexed(1,
			wg.Inset(0.25,
				func(gtx l.Context) l.Dimensions {
					return wg.lists["console"].Vertical().Length(len(wg.sendAddresses)).ListElement(wg.consoleRow).Fn(gtx)
				},
			).Fn,
		).
		Rigid(
			wg.consoleInput(),
		).Fn
}

func (wg *WalletGUI) consoleRow(gtx l.Context, i int) l.Dimensions {
	t := wg.console.Commands[i]
	return wg.Flex().Vertical().AlignEnd().
		Rigid(
			wg.th.Caption("ds://" + t.ComID).
				Font("go regular").
				Color("PanelText").
				TextScale(0.66).Fn,
		).
		Rigid(
			wg.th.Caption(t.Out).
				Font("go regular").
				Color("PanelText").
				TextScale(0.66).Fn,
		).Fn(gtx)
}

func (wg *WalletGUI) consoleInput() l.Widget {
	return wg.inputs["console"].Input("", "Run command", "Secondary", "Primary", 25, func(pass string) {
		//func(e gel.SubmitEvent) {
		wg.console.Commands = append(
			wg.console.Commands,
			ConsoleCommand{
				ComID: wg.inputs["console"].GetText(),
				Time:  time.Time{},
				Out:   wg.ConsoleCmd(wg.inputs["console"].GetText()),
			})
	}).Fn
}
func (wg *WalletGUI) ConsoleCmd(com string) (o string) {
	split := strings.Split(com, " ")
	method := split[0]
	args := split[1:]
	var cmd, res interface{}
	var err error
	var errString, prev string
	if method == "help" {
		if len(args) < 1 {
			method = ""
			cmd = &btcjson.HelpCmd{Command: &method}
			if res, err = chainrpc.RPCHandlers["help"].Fn(wg.cx.RPCServer, cmd, nil); Check(err) {
				errString += fmt.Sprintln(err)
			}
			o += fmt.Sprintln(res)
			if res, err = legacy.RPCHandlers["help"].
				Handler(cmd, wg.cx.WalletServer, wg.cx.ChainClient); Check(err) {
				errString += fmt.Sprintln(err)
			}
			o += fmt.Sprintln(res)
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
				o += "Error:\n"
				o += errString
			}
		} else {
			method = args[0]
			//L.Debug("finding help for command", method)
			if help, err := wg.cx.RPCServer.HelpCacher.RPCMethodHelp(
				method); Check(err) {
				o += err.Error() + "\n"
				o += fmt.Sprintln(res)
				cmd = &btcjson.HelpCmd{Command: &method}
				if res, err = legacy.RPCHandlers["help"].
					Handler(cmd, wg.cx.WalletServer, wg.cx.ChainClient); Check(err) {
					errString += fmt.Sprintln(err)
				}
				o += fmt.Sprintln(res)
			} else {
				o += help
			}
			// if _, ok := legacy.RPCHandlers[method]; ok {
			// 	o += "wallet server:\n"
			// 	o += legacy.HelpDescsEnUS()[method]
			// }
			// if _, ok := rpc.RPCHandlers[method]; ok {
			// 	o += "chain server:\n"
			// 	o += rpc.HelpDescsEnUS[method]
			// }
		}
		return
	}
	params := make([]interface{}, 0, len(split[1:]))
	for _, arg := range args {
		params = append(params, arg)
	}
	if cmd, err = btcjson.NewCmd(method, params...); Check(err) {
		o += fmt.Sprintln(err)
	}
	if x, ok := chainrpc.RPCHandlers[method]; !ok {
		if x, ok := legacy.RPCHandlers[method]; ok {
			if res, err = x.Handler(cmd, wg.cx.WalletServer,
				wg.cx.ChainClient); Check(err) {
				o += err.Error()
			}
			// o += fmt.Sprintln(res)
		}
	} else {
		if res, err = x.Fn(wg.cx.RPCServer, cmd, nil); Check(err) {
			o += err.Error()
		}
		// o += fmt.Sprintln(res)
	}
	if res != nil {
		if j, err := json.MarshalIndent(res, "",
			"  "); !Check(err) {
			o += string(j)
		}
	}
	return

}
