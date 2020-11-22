package p9

import l "gioui.org/layout"

type TextTableHeader []string

type TextTableRow []string

type TextTableBody []TextTableRow

type TextTable struct {
	*Theme
	Header           TextTableHeader
	Body             TextTableBody
	HeaderColor      string
	HeaderBackground string
	HeaderFont       string
	HeaderFontScale  float32
	CellColor        string
	CellBackground   string
	CellFont         string
	CellFontScale    float32
	Inset            float32
	List             *List
}

func (tt *TextTable) Fn(gtx l.Context) l.Dimensions {
	// set defaults if unset
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
	var header CellRow
	for i := range tt.Header {
		header = append(header, Cell{
			Widget: // tt.Theme.Fill(tt.HeaderBackground,
			tt.Theme.Inset(tt.Inset,
				tt.Theme.Body1(tt.Header[i]).
					Color(tt.HeaderColor).
					TextScale(tt.HeaderFontScale).
					Font(tt.HeaderFont).MaxLines(1).
					Fn,
			).Fn,
			// ).Fn,
		})
	}
	var body CellGrid
	for i := range tt.Body {
		row := CellRow{}
		for j := range tt.Body[i] {
			row = append(row, Cell{
				Widget: tt.Theme.Inset(0.25,
					tt.Theme.Body1(tt.Body[i][j]).
						Color(tt.CellColor).
						TextScale(tt.CellFontScale).
						Font(tt.CellFont).MaxLines(1).
						Fn,
				).Fn,
			})
		}
		body = append(body, row)
	}
	table := Table{
		th:               tt.Theme,
		header:           header,
		body:             body,
		list:             tt.List,
		headerBackground: tt.HeaderBackground,
		cellBackground:   tt.CellBackground,
	}
	return table.Fn(gtx)
}
