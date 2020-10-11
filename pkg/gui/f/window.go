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

type window struct {
	Ctx   layout.Context
	w     *app.Window
	opts  []app.Option
	scale *scaledConfig
}

// Window creates a new window
func Window() (out *window) {
	var ops op.Ops
	var e system.FrameEvent
	out = &window{
		Ctx:   layout.NewContext(&ops, e),
		scale: &scaledConfig{1},
	}
	// out.set()
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

// Open sets the window options and initialise the app.window
func (w *window) Open() (out *window) {
	if w.scale == nil {
		w.Scale(1)
	}
	if w.opts != nil {
		w.w = app.NewWindow(w.opts...)
		w.opts = nil
	}
	return w
}

func (w *window) Run(frame func(ctx layout.Context), destroy func()) (err error) {
	// w.set()
	var ops op.Ops
	for {
		e := <-w.w.Events()
		switch e := e.(type) {
		case system.DestroyEvent:
			destroy()
			return e.Err
		case system.FrameEvent:
			ctx := layout.NewContext(&ops, e)
			frame(ctx)
			e.Frame(ctx.Ops)
		}
	}
}
