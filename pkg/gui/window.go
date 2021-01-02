package gui

import (
	"math"
	"time"
	
	"github.com/p9c/pod/pkg/gui/fonts/p9fonts"
	qu "github.com/p9c/pod/pkg/util/quit"
	
	"gioui.org/app"
	"gioui.org/io/system"
	l "gioui.org/layout"
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
	*Theme
	l.Context
	*app.Window
	opts   []app.Option
	scale  *scaledConfig
	Width  int // stores the width at the beginning of render
	Height int
	ops    op.Ops
	evQ    system.FrameEvent
}

// NewWindowP9 creates a new window
func NewWindowP9(quit chan struct{}) (out *Window) {
	out = &Window{
		Theme: NewTheme(p9fonts.Collection(), quit),
		scale: &scaledConfig{1},
	}
	out.WidgetPool = out.NewPool()
	out.Context = l.NewContext(&out.ops, out.evQ)
	return
}

// NewWindow creates a new window
func NewWindow(th *Theme) (out *Window) {
	out = &Window{
		Theme: th,
		scale: &scaledConfig{1},
	}
	out.Context = l.NewContext(&out.ops, out.evQ)
	return
}

// Title sets the title of the window
func (w *Window) Title(title string) (out *Window) {
	w.opts = append(w.opts, app.Title(title))
	return w
}

// Size sets the dimensions of the window
func (w *Window) Size(width, height float32) (out *Window) {
	w.opts = append(w.opts,
		app.Size(w.Theme.TextSize.Scale(width), w.Theme.TextSize.Scale(height)))
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

func (w *Window) Run(frame func(ctx l.Context) l.Dimensions,
	overlay func(ctx l.Context), destroy func(), quit qu.C) (err error) {
	var ops op.Ops
	for {
		select {
		case <-quit:
			return nil
		case e := <-w.Window.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				destroy()
				return e.Err
			case system.FrameEvent:
				ctx := l.NewContext(&ops, e)
				// update dimensions for responsive sizing widgets
				w.Width = ctx.Constraints.Max.X
				w.Height = ctx.Constraints.Max.Y
				frame(ctx)
				overlay(ctx)
				e.Frame(ctx.Ops)
			}
		}
	}
}
