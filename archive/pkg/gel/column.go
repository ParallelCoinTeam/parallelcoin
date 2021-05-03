package gel

import (
	l "github.com/p9c/gio/layout"
)

type ColumnRow struct {
	Label string
	W     l.Widget
}

type Rows []ColumnRow

type Column struct {
	*Window
	rows              []ColumnRow
	font              string
	scale             float32
	color, background string
	list              *List
}

func (w *Window) Column(rows Rows, font string, scale float32, color string,
	background string) *Column {
	return &Column{Window: w, rows: rows, font: font, scale: scale, color: color,
		background: background, list: w.List()}
}

func (c *Column) Fn(gtx l.Context) l.Dimensions {
	max, list := c.List(gtx)
	out := c.Theme.SliceToWidget(list, l.Vertical)
	gtx.Constraints.Max.X = max
	return out(gtx)
}

func (c *Column) List(gtx l.Context) (max int, out []l.Widget) {
	le := func(gtx l.Context, index int) l.Dimensions {
		return c.H6(c.rows[index].Label).Font(c.font).TextScale(c.scale).Fn(gtx)
	}
	// render the widgets onto a second context to get their dimensions
	gtx1 := CopyContextDimensionsWithMaxAxis(gtx, l.Vertical)
	// generate the dimensions for all the list elements
	dims := GetDimensionList(gtx1, len(c.rows), le)
	D.S(dims)
	for i := range dims {
		if dims[i].Size.X > max {
			max = dims[i].Size.X
		}
	}
	for x := range c.rows {
		i := x
		_ = i
		out = append(out, func(gtx l.Context) l.Dimensions {
			return c.Flex(). // AlignEnd().SpaceStart().
				Rigid(
					c.Fill("red", l.Center, c.TextSize.V, 0, EmptySpace(max-dims[i].Size.X, dims[i].Size.Y)).Fn,
				).
				Rigid(
					c.Inset(0.5, func(gtx l.Context) l.Dimensions {
						// D.Ln("max!", m)
						// // gtx.Constraints.Max.X = max
						// gtx.Constraints.Max.Y = dims[i].Size.Y
						// gtx.Constraints.Min.Y = dims[i].Size.Y
						// gtx.Constraints.Min.X = max
						// gtx.Constraints.Constrain(image.Point{X: max, Y: dims[i].Size.Y})
						// gtx.Constraints.Max.X = max
						// gtx.Constraints.Min.Y = 0
						// gtx.Constraints.Max.Y = dims[i].Size.Y
						// gtx.Constraints.Constrain(dims[i].Size)
						dms := c.Label().
							Text(c.rows[i].Label).
							Font(c.font).
							TextScale(c.scale).
							Color(c.color).
							Fn(gtx)
						return dms
					}).Fn,
				).
				Rigid(
					c.Inset(0.5,
						c.rows[i].W,
					).Fn,
				).
				Fn(gtx)
			// return c.Theme.Fill("Primary",
			// 	c.Theme.Flex().AlignEnd().SpaceBetween().
			// 		Rigid(
			// 		).Fn,
			// ).
			// Fn(gtx)
		})
	}
	// le = func(gtx l.Context, index int) l.Dimensions {
	// 	return out[index](gtx)
	// }
	return max, out
	
	// // render the widgets onto a second context to get their dimensions
	// gtx1 = CopyContextDimensionsWithMaxAxis(gtx, gtx.Constraints.Max, l.Vertical)
	// dim := GetDimension(gtx1, c.Theme.SliceToWidget(out, l.Vertical))
	// max = dim.Size.X
	// return
}
