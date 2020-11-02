package p9

import l "gioui.org/layout"

type ColumnRow struct {
	Label string
	W     l.Widget
}

type Rows []ColumnRow

type Column struct {
	th    *Theme
	rows  []ColumnRow
	font  string
	scale float32
}

func (th *Theme) Column(rows Rows, font string, scale float32) *Column {
	return &Column{th: th, rows: rows, font: font, scale: scale}
}

func (c *Column) Fn(gtx l.Context) l.Dimensions {
	max, list := c.List(gtx)
	out := c.th.SliceToWidget(list, l.Vertical)
	gtx.Constraints.Max.X = max
	return out(gtx)
}

func (c *Column) List(gtx l.Context) (max int, out []l.Widget) {
	le := func(gtx l.Context, index int) l.Dimensions {
		return c.th.Label().Text(c.rows[index].Label).Font(c.font).TextScale(c.scale).Fn(gtx)
	}
	// render the widgets onto a second context to get their dimensions
	gtx1 := CopyContextDimensions(gtx, gtx.Constraints.Max, l.Horizontal)
	// generate the dimensions for all the list elements
	dims := GetDimensionList(gtx1, len(c.rows), le)
	for i := range dims {
		if dims[i].Size.X > max {
			max = dims[i].Size.X
		}
	}
	for x := range c.rows {
		i := x
		out = append(out, func(gtx l.Context) l.Dimensions {
			return c.th.Inset(0.25, func(gtx l.Context) l.Dimensions {
				// gtx.Constraints.Max.Y = dims[i].Size.Y
				// gtx.Constraints.Min.Y = dims[i].Size.Y
				// gtx.Constraints.Min.X = max
				return c.th.Flex().AlignBaseline().
					Rigid(
						func(gtx l.Context) l.Dimensions {
							gtx.Constraints.Min.X = max // dims[i].Size.X
							// gtx.Constraints.Max.X = max
							return c.th.Inset(0.25,
								c.th.Label().
									Text(c.rows[i].Label).
									Font(c.font).
									TextScale(c.scale).Fn,
							).Fn(gtx)
						},
					).
					Rigid(
						c.rows[i].W,
					).
					Fn(gtx)
			}).Fn(gtx)
		})
	}
	// // render the widgets onto a second context to get their dimensions
	// gtx1 = CopyContextDimensions(gtx, gtx.Constraints.Max, l.Vertical)
	// dim := GetDimension(gtx1, c.th.SliceToWidget(out, l.Vertical))
	// max = dim.Size.X
	return
}
