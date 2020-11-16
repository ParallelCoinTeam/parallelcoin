package p9

import (
	"sort"

	l "gioui.org/layout"
	"gioui.org/op"
)

type Cell struct {
	l.Widget
	dims l.Dimensions
	// priority only has meaning for the header row in defining an order of eliminating elements to fit a width.
	// When trimming size to fit width add from highest to lowest priority and stop when dimensions exceed the target.
	Priority int
}

func (c *Cell) getWidgetDimensions(gtx l.Context) {
	// gather the dimensions of the list elements
	gtx.Ops.Reset()
	child := op.Record(gtx.Ops)
	c.dims = c.Widget(gtx)
	_ = child.Stop()
	return
}

type CellRow []Cell

func (c CellRow) GetPriority() (out CellPriority) {
	for i := range c {
		out = append(out, c[i].Priority)
	}
	sort.Sort(out)
	return
}

type CellPriority []int

// sort a cell row by priority
func (c CellPriority) Len() int {
	return len(c)
}
func (c CellPriority) Less(i, j int) bool {
	return c[i] < c[j]
}
func (c CellPriority) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CellGrid []CellRow

// Table is a super simple table widget that finds the dimensions of all cells, sets all to max of each axis, and then
// scales the remaining space evenly after pruning off columns that cause the table to exceed the maximum width by
// adding columns with the highest priority (lowest value) first.
type Table struct {
	th     *Theme
	header CellRow
	body   CellGrid
	Y, X   []int
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
	for i := range t.body {
		if len(t.header) != len(t.body[i]) {
			// this should never happen hence panic
			panic("not all rows are equal number of cells")
		}
	}
	gtx1 := GetInfContext(gtx)
	// gather the dimensions from all cells
	for i := range t.header {
		t.header[i].getWidgetDimensions(gtx1)
	}
	Debugs(t.header)
	for i := range t.body {
		for j := range t.body[i] {
			t.body[i][j].getWidgetDimensions(gtx1)
		}
	}
	Debugs(t.body)
	// find the max of each row and column
	var table []CellRow
	table = append(table, t.header)
	table = append(table, t.body...)
	t.Y = make([]int, len(table))
	t.X = make([]int, len(table[0]))
	for i := range table {
		for j := range table[i] {
			y := table[i][j].dims.Size.Y
			if y > t.Y[i] {
				t.Y[i] = y
			}
			x := table[i][j].dims.Size.X
			if x > t.X[j] {
				t.X[j] = x
			}
		}
	}
	// find the columns that will be rendered into the existing width
	maxWidth := gtx.Constraints.Max.X
	priorities := t.header.GetPriority()
	var runningTotal, prev int
	var columnsToRender []int
	for i := range priorities {
		prev = runningTotal
		runningTotal += t.header[priorities[i]].dims.Size.X
		if runningTotal > maxWidth {
			break
		}
		columnsToRender = append(columnsToRender, priorities[i])
	}
	// All fields will be expanded by the following ratio to reach the target width
	expansionFactor := float32(maxWidth) / float32(prev)
	for i := range t.X {
		t.X[i] = int(float32(t.X[i]) * expansionFactor)
	}

	return l.Dimensions{}
}
