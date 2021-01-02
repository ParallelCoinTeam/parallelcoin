package gui

import (
	"image"
	"time"
	
	"gioui.org/f32"
	"gioui.org/gesture"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	l "gioui.org/layout"
	"gioui.org/op"
)

type clickEvents struct {
	Click, Cancel, Press func()
}

// Clickable represents a clickable area.
type Clickable struct {
	*Window
	click  gesture.Click
	clicks []click
	// prevClicks is the index into clicks that marks the clicks from the most recent Fn call. prevClicks is used to
	// keep clicks bounded.
	prevClicks int
	history    []press
	Events     clickEvents
}

func (w *Window) Clickable() (c *Clickable) {
	c = &Clickable{
		Window:     w,
		click:      gesture.Click{},
		clicks:     nil,
		prevClicks: 0,
		history:    nil,
		Events: clickEvents{
			Click: func() {
				Debug("click event")
			},
			Cancel: func() {
				Debug("cancel event")
			},
			Press: func() {
				Debug("press event")
			},
		},
	}
	return
}

func (c *Clickable) SetClick(fn func()) *Clickable {
	c.Events.Click = fn
	return c
}

func (c *Clickable) SetCancel(fn func()) *Clickable {
	c.Events.Cancel = fn
	return c
}

func (c *Clickable) SetPress(fn func()) *Clickable {
	c.Events.Press = fn
	return c
}

// click represents a click.
type click struct {
	Modifiers key.Modifiers
	NumClicks int
}

// press represents a past pointer press.
type press struct {
	// Position of the press.
	Position f32.Point
	// Start is when the press began.
	Start time.Time
	// End is when the press was ended by a release or Cancel. A zero End means it hasn't ended yet.
	End time.Time
	// Cancelled is true for cancelled presses.
	Cancelled bool
}

// Clicked reports whether there are pending clicks as would be reported by Clicks. If so, Clicked removes the earliest
// click.
func (c *Clickable) Clicked() bool {
	if len(c.clicks) == 0 {
		return false
	}
	n := copy(c.clicks, c.clicks[1:])
	c.clicks = c.clicks[:n]
	if c.prevClicks > 0 {
		c.prevClicks--
	}
	return true
}

// Clicks returns and clear the clicks since the last call to Clicks.
func (c *Clickable) Clicks() []click {
	clicks := c.clicks
	c.clicks = nil
	c.prevClicks = 0
	return clicks
}

// History is the past pointer presses useful for drawing markers. History is retained for a short duration (about a
// second).
func (c *Clickable) History() []press {
	return c.history
}

func (c *Clickable) Fn(gtx l.Context) l.Dimensions {
	c.update(gtx)
	stack := op.Push(gtx.Ops)
	pointer.Rect(image.Rectangle{Max: gtx.Constraints.Min}).Add(gtx.Ops)
	c.click.Add(gtx.Ops)
	stack.Pop()
	for len(c.history) > 0 {
		cc := c.history[0]
		if cc.End.IsZero() || gtx.Now.Sub(cc.End) < 1*time.Second {
			break
		}
		n := copy(c.history, c.history[1:])
		c.history = c.history[:n]
	}
	return l.Dimensions{Size: gtx.Constraints.Min}
}

// update the button changeState by processing clickEvents.
func (c *Clickable) update(gtx l.Context) {
	// Flush clicks from before the last update.
	n := copy(c.clicks, c.clicks[c.prevClicks:])
	c.clicks = c.clicks[:n]
	c.prevClicks = n
	for _, e := range c.click.Events(gtx) {
		switch e.Type {
		case gesture.TypeClick:
			click := click{
				Modifiers: e.Modifiers,
				NumClicks: e.NumClicks,
			}
			c.clicks = append(c.clicks, click)
			if ll := len(c.history); ll > 0 {
				c.history[ll-1].End = gtx.Now
			}
			c.Window.Runner <- func() error { c.Events.Click(); return nil }
		case gesture.TypeCancel:
			for i := range c.history {
				c.history[i].Cancelled = true
				if c.history[i].End.IsZero() {
					c.history[i].End = gtx.Now
				}
			}
			c.Window.Runner <- func() error { c.Events.Cancel(); return nil }
		case gesture.TypePress:
			c.history = append(c.history, press{
				Position: e.Position,
				Start:    gtx.Now,
			})
			c.
				Window.
				Runner <- func() error {
				c.
					Events.
					Press()
				return nil
			}
		}
	}
}
