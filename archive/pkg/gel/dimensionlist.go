package gel

import (
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
)

type DimensionList []l.Dimensions

func (d DimensionList) GetTotal(axis l.Axis) (total int) {
	for i := range d {
		total += axisMain(axis, d[i].Size)
	}
	return total
}

// PositionToCoordinate converts a list position to absolute coordinate
func (d DimensionList) PositionToCoordinate(position Position, axis l.Axis) (coordinate int) {
	for i := 0; i < position.First; i++ {
		coordinate += axisMain(axis, d[i].Size)
	}
	return coordinate + position.Offset
}

// CoordinateToPosition converts an absolute coordinate to a list position
func (d DimensionList) CoordinateToPosition(coordinate int, axis l.Axis) (position Position) {
	cursor := 0
	if coordinate < 0 {
		coordinate = 0
		return
	}
	tot := d.GetTotal(axis)
	if coordinate > tot {
		position.First = len(d) - 1
		position.Offset = axisMain(axis, d[len(d)-1].Size)
		position.BeforeEnd = false
		return
	}
	var i int
	for i = range d {
		cursor += axisMain(axis, d[i].Size)
		if cursor >= coordinate {
			position.First = i
			position.Offset = coordinate - cursor
			position.BeforeEnd = true
			return
		}
	}
	// if it overshoots, stop it, if it is at the end, mark it
	if coordinate >= cursor {
		position.First = len(d) - 1
		position.Offset = axisMain(axis, d[len(d)-1].Size)
		position.BeforeEnd = false
	}
	return
}

// GetDimensionList returns a dimensionlist based on the given listelement
func GetDimensionList(gtx l.Context, length int, listElement ListElement) (dims DimensionList) {
	// gather the dimensions of the list elements
	for i := 0; i < length; i++ {
		child := op.Record(gtx.Ops)
		d := listElement(gtx, i)
		_ = child.Stop()
		dims = append(dims, d)
	}
	return
}

func GetDimension(gtx l.Context, w l.Widget) (dim l.Dimensions) {
	child := op.Record(gtx.Ops)
	dim = w(gtx)
	_ = child.Stop()
	return
}

func (d DimensionList) GetSizes(position Position, axis l.Axis) (total, before int) {
	for i := range d {
		inc := axisMain(axis, d[i].Size)
		total += inc
		if i < position.First {
			before += inc
		}
	}
	before += position.Offset
	return
}
