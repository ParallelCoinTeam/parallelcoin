package f

import (
	"github.com/p9c/pod/pkg/util/interrupt"
	qu "github.com/p9c/pod/pkg/util/quit"
	"math"
	"time"
	
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/unit"
	"github.com/p9c/pod/pkg/gui/p9"
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
	Theme  *p9.Theme
	Window *app.Window
	opts   []app.Option
	scale  *scaledConfig
	Width  *int // stores the width at the beginning of render
	Height *int
}

// NewWindow creates a new window
func NewWindow(th *p9.Theme) (out *Window) {
	var ops op.Ops
	var e system.FrameEvent
	var width, height int
	out = &Window{
		Ctx:    layout.NewContext(&ops, e),
		Theme:  th,
		scale:  &scaledConfig{1},
		Width:  &width,
		Height: &height,
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
func (w *Window) Size(width, height float32) (out *Window) {
	w.opts = append(
		w.opts,
		app.Size(w.Theme.TextSize.Scale(width), w.Theme.TextSize.Scale(height)),
	)
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

func (w *Window) Run(
	frame func(ctx layout.Context) layout.Dimensions,
	overlay func(ctx layout.Context), destroy func(), quit qu.C,
) (err error) {
	var lastFrameTime time.Time
	_ = lastFrameTime
	frameTicker := time.NewTicker(time.Second / 60)
	var ops op.Ops
	var recently bool
	go func() {
		for {
			select {
			case <-quit:
				return
			case <-frameTicker.C:
				// if the last time a frame event was received is more than what should have been 5 frames, print
				// goroutine dump
				timeSinceLastFrame := time.Now().Sub(lastFrameTime)
				if timeSinceLastFrame > time.Second/5 && recently {
					Debug(interrupt.GoroutineDump())
					recently = false
				}
			}
		}
	}()
	for {
		select {
		case <-quit:
			return nil
		case e := <-w.Window.Events():
			switch e := e.(type) {
			case system.StageEvent:
				Debug("StageEvent", e.Stage)
			case system.ClipboardEvent:
				Debug("ClipboardEvent", e)
			case system.DestroyEvent:
				destroy()
				w.Window.Close()
				return e.Err
			case system.FrameEvent:
				lastFrameTime = time.Now()
				recently = true
				gtx := layout.NewContext(&ops, e)
				// update dimensions for responsive sizing widgets
				*w.Width = gtx.Constraints.Max.X
				*w.Height = gtx.Constraints.Max.Y
				frame(gtx)
				overlay(gtx)
				e.Frame(gtx.Ops)
				
			}
		}
	}
}
