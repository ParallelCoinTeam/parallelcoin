package fui

import (
	"gioui.org/app"
	"gioui.org/io/event"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/unit"
	"math"
	"time"
)

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

type window struct {
	w     *app.Window
	opts  []app.Option
	scale *scaledConfig
}

// Window creates a new window
func Window() (out *window) {
	out = &window{}
	return
}

// Title sets the title of the window
func (w *window) Title(title string) (out *window) {
	w.opts = append(w.opts, app.Title(title))
	return w
}

// Size sets the dimensions of the window
func (w *window) Size(width, height int) (out *window) {
	w.opts = append(w.opts,
		app.Size(unit.Dp(float32(width)), unit.Dp(float32(height))))
	return w
}

// Scale sets the scale factor for rendering
func (w *window) Scale(s float32) *window {
	w.scale = &scaledConfig{s}
	return w
}

// Set the window options and initialise the app.window
func (w *window) set() (out *window) {
	if w.opts != nil {
		w.w = app.NewWindow(w.opts...)
		w.opts = nil
	}
	if w.scale == nil {
		w.Scale(1)
	}
	return w
}

// These override app.window methods to ensure the options are set first

// Queue returns the window's event queue
func (w *window) Queue() (out *app.Queue) {
	w.set()
	return w.w.Queue()
}

// Events returns the channel for events registered with the window
func (w *window) Events() (out <-chan event.Event) {
	w.set()
	return w.w.Events()
}

// Context for the window
func (w *window) Context() (out *layout.Context) {
	w.set()
	return layout.NewContext(w.w.Queue())
}

func (w *window) Run(frame func(ctx *layout.Context), destroy func()) {
	w.set()
	ctx := w.Context()
	go func() {
		for {
			select {
			case e := <-w.Events():
				switch e := e.(type) {
				case system.DestroyEvent:
					destroy()
				case system.FrameEvent:
					ctx.Reset(w.scale, e.Size)
					frame(ctx)
					e.Frame(ctx.Ops)
				}
			}
		}
	}()
}
