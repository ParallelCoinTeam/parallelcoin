package p9

import (
	l "gioui.org/layout"
	"gioui.org/op"
)

type DimensionList []l.Dimensions

func (d DimensionList) GetTotal(gtx l.Context, axis l.Axis) (total int) {
	for i := range d {
		total += axisMain(axis, d[i].Size)
	}
	_ = gtx.Metric
	return total
}

func (d DimensionList) PositionToCoordinate(position Position, axis l.Axis) (coordinate int) {
	for i := 0; i < position.First; i++ {
		coordinate += axisMain(axis, d[i].Size)
	}
	return coordinate + position.Offset
}

func (d DimensionList) CoordinateToPosition(coordinate int, axis l.Axis) (position Position) {
	cursor := 0
	for i := range d {
		cursor += axisMain(axis, d[i].Size)
		if cursor > coordinate {
			if i == 0 {
				position.First = 0
				position.Offset = coordinate - cursor
				position.BeforeEnd = true
				break
			}
			// step back
			cursor -= axisMain(axis, d[i].Size)
			position.First = i - 1
			position.Offset = coordinate - cursor
			position.BeforeEnd = true
			break
		}
	}
	return
}

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
