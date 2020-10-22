package p9

import l "gioui.org/layout"

func (th *Theme) Card(background string, w l.Widget) func(gtx l.Context) l.Dimensions {
	return th.Inset(0.25,
		th.Fill(background,
			th.Inset(0.5,
				w,
			).Fn,
		).Fn,
	).Fn
}

func (th *Theme) CardList(list *List, background string, widgets ...l.Widget) func(gtx l.Context) l.Dimensions {
	out := list.Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
		return th.Card(background, widgets[index])(gtx)
	}).Length(len(widgets))
	return out.Fn
}

func (th *Theme) CardContent(title, color string, w l.Widget) func(gtx l.Context) l.Dimensions {
	out := th.VFlex()
	if title != "" {
		out.Rigid(th.H5(title).Color(color).Fn)
	}
	out.Rigid(w)
	return out.Fn
}
