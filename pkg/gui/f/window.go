package f

import (
	"math"
	"time"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
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

type Window struct {
	Ctx    layout.Context
	Window *app.Window
	opts   []app.Option
	scale  *scaledConfig
}

// NewWindow creates a new window
func NewWindow() (out *Window) {
	var ops op.Ops
	var e system.FrameEvent
	out = &Window{
		Ctx:   layout.NewContext(&ops, e),
		scale: &scaledConfig{1},
	}
	// out.set()
	return
}

// Title sets the title of the window
func (w *Window) Title(title string) (out *Window) {
	w.opts = append(w.opts, app.Title(title))
	return w
}

// Size sets the dimensions of the window
func (w *Window) Size(width, height int) (out *Window) {
	w.opts = append(w.opts,
		app.Size(unit.Sp(float32(width)), unit.Sp(float32(height))))
	return w
}

// Scale sets the scale factor for rendering
func (w *Window) Scale(s float32) *Window {
	w.scale = &scaledConfig{s}
	return w
}

// Open sets the window options and initialise the app.window
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

func (w *Window) Run(frame func(ctx layout.Context) layout.Dimensions, destroy func(), quit chan struct{}) (err error) {
	var ops op.Ops
	for {
		select {
		case <-quit:
			return nil
		case e := <-w.Window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				go destroy()
				return e.Err
			case system.FrameEvent:
				ctx := layout.NewContext(&ops, e)
				frame(ctx)
				e.Frame(ctx.Ops)
			}
		}
	}
}
