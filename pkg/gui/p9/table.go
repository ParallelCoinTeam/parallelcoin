package p9

import (
	l "gioui.org/layout"
	"gioui.org/op"
)

type Cell struct {
	l.Widget
	dims l.Dimensions
}

func (c *Cell) getWidgetDimensions(gtx l.Context) {
	// gather the dimensions of the list elements
	gtx.Ops.Reset()
	child := op.Record(gtx.Ops)
	c.dims = c.Widget(gtx)
	_ = child.Stop()
	return
}

type CellRow struct {
	cells []Cell
}

type CellGrid struct {
	rows []CellRow
}

type Table struct {
	th     *Theme
	header CellRow
	body   CellGrid
}

func (th *Theme) Table() *Table {
	return &Table{
		th: th,
	}
}

func (t *Table) Header(h CellRow) *Table {
	t.header = h
	return t
}

func (t *Table) Body(g CellGrid) *Table {
	t.body = g
	return t
}

func (t *Table) Fn(gtx l.Context) l.Dimensions {
	for i := range t.body.rows {
		if len(t.header.cells) != len(t.body.rows[i].cells) {
			panic("not all rows are equal number of cells")
		}
	}
	gtx1 := GetInfContext(gtx)
	for i := range t.header.cells {
		t.header.cells[i].getWidgetDimensions(gtx1)
	}
	Debugs(t.header)
	for i := range t.body.rows {
		for j := range t.body.rows[i].cells {
			t.body.rows[i].cells[j].getWidgetDimensions(gtx1)
		}
	}
	Debugs(t.body)
	// find the max of each row and column

	return l.Dimensions{}
}
