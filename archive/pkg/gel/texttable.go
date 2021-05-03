package gel

import l "github.com/p9c/gio/layout"

type TextTableHeader []string

type TextTableRow []string

type TextTableBody []TextTableRow

// TextTable is a widget that renders a scrolling list of rows of data labeled by a header. Note that for the reasons of
// expedience and performance this widget assumes a growing but immutable list of rows of items. If this is used on
// data that is not immutable, nilling the body will cause it to be wholly regenerated, updating older content than the
// longest length the list has reached.
type TextTable struct {
	*Window
	Header           TextTableHeader
	Body             TextTableBody
	HeaderColor      string
	HeaderDarkTheme  bool
	HeaderBackground string
	HeaderFont       string
	HeaderFontScale  float32
	CellColor        string
	CellBackground   string
	CellFont         string
	CellFontScale    float32
	CellInset        float32
	List             *List
	Table            *Table
}

// Regenerate the text table.
func (tt *TextTable) Regenerate(fully bool) {
	if len(tt.Header) == 0 || len(tt.Body) == 0 {
		return
	}
	// // set defaults if unset
	tt.SetDefaults()
	if tt.Table.header == nil || len(tt.Table.header) < 1 || tt.HeaderDarkTheme != tt.Theme.Dark.True() {
		tt.HeaderDarkTheme = tt.Theme.Dark.True()
		// if this is being regenerated due to theme change
		tt.Table.header = tt.Table.header[:0]
		// this only has to be created once
		for i := range tt.Header {
			tt.Table.header = append(tt.Table.header, Cell{
				Widget: // tt.Theme.Fill(tt.HeaderBackground,
				tt.Inset(tt.CellInset,
					tt.Body1(tt.Header[i]).
						Color(tt.HeaderColor).
						TextScale(tt.HeaderFontScale).
						Font(tt.HeaderFont).MaxLines(1).
						Fn,
				).Fn,
				// ).Fn,
			})
		}
	}
	// var startIndex int
	// if tt.Table.body == nil || len(tt.Table.body) < 1 {
	// 	// tt.Table.body = tt.Table.body[:0]
	// } else {
	// if fully {
	// tt.Body = tt.Body[:0]
	// tt.Table.body = tt.Table.body[:0]
	// }
	// startIndex = len(tt.Table.body)
	// D.Ln("startIndex", startIndex, len(tt.Body))
	// if startIndex < len(tt.Body) {
	
	// bd := tt.Body // [startIndex:]
	diff := len(tt.Body) - len(tt.Table.body)
	// D.Ln(len(tt.Table.body), len(tt.Body), diff)
	if diff > 0 {
		cg := make(CellGrid, diff)
		for i := range cg {
			cg[i] = make(CellRow, len(tt.Header))
		}
		tt.Table.body = append(tt.Table.body, cg...)
	}
	// D.Ln(len(tt.Table.body), len(tt.Body))
	var body CellGrid
	for i := range tt.Body {
		var row CellRow
		for j := range tt.Body[i] {
			tt.Table.body[i][j] = Cell{
				Widget: tt.Inset(0.25,
					tt.Body1(tt.Body[i][j]).
						Color(tt.CellColor).
						TextScale(tt.CellFontScale).
						Font(tt.CellFont).MaxLines(1).
						Fn,
				).Fn,
			}
		}
		body = append(body, row)
	}
	// tt.Table.body = append(tt.Table.body, body...)
	// }
	// }
}

func (tt *TextTable) SetReverse() *TextTable {
	tt.Table.reverse = true
	return tt
}

func (tt *TextTable) SetDefaults() *TextTable {
	if tt.HeaderColor == "" {
		tt.HeaderColor = "PanelText"
	}
	if tt.HeaderBackground == "" {
		tt.HeaderBackground = "PanelBg"
	}
	if tt.HeaderFont == "" {
		tt.HeaderFont = "bariol bold"
	}
	if tt.HeaderFontScale == 0 {
		tt.HeaderFontScale = Scales["Caption"]
	}
	if tt.CellColor == "" {
		tt.CellColor = "DocText"
	}
	if tt.CellBackground == "" {
		tt.CellBackground = "DocBg"
	}
	if tt.CellFont == "" {
		tt.CellFont = "go regular"
	}
	if tt.CellFontScale == 0 {
		tt.CellFontScale = Scales["Caption"]
	}
	// we assume the caller has intended a zero inset if it is zero
	if tt.Table == nil {
		tt.Table = &Table{
			Window:           tt.Window,
			list:             tt.List,
			headerBackground: tt.HeaderBackground,
			cellBackground:   tt.CellBackground,
		}
	}
	return tt
}

func (tt *TextTable) Fn(gtx l.Context) l.Dimensions {
	return tt.Table.Fn(gtx)
}
