package widget

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

type ClickEvents struct {
	Click, Cancel, Press func()
}

// Clickable represents a clickable area.
type Clickable struct {
	click  gesture.Click
	clicks []Click
	// prevClicks is the index into clicks that marks the clicks from the most recent Fn call. prevClicks is used to
	// keep clicks bounded.
	prevClicks int
	history    []Press
	Events     ClickEvents
}

func NewClickable() (c *Clickable) {
	c = &Clickable{
		click:      gesture.Click{},
		clicks:     nil,
		prevClicks: 0,
		history:    nil,
		Events: ClickEvents{
			Click: func() {
				// Debug("click event")
			},
			Cancel: func() {
				// Debug("cancel event")
			},
			Press: func() {
				// Debug("press event")
			},
		},
	}
	return
}

func (c *Clickable) SetClick(fn func()) *Clickable {
	if c.Events.Click == nil {
		c.Events.Click = func() {}
	}
	c.Events.Click = fn
	return c
}

func (c *Clickable) SetCancel(fn func()) *Clickable {
	if c.Events.Cancel == nil {
		c.Events.Cancel = func() {}
	}
	c.Events.Cancel = fn
	return c
}

func (c *Clickable) SetPress(fn func()) *Clickable {
	if c.Events.Press == nil {
		c.Events.Press = func() {}
	}
	c.Events.Press = fn

	return c
}

// Click represents a Click.
type Click struct {
	Modifiers key.Modifiers
	NumClicks int
}

// Press represents a past pointer Press.
type Press struct {
	// Position of the Press.
	Position f32.Point
	// Start is when the Press began.
	Start time.Time
	// End is when the Press was ended by a release or Cancel. A zero End means it hasn't ended yet.
	End time.Time
	// Cancelled is true for cancelled presses.
	Cancelled bool
}

// Clicked reports whether there are pending clicks as would be reported by Clicks. If so, Clicked removes the earliest
// Click.
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
func (c *Clickable) Clicks() []Click {
	clicks := c.clicks
	c.clicks = nil
	c.prevClicks = 0
	return clicks
}

// History is the past pointer presses useful for drawing markers. History is retained for a short duration (about a
// second).
func (c *Clickable) History() []Press {
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

// update the button changeState by processing ClickEvents.
func (c *Clickable) update(gtx l.Context) {
	// if this is used by old code these functions have to be empty as they are called, not nil (which will panic)
	if c.Events.Click == nil {
		c.Events.Click = func() {
			Debug("click event")
		}
	}
	if c.Events.Cancel == nil {
		c.Events.Cancel = func() {
			Debug("cancel event")
		}
	}
	if c.Events.Press == nil {
		c.Events.Press = func() {
			Debug("press event")
		}
	}
	// Flush clicks from before the last update.
	n := copy(c.clicks, c.clicks[c.prevClicks:])
	c.clicks = c.clicks[:n]
	c.prevClicks = n
	for _, e := range c.click.Events(gtx) {
		switch e.Type {
		case gesture.TypeClick:
			click := Click{
				Modifiers: e.Modifiers,
				NumClicks: e.NumClicks,
			}
			c.clicks = append(c.clicks, click)
			if l := len(c.history); l > 0 {
				c.history[l-1].End = gtx.Now
			}
			c.Events.Click()
		case gesture.TypeCancel:
			for i := range c.history {
				c.history[i].Cancelled = true
				if c.history[i].End.IsZero() {
					c.history[i].End = gtx.Now
				}
			}
			c.Events.Cancel()
		case gesture.TypePress:
			c.history = append(c.history, Press{
				Position: e.Position,
				Start:    gtx.Now,
			})
			c.Events.Press()
		}
	}
}
