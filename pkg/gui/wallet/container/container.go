package container

import (
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/gioapp/gel/helper"
)

type container struct {
	outsideColor, borderColor, insideColor string
	margin, border, padding                int
}

func C() *container {
	return &container{
		outsideColor: "",
		borderColor:  "",
		insideColor:  "",
		margin:       0,
		border:       0,
		padding:      0,
	}
}

func (c *container) OutsideColor(color string) *container {
	c.outsideColor = color
	return c
}
func (c *container) BorderColor(color string) *container {
	c.borderColor = color
	return c
}
func (c *container) InsideColor(color string) *container {
	c.insideColor = color
	return c
}

func (c *container) Margin(size int) *container {
	c.margin = size
	return c
}
func (c *container) Border(size int) *container {
	c.border = size
	return c
}
func (c *container) Padding(size int) *container {
	c.padding = size
	return c
}

func (c *container) Layout(content layout.Widget) layout.Widget {
	return stack(c.outsideColor, c.margin, stack(c.borderColor, c.border, stack(c.insideColor, c.padding, content)))
}

func stack(color string, size int, content layout.Widget) layout.Widget {
	return func(gtx layout.Context) layout.Dimensions {
		return layout.Stack{Alignment: layout.W}.Layout(gtx,
			layout.Expanded(func(gtx layout.Context) layout.Dimensions {
				return helper.Fill(gtx, helper.HexARGB(color))
			}),
			layout.Stacked(func(gtx layout.Context) layout.Dimensions {
				return layout.Inset{
					Top:    unit.Dp(float32(size)),
					Right:  unit.Dp(float32(size)),
					Bottom: unit.Dp(float32(size)),
					Left:   unit.Dp(float32(size)),
				}.Layout(gtx, content)
			}),
		)
	}
}
