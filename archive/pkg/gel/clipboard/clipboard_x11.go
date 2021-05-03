// +build linux freebsd openbsd

package clipboard

import (
	"fmt"
	"os"
	"time"
	
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
)

// todo: only X is required from this package, the rest runs off the built-in Gio clipboard

const debugClipboardRequests = false

var X *xgb.Conn
var win xproto.Window
var clipboardText string
var selnotify chan bool

var clipboardAtom, primaryAtom, textAtom, targetsAtom, atomAtom xproto.Atom
var targetAtoms []xproto.Atom
var clipboardAtomCache = map[xproto.Atom]string{}

var RunningX bool

// Start up the clipboard watcher
func Start() {
	// first, check if X is running as only X has the Primary buffer
	env := os.Environ()
	for i := range env {
		if env[i]=="XDG_SESSION_TYPE=x11" {
			I.Ln("running X11, primary selection buffer enabled")
			RunningX=true
			break
		}
	}
	if !RunningX {
		return
	}
	var e error
	X, e = xgb.NewConnDisplay("")
	if e != nil {
		panic(e)
	}

	selnotify = make(chan bool, 1)

	win, e = xproto.NewWindowId(X)
	if e != nil {
		panic(e)
	}

	setup := xproto.Setup(X)
	s := setup.DefaultScreen(X)
	e = xproto.CreateWindowChecked(
		X,
		s.RootDepth,
		win,
		s.Root,
		100,
		100,
		1,
		1,
		0,
		xproto.WindowClassInputOutput,
		s.RootVisual,
		0,
		[]uint32{},
	).Check()
	if e != nil {
		panic(e)
	}

	clipboardAtom = internAtom(X, "CLIPBOARD")
	primaryAtom = internAtom(X, "PRIMARY")
	textAtom = internAtom(X, "UTF8_STRING")
	targetsAtom = internAtom(X, "TARGETS")
	atomAtom = internAtom(X, "ATOM")

	targetAtoms = []xproto.Atom{targetsAtom, textAtom}

	go eventLoop()
}

func Set(text string) (e error){
	if !RunningX {
		return
	}
	clipboardText = text
	ssoc := xproto.SetSelectionOwnerChecked(X, win, clipboardAtom, xproto.TimeCurrentTime)
	if e = ssoc.Check(); E.Chk(e) {
		fmt.Fprintf(os.Stderr, "Error setting clipboard: %v", e)
	}
	ssoc = xproto.SetSelectionOwnerChecked(X, win, primaryAtom, xproto.TimeCurrentTime)
	if e = ssoc.Check(); E.Chk(e) {
		fmt.Fprintf(os.Stderr, "Error setting primary selection: %v", e)
	}
	return
}

func SetPrimary(text string) (e error){
	if !RunningX {
		return
	}
	clipboardText = text
	ssoc := xproto.SetSelectionOwnerChecked(X, win, primaryAtom, xproto.TimeCurrentTime)
	if e = ssoc.Check(); E.Chk(e) {
		fmt.Fprintf(os.Stderr, "Error setting primary selection: %v", e)
	}
	return
}

func Get() string {
	if !RunningX {
		return ""
	}
	return getSelection(clipboardAtom)
}

func GetPrimary() string {
	if !RunningX {
		return ""
	}
	return getSelection(primaryAtom)
}

func getSelection(selAtom xproto.Atom) string {
	csc := xproto.ConvertSelectionChecked(X, win, selAtom, textAtom, selAtom, xproto.TimeCurrentTime)
	e := csc.Check()
	if e != nil {
		fmt.Fprintln(os.Stderr, e)
		return ""
	}
	
	select {
	case r := <-selnotify:
		if !r {
			return ""
		}
		gpc := xproto.GetProperty(X, true, win, selAtom, textAtom, 0, 5*1024*1024)
		gpr, e := gpc.Reply()
		if e != nil {
			fmt.Fprintln(os.Stderr, e)
			return ""
		}
		if gpr.BytesAfter != 0 {
			fmt.Fprintln(os.Stderr, "Clipboard too large")
			return ""
		}
		return string(gpr.Value[:gpr.ValueLen])
	case <-time.After(1 * time.Second):
		fmt.Fprintln(os.Stderr, "Clipboard retrieval failed, timeout")
		return ""
	}
}

func eventLoop() {
	for {
		ev, e := X.WaitForEvent()
		if e != nil {
			continue
		}
		switch evt := ev.(type) {
		case xproto.SelectionRequestEvent:
			if debugClipboardRequests {
				tgtname := lookupAtom(evt.Target)
				fmt.Fprintln(
					os.Stderr,
					"SelectionRequest",
					ev,
					textAtom,
					tgtname,
					"isPrimary:",
					evt.Selection == primaryAtom,
					"isClipboard:",
					evt.Selection == clipboardAtom,
				)
			}
			t := clipboardText
			
			switch evt.Target {
			case textAtom:
				if debugClipboardRequests {
					fmt.Fprintln(os.Stderr, "Sending as text")
				}
				cpc := xproto.ChangePropertyChecked(
					X,
					xproto.PropModeReplace,
					evt.Requestor,
					evt.Property,
					textAtom,
					8,
					uint32(len(t)),
					[]byte(t),
				)
				e := cpc.Check()
				if e == nil {
					sendSelectionNotify(evt)
				} else {
					fmt.Fprintln(os.Stderr, e)
				}
			
			case targetsAtom:
				if debugClipboardRequests {
					fmt.Fprintln(os.Stderr, "Sending targets")
				}
				buf := make([]byte, len(targetAtoms)*4)
				for i, atom := range targetAtoms {
					xgb.Put32(buf[i*4:], uint32(atom))
				}
				
				_ = xproto.ChangePropertyChecked(
					X,
					xproto.PropModeReplace,
					evt.Requestor,
					evt.Property,
					atomAtom,
					32,
					uint32(len(targetAtoms)),
					buf,
				).Check()
				if e == nil {
					sendSelectionNotify(evt)
				} else {
					fmt.Fprintln(os.Stderr, e)
				}
			
			default:
				if debugClipboardRequests {
					fmt.Fprintln(os.Stderr, "Skipping")
				}
				evt.Property = 0
				sendSelectionNotify(evt)
			}
		
		case xproto.SelectionNotifyEvent:
			selnotify <- (evt.Property == clipboardAtom) || (evt.Property == primaryAtom)
		}
	}
}

func lookupAtom(at xproto.Atom) string {
	if s, ok := clipboardAtomCache[at]; ok {
		return s
	}
	
	reply, e := xproto.GetAtomName(X, at).Reply()
	if e != nil {
		panic(e)
	}
	
	// If we're here, it means we didn't have the ATOM id cached. So cache it.
	atomName := string(reply.Name)
	clipboardAtomCache[at] = atomName
	return atomName
}

func sendSelectionNotify(ev xproto.SelectionRequestEvent) {
	sn := xproto.SelectionNotifyEvent{
		Time:      xproto.TimeCurrentTime,
		Requestor: ev.Requestor,
		Selection: ev.Selection,
		Target:    ev.Target,
		Property:  ev.Property,
	}
	var e error
	sec := xproto.SendEventChecked(X, false, ev.Requestor, 0, string(sn.Bytes()))
	if e = sec.Check(); E.Chk(e) {
		fmt.Fprintln(os.Stderr, e)
	}
}

func internAtom(conn *xgb.Conn, n string) xproto.Atom {
	iac := xproto.InternAtom(conn, true, uint16(len(n)), n)
	iar, e := iac.Reply()
	if e != nil {
		panic(e)
	}
	return iar.Atom
}
