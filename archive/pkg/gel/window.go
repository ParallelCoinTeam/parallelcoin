package gel

import (
	"math"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/p9c/opts/binary"
	"github.com/p9c/opts/meta"

	clipboard2 "github.com/p9c/gio/io/clipboard"

	"github.com/p9c/gel/clipboard"
	"github.com/p9c/gel/fonts/p9fonts"

	"github.com/p9c/gio/io/event"

	"github.com/p9c/qu"

	uberatomic "go.uber.org/atomic"

	"github.com/p9c/gio/app"
	"github.com/p9c/gio/io/system"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/unit"
)

type CallbackQueue chan func() error

func NewCallbackQueue(bufSize int) CallbackQueue {
	return make(CallbackQueue, bufSize)
}

type scaledConfig struct {
	Scale float32
}

func (s *scaledConfig) Now() time.Time {
	return time.Now()
}

func (s *scaledConfig) Px(v unit.Value) int {
	scale := s.Scale
	if v.U == unit.UnitPx {
		scale = 1
	}
	return int(math.Round(float64(scale * v.V)))
}

type Window struct {
	*Theme
	*app.Window
	opts                  []app.Option
	scale                 *scaledConfig
	Width                 *uberatomic.Int32 // stores the width at the beginning of render
	Height                *uberatomic.Int32
	ops                   op.Ops
	evQ                   system.FrameEvent
	Runner                CallbackQueue
	overlay               []*func(gtx l.Context)
	ClipboardWriteReqs    chan string
	ClipboardReadReqs     chan func(string)
	ClipboardContent      string
	clipboardReadReady    qu.C
	clipboardReadResponse chan string
}

func (w *Window) PushOverlay(overlay *func(gtx l.Context)) {
	w.overlay = append(w.overlay, overlay)
}

func (w *Window) PopOverlay(overlay *func(gtx l.Context)) {
	if len(w.overlay) == 0 {
		return
	}
	index := -1
	for i := range w.overlay {
		if overlay == w.overlay[i] {
			index = i
			break
		}
	}
	if index != -1 {
		if index == len(w.overlay)-1 {
			w.overlay = w.overlay[:index]
		} else if index == 0 {
			w.overlay = w.overlay[1:]
		} else {
			w.overlay = append(w.overlay[:index], w.overlay[index+1:]...)
		}
	}
}

func (w *Window) Overlay(gtx l.Context) {
	for _, overlay := range w.overlay {
		(*overlay)(gtx)
	}
}

// NewWindowP9 creates a new window
func NewWindowP9(quit chan struct{}) (out *Window) {
	out = &Window{
		scale:              &scaledConfig{1},
		Runner:             NewCallbackQueue(32),
		Width:              uberatomic.NewInt32(0),
		Height:             uberatomic.NewInt32(0),
		ClipboardWriteReqs: make(chan string, 1),
		ClipboardReadReqs:  make(chan func(string), 32),
		clipboardReadReady: qu.Ts(1),
		clipboardReadResponse: make(chan string,1),
	}
	out.Theme = NewTheme(
		binary.New(meta.Data{}, false, nil),
		p9fonts.Collection(), quit,
	)
	out.Theme.WidgetPool = out.NewPool()
	clipboard.Start()
	return
}

// NewWindow creates a new window
func NewWindow(th *Theme) (out *Window) {
	out = &Window{
		Theme: th,
		scale: &scaledConfig{1},
	}
	return
}

// Title sets the title of the window
func (w *Window) Title(title string) (out *Window) {
	w.opts = append(w.opts, app.Title(title))
	return w
}

// Size sets the dimensions of the window
func (w *Window) Size(width, height float32) (out *Window) {
	w.opts = append(
		w.opts,
		app.Size(w.TextSize.Scale(width), w.TextSize.Scale(height)),
	)
	return w
}

// Scale sets the scale factor for rendering
func (w *Window) Scale(s float32) *Window {
	w.scale = &scaledConfig{s}
	return w
}

// Open sets the window options and initialise the node.window
func (w *Window) Open() (out *Window) {
	if w.scale == nil {
		w.Scale(1)
	}
	if w.opts != nil {
		w.Window = app.NewWindow(w.opts...)
		w.opts = nil
	}

	return w
}

func (w *Window) Run(frame func(ctx l.Context) l.Dimensions, destroy func(), quit qu.C,) (e error) {
	runner := func() {
		ticker := time.NewTicker(time.Second)
		for {
			select {
			case content := <-w.ClipboardWriteReqs:
				w.WriteClipboard(content)
			case fn := <-w.ClipboardReadReqs:
				go func() {
					w.ReadClipboard()
					fn(<-w.clipboardReadResponse)
				}()
			case <-ticker.C:
				if runtime.GOOS == "linux" {
					var e error
					var b []byte
					textSize := unit.Sp(16)
					runner := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "text-scaling-factor")
					if b, e = runner.CombinedOutput(); D.Chk(e) {
					}
					var factor float64
					numberString := strings.TrimSpace(string(b))
					if factor, e = strconv.ParseFloat(numberString, 10); D.Chk(e) {
					}
					w.TextSize = textSize.Scale(float32(factor))
					// I.Ln(w.TextSize)
				}
				w.Invalidate()
			case fn := <-w.Runner:
				if e = fn(); E.Chk(e) {
					return
				}
			case <-quit.Wait():
				return
			// by repeating selectors we decrease the chance of a runner delaying
			// a frame event hitting the physical frame deadline
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			case ev := <-w.Window.Events():
				if e = w.processEvents(ev, frame, destroy); E.Chk(e) {
					return
				}
			}
		}
	}
	switch runtime.GOOS {
	case "ios", "android":
		go runner()
	default:
		runner()
	}
	return
}

func (w *Window) processEvents(e event.Event, frame func(ctx l.Context) l.Dimensions, destroy func()) error {
	switch ev := e.(type) {
	case system.DestroyEvent:
		D.Ln("received destroy event", ev.Err)
		destroy()
		return ev.Err
	case system.FrameEvent:
		ops := op.Ops{}
		c := l.NewContext(&ops, ev)
		// update dimensions for responsive sizing widgets
		w.Width.Store(int32(c.Constraints.Max.X))
		w.Height.Store(int32(c.Constraints.Max.Y))
		frame(c)
		w.Overlay(c)
		ev.Frame(c.Ops)
	case clipboard2.Event:
		// w.ClipboardContent = ev.Text
		// w.clipboardReadReady.Signal()
		w.clipboardReadResponse <- ev.Text
	}
	return nil
}
