package gel

import (
	"image"
	"sort"
	
	l "github.com/p9c/gio/layout"
	"github.com/p9c/gio/op"
)

type Cell struct {
	l.Widget
	dims     l.Dimensions
	computed bool
	// priority only has meaning for the header row in defining an order of eliminating elements to fit a width.
	// When trimming size to fit width add from highest to lowest priority and stop when dimensions exceed the target.
	Priority int
}

func (c *Cell) getWidgetDimensions(gtx l.Context) {
	if c.Widget == nil {
		// this happens when new items are added if a frame reads the cell, it just - can't - be rendered!
		return
	}
	if c.computed {
		return
	}
	// gather the dimensions of the list elements
	gtx.Ops.Reset()
	child := op.Record(gtx.Ops)
	c.dims = c.Widget(gtx)
	c.computed = true
	_ = child.Stop()
	return
}

type CellRow []Cell

func (c CellRow) GetPriority() (out CellPriorities) {
	for i := range c {
		var cp CellPriority
		cp.Priority = c[i].Priority
		cp.Column = i
		out = append(out, cp)
	}
	sort.Sort(out)
	return
}

type CellPriority struct {
	Column   int
	Priority int
}

type CellPriorities []CellPriority

// Len sorts a cell row by priority
func (c CellPriorities) Len() int {
	return len(c)
}
func (c CellPriorities) Less(i, j int) bool {
	return c[i].Priority < c[j].Priority
}
func (c CellPriorities) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type CellGrid []CellRow

// Table is a super simple table widget that finds the dimensions of all cells, sets all to max of each axis, and then
// scales the remaining space evenly
type Table struct {
	*Window
	header           CellRow
	body             CellGrid
	list             *List
	Y, X             []int
	headerBackground string
	cellBackground   string
	reverse          bool
}

func (w *Window) Table() *Table {
	return &Table{
		Window: w,
		list:   w.List(),
	}
}

func (t *Table) SetReverse(color string) *Table {
	t.reverse = true
	return t
}

func (t *Table) HeaderBackground(color string) *Table {
	t.headerBackground = color
	return t
}

func (t *Table) CellBackground(color string) *Table {
	t.cellBackground = color
	return t
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
	// D.Ln(len(t.body), len(t.header))
	if len(t.header) == 0 {
		return l.Dimensions{}
	}
	// if len(t.body) == 0 || len(t.header) == 0 {
	// 	return l.Dimensions{}
	// }
	
	for i := range t.body {
		if len(t.header) != len(t.body[i]) {
			// this should never happen hence panic
			panic("not all rows are equal number of cells")
		}
	}
	gtx1 := CopyContextDimensionsWithMaxAxis(gtx, l.Vertical)
	gtx1.Constraints.Max = image.Point{X: Inf, Y: Inf}
	// gather the dimensions from all cells
	for i := range t.header {
		t.header[i].getWidgetDimensions(gtx1)
	}
	// D.S(t.header)
	for i := range t.body {
		for j := range t.body[i] {
			t.body[i][j].getWidgetDimensions(gtx1)
		}
	}
	// D.S(t.body)
	
	// find the max of each row and column
	var table CellGrid
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
	// // D.S(t.Y)
	// D.S(t.X)
	var total int
	for i := range t.X {
		total += t.X[i]
	}
	// D.S(t.X)
	// D.Ln(total)
	maxWidth := gtx.Constraints.Max.X
	for i := range t.X {
		t.X[i] = int(float32(t.X[i]) * float32(maxWidth) / float32(total))
	}
	// D.S(t.X)
	// D.Ln(maxWidth)
	// // find the columns that will be rendered into the existing width
	// // D.S(t.header)
	// priorities := t.header.GetPriority()
	// // D.S(priorities)
	// var runningTotal, prev int
	// columnsToRender := make([]int, 0)
	// for i := range priorities {
	// 	prev = runningTotal
	// 	x := t.header[priorities[i].Column].dims.Size.X
	// 	// D.Ln(priorities[i], x)
	// 	runningTotal += x
	//
	// 	if runningTotal > maxWidth {
	// 		// D.Ln(runningTotal, prev, maxWidth)
	// 		break
	// 	}
	// 	columnsToRender = append(columnsToRender, priorities[i].Column)
	// }
	// // txsort the columns to render into their original order
	// txsort.Ints(columnsToRender)
	// // D.S(columnsToRender)
	// // D.Ln(len(columnsToRender))
	// // All fields will be expanded by the following ratio to reach the target width
	// expansionFactor := float32(maxWidth) / float32(prev)
	// outColWidths := make([]int, len(columnsToRender))
	// for i := range columnsToRender {
	// 	outColWidths[i] = int(float32(t.X[columnsToRender[i]]) * expansionFactor)
	// }
	// // D.Ln(outColWidths)
	// // assemble the grid to be rendered as a two dimensional slice
	// grid := make([][]l.Widget, len(t.body)+1)
	// for i := 0; i < len(columnsToRender); i++ {
	// 	grid[0] = append(grid[0], t.header[columnsToRender[i]].Widget)
	// }
	// // for i := 0; i < len(columnsToRender); i++ {
	// // 	for j := range t.body[i] {
	// // 		grid[i+1] = append(grid[i+1], t.body[i][j].Widget)
	// // 	}
	// // }
	// // D.S(grid)
	// // assemble each row into a flex
	// out := make([]l.Widget, len(grid))
	// for i := range grid {
	// 	outFlex := t.Theme.Flex()
	// 	for jj, j := range grid[i] {
	// 		x := j
	// 		_ = jj
	// 		// outFlex.Rigid(x)
	// 		outFlex.Rigid(func(gtx l.Context) l.Dimensions {
	// 			// lock the cell to the calculated width.
	// 			gtx.Constraints.Max.X = outColWidths[jj]
	// 			gtx.Constraints.Min.X = gtx.Constraints.Max.X
	// 			return x(gtx)
	// 		})
	// 	}
	// 	out[i] = outFlex.Fn
	// }
	header := t.Theme.Flex() // .SpaceEvenly()
	for x, oi := range t.header {
		i := x
		// header is not in the list but drawn above it
		oie := oi
		txi := t.X[i]
		tyi := t.Y[0]
		header.Rigid(func(gtx l.Context) l.Dimensions {
			cs := gtx.Constraints
			cs.Max.X = txi
			cs.Min.X = gtx.Constraints.Max.X
			cs.Max.Y = tyi
			cs.Min.Y = gtx.Constraints.Max.Y
			// gtx.Constraints.Constrain(image.Point{X: txi, Y: tyi})
			dims := t.Fill(t.headerBackground, l.Center, t.TextSize.V, 0, EmptySpace(txi, tyi)).Fn(gtx)
			oie.Widget(gtx)
			return dims
		})
	}
	
	var out CellGrid
	out = CellGrid{t.header}
	if t.reverse {
		// append the body elements in reverse order stored
		lb := len(t.body) - 1
		for i := range t.body {
			out = append(out, t.body[lb-i])
		}
	} else {
		out = append(out, t.body...)
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		f := t.Theme.Flex() // .SpaceEvenly()
		oi := out[index]
		for x, oiee := range oi {
			i := x
			if index == 0 {
				// we skip the header, not implemented but the header could be part of the scrollable area if need
				// arises later, unwrap this block on a flag
			} else {
				if index >= len(t.Y) {
					break
				}
				oie := oiee
				txi := t.X[i]
				tyi := t.Y[index]
				f.Rigid(t.Fill(t.cellBackground, l.Center, t.TextSize.V, 0, func(gtx l.Context) l.Dimensions {
					cs := gtx.Constraints
					cs.Max.X = txi
					cs.Min.X = gtx.Constraints.Max.X
					cs.Max.Y = tyi
					cs.Min.Y = gtx.Constraints.Max.Y // gtx.Constraints.Constrain(image.Point{
					// 	X: t.X[i],
					// 	Y: t.Y[index],
					// })
					gtx.Constraints.Max.X = txi
					// gtx.Constraints.Min.X = gtx.Constraints.Max.X
					gtx.Constraints.Max.Y = tyi
					// gtx.Constraints.Min.Y = gtx.Constraints.Max.Y
					dims := EmptySpace(txi, tyi)(gtx)
					// dims
					oie.Widget(gtx)
					return dims
				}).Fn)
			}
		}
		return f.Fn(gtx)
	}
	return t.Theme.VFlex().
		Rigid(func(gtx l.Context) l.Dimensions {
			// header is fixed to the top of the widget
			return t.Fill(t.headerBackground, l.Center, t.TextSize.V, 0, header.Fn).Fn(gtx)
		}).
		Flexed(1,
			t.Fill(t.cellBackground, l.Center, t.TextSize.V, 0, func(gtx l.Context) l.Dimensions {
				return t.list.Vertical().
					Length(len(out)).
					Background(t.cellBackground).
					ListElement(le).
					Fn(gtx)
			}).Fn,
		).
		Fn(gtx)
}
