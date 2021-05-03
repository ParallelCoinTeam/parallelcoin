package gel

import l "github.com/p9c/gio/layout"

func (w *Window) Card(background string, embed l.Widget,
) func(gtx l.Context) l.Dimensions {
	return w.Inset(0.0,
		w.Fill(background, l.Center, w.TextSize.V, 0, w.Inset(0.25,
			embed,
		).Fn).Fn,
	).Fn
}

func (w *Window) CardList(list *List, background string,
	widgets ...l.Widget) func(gtx l.Context) l.Dimensions {
	out := list.Vertical().ListElement(func(gtx l.Context, index int) l.Dimensions {
		return w.Card(background, widgets[index])(gtx)
	}).Length(len(widgets))
	return out.Fn
}

func (w *Window) CardContent(title, color string, embed l.Widget) func(gtx l.Context) l.Dimensions {
	out := w.VFlex()
	if title != "" {
		out.Rigid(w.H6(title).Color(color).Fn)
	}
	return out.Rigid(embed).Fn
}
