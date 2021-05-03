package gel

import (
	"image"
	"image/color"

	"github.com/p9c/gio/f32"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op/clip"
	"github.com/p9c/gio/op/paint"
	"github.com/p9c/gio/unit"
	
	"github.com/p9c/gel/f32color"
)

type ProgressBar struct {
	*Window
	color    color.NRGBA
	progress int
}

// ProgressBar renders a horizontal bar with an indication of completion of a process
func (w *Window) ProgressBar() *ProgressBar {
	return &ProgressBar{
		Window:   w,
		progress: 0,
		color:    w.Colors.GetNRGBAFromName("Primary"),
	}
}

// SetProgress sets the progress of the progress bar
func (p *ProgressBar) SetProgress(progress int) *ProgressBar {
	p.progress = progress
	return p
}

// Color sets the color to render the bar in
func (p *ProgressBar) Color(c string) *ProgressBar {
	p.color = p.Theme.Colors.GetNRGBAFromName(c)
	return p
}

// Fn renders the progress bar as it is currently configured
func (p *ProgressBar) Fn(gtx l.Context) l.Dimensions {
	shader := func(width float32, color color.NRGBA) l.Dimensions {
		maxHeight := unit.Dp(4)
		rr := float32(gtx.Px(unit.Dp(2)))
		
		d := image.Point{X: int(width), Y: gtx.Px(maxHeight)}
		
		clip.RRect{
			Rect: f32.Rectangle{Max: f32.Point{X: width, Y: float32(gtx.Px(maxHeight))}},
			NE:   rr, NW: rr, SE: rr, SW: rr,
		}.Add(gtx.Ops)
		
		paint.ColorOp{Color: color}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)
		
		return l.Dimensions{Size: d}
	}
	
	progress := p.progress
	if progress > 100 {
		progress = 100
	} else if progress < 0 {
		progress = 0
	}
	
	progressBarWidth := float32(gtx.Constraints.Max.X)
	
	return l.Stack{Alignment: l.W}.Layout(gtx,
		l.Stacked(func(gtx l.Context) l.Dimensions {
			// Use a transparent equivalent of progress color.
			bgCol := f32color.MulAlpha(p.color, 150)
			
			return shader(progressBarWidth, bgCol)
		}),
		l.Stacked(func(gtx l.Context) l.Dimensions {
			fillWidth := (progressBarWidth / 100) * float32(progress)
			fillColor := p.color
			if gtx.Queue == nil {
				fillColor = f32color.MulAlpha(fillColor, 200)
			}
			return shader(fillWidth, fillColor)
		}),
	)
}
