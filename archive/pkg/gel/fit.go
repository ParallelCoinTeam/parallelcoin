// SPDX-License-Identifier: Unlicense OR MIT

package gel

import (
	"image"
	
	"github.com/p9c/gio/f32"
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
	"github.com/p9c/gio/op/clip"
)

// Fit scales a widget to fit and clip to the constraints.
type Fit uint8

const (
	// Unscaled does not alter the scale of a widget.
	Unscaled Fit = iota
	// Contain scales widget as large as possible without cropping
	// and it preserves aspect-ratio.
	Contain
	// Cover scales the widget to cover the constraint area and
	// preserves aspect-ratio.
	Cover
	// ScaleDown scales the widget smaller without cropping,
	// when it exceeds the constraint area.
	// It preserves aspect-ratio.
	ScaleDown
	// Stretch stretches the widget to the constraints and does not
	// preserve aspect-ratio.
	Stretch
)

// scale adds clip and scale operations to fit dims to the constraints.
// It positions the widget to the appropriate position.
// It returns dimensions modified accordingly.
func (fit Fit) scale(gtx l.Context, pos l.Direction, dims l.Dimensions) l.Dimensions {
	widgetSize := dims.Size
	
	if fit == Unscaled || dims.Size.X == 0 || dims.Size.Y == 0 {
		dims.Size = gtx.Constraints.Constrain(dims.Size)
		clip.Rect{Max: dims.Size}.Add(gtx.Ops)
		
		offset := pos.Position(widgetSize, dims.Size)
		op.Offset(l.FPt(offset)).Add(gtx.Ops)
		dims.Baseline += offset.Y
		return dims
	}
	
	scale := f32.Point{
		X: float32(gtx.Constraints.Max.X) / float32(dims.Size.X),
		Y: float32(gtx.Constraints.Max.Y) / float32(dims.Size.Y),
	}
	
	switch fit {
	case Contain:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case Cover:
		if scale.Y > scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
	case ScaleDown:
		if scale.Y < scale.X {
			scale.X = scale.Y
		} else {
			scale.Y = scale.X
		}
		
		// The widget would need to be scaled up, no change needed.
		if scale.X >= 1 {
			dims.Size = gtx.Constraints.Constrain(dims.Size)
			clip.Rect{Max: dims.Size}.Add(gtx.Ops)
			
			offset := pos.Position(widgetSize, dims.Size)
			op.Offset(l.FPt(offset)).Add(gtx.Ops)
			dims.Baseline += offset.Y
			return dims
		}
	case Stretch:
	}
	
	var scaledSize image.Point
	scaledSize.X = int(float32(widgetSize.X) * scale.X)
	scaledSize.Y = int(float32(widgetSize.Y) * scale.Y)
	dims.Size = gtx.Constraints.Constrain(scaledSize)
	dims.Baseline = int(float32(dims.Baseline) * scale.Y)
	
	clip.Rect{Max: dims.Size}.Add(gtx.Ops)
	
	offset := pos.Position(scaledSize, dims.Size)
	op.Affine(f32.Affine2D{}.
		Scale(f32.Point{}, scale).
		Offset(l.FPt(offset)),
	).Add(gtx.Ops)
	
	dims.Baseline += offset.Y
	
	return dims
}
