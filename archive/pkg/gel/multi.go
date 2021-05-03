package gel

import (
	l "github.com/p9c/gio/layout"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Multi struct {
	*Window
	lines            *[]string
	clickables       []*Clickable
	buttons          []*ButtonLayout
	input            *Input
	inputLocation    int
	addClickable     *Clickable
	removeClickables []*Clickable
	removeButtons    []*IconButton
	handle           func(txt []string)
}

func (w *Window) Multiline(
	txt *[]string,
	borderColorFocused, borderColorUnfocused, backgroundColor string,
	size float32,
	handle func(txt []string),
) (m *Multi) {
	if handle == nil {
		handle = func(txt []string) {
			D.Ln(txt)
		}
	}
	addClickable := w.Clickable()
	m = &Multi{
		Window:        w,
		lines:         txt,
		inputLocation: -1,
		addClickable:  addClickable,
		handle:        handle,
	}
	handleChange := func(txt string) {
		D.Ln("handleChange", m.inputLocation)
		(*m.lines)[m.inputLocation] = txt
		// after submit clear the editor
		m.inputLocation = -1
		m.handle(*m.lines)
	}
	m.input = w.Input("", "", borderColorFocused, borderColorUnfocused, backgroundColor, handleChange, nil)
	m.clickables = append(m.clickables, (*Clickable)(nil))
	// m.buttons = append(m.buttons, (*ButtonLayout)(nil))
	m.removeClickables = append(m.removeClickables, (*Clickable)(nil))
	m.removeButtons = append(m.removeButtons, (*IconButton)(nil))
	for i := range *m.lines {
		// D.Ln("making clickables")
		x := i
		clickable := m.Clickable().SetClick(
			func() {
				m.inputLocation = x
				D.Ln("button clicked", x, m.inputLocation)
			})
		if len(*m.lines) > len(m.clickables) {
			m.clickables = append(m.clickables, clickable)
		} else {
			m.clickables[i] = clickable
		}
		// D.Ln("making button")
		btn := m.ButtonLayout(clickable).CornerRadius(0).Background(
			backgroundColor).
			Embed(
				m.Theme.Flex().AlignStart().
					Flexed(1,
						m.Fill("Primary", l.Center, m.TextSize.V, 0, m.Inset(0.25,
							m.Body2((*m.lines)[i]).Color("DocText").Fn,
						).Fn).Fn,
					).Fn,
			)
		if len(*m.lines) > len(m.buttons) {
			m.buttons = append(m.buttons, btn)
		} else {
			m.buttons[i] = btn
		}
		// D.Ln("making clickables")
		removeClickable := m.Clickable()
		if len(*m.lines) > len(m.removeClickables) {
			m.removeClickables = append(m.removeClickables, removeClickable)
		} else {
			m.removeClickables[i] = removeClickable
		}
		// D.Ln("making remove button")
		y := i
		removeBtn := m.IconButton(removeClickable).
			Icon(
				m.Icon().Scale(1.5).Color("DocText").Src(&icons.ActionDelete),
			).
			Background("").
			SetClick(func() {
				D.Ln("remove button", y, "clicked", len(*m.lines))
				m.inputLocation = -1
				if len(*m.lines)-1 == y {
					*m.lines = (*m.lines)[:len(*m.lines)-1]
				} else if len(*m.lines)-2 == y {
					*m.lines = (*m.lines)[:len(*m.lines)-2]
				} else {
					*m.lines = append((*m.lines)[:y+1], (*m.lines)[y+2:]...)
				}
				m.handle(*m.lines)
				// D.Ln("remove button", i, "clicked")
				// m.inputLocation = -1
				// ll := len(*m.lines)-1
				// if i == ll {
				// 	*m.lines = (*m.lines)[:len(*m.lines)-1]
				// 	m.clickables = m.clickables[:len(m.clickables)-1]
				// 	m.buttons = m.buttons[:len(m.buttons)-1]
				// 	m.removeClickables = m.removeClickables[:len(m.removeClickables)-1]
				// 	m.removeButtons = m.removeButtons[:len(m.removeButtons)-1]
				// } else {
				// 	if len(*m.lines)-1 < i {
				// 		return
				// 	}
				// 	*m.lines = append((*m.lines)[:i], (*m.lines)[i+1:]...)
				// 	m.clickables = append(m.clickables[:i], m.clickables[i+1:]...)
				// 	m.buttons = append(m.buttons[:i], m.buttons[i+1:]...)
				// 	m.removeClickables = append(m.removeClickables[:i], m.removeClickables[i+1:]...)
				// 	m.removeButtons = append(m.removeButtons[:i], m.removeButtons[i+1:]...)
				// }
			})
		if len(*m.lines) > len(m.removeButtons) {
			m.removeButtons = append(m.removeButtons, removeBtn)
		} else {
			m.removeButtons[x] = removeBtn
		}
	}
	return m
}

func (m *Multi) UpdateWidgets() *Multi {
	if len(m.clickables) < len(*m.lines) {
		D.Ln("allocating new clickables")
		m.clickables = append(m.clickables, (*Clickable)(nil))
	}
	if len(m.buttons) < len(*m.lines) {
		D.Ln("allocating new buttons")
		m.buttons = append(m.buttons, (*ButtonLayout)(nil))
	}
	if len(m.removeClickables) < len(*m.lines) {
		D.Ln("allocating new removeClickables")
		m.removeClickables = append(m.clickables, (*Clickable)(nil))
	}
	if len(m.removeButtons) < len(*m.lines) {
		D.Ln("allocating new removeButtons")
		m.removeButtons = append(m.removeButtons, (*IconButton)(nil))
	}
	return m
}

func (m *Multi) PopulateWidgets() *Multi {
	added := false
	for i := range *m.lines {
		if m.clickables[i] == nil {
			added = true
			D.Ln("making clickables", i)
			x := i
			m.clickables[i] = m.Clickable().SetClick(
				func() {
					D.Ln("clicked", x, m.inputLocation)
					m.inputLocation = x
					m.input.editor.SetText((*m.lines)[x])
					m.input.editor.Focus()
					// m.input.editor.SetFocus(func(is bool) {
					// 	if !is {
					// 		m.inputLocation = -1
					// 	}
					// })
				})
		}
		// m.clickables[i]
		if m.buttons[i] == nil {
			added = true
			btn := m.ButtonLayout(m.clickables[i]).CornerRadius(0).Background("Transparent")
			m.buttons[i] = btn
		}
		m.buttons[i].Embed(
			m.Theme.Flex().
				Rigid(
					m.Inset(0.25,
						m.Body2((*m.lines)[i]).Color("DocText").Fn,
					).Fn,
				).Fn,
		)
		if m.removeClickables[i] == nil {
			added = true
			removeClickable := m.Clickable()
			m.removeClickables[i] = removeClickable
		}
		if m.removeButtons[i] == nil {
			added = true
			D.Ln("making remove button", i)
			x := i
			m.removeButtons[i] = m.IconButton(m.removeClickables[i].
				SetClick(func() {
					D.Ln("remove button", x, "clicked", len(*m.lines))
					m.inputLocation = -1
					if len(*m.lines)-1 == i {
						*m.lines = (*m.lines)[:len(*m.lines)-1]
					} else {
						*m.lines = append((*m.lines)[:x], (*m.lines)[x+1:]...)
					}
					m.handle(*m.lines)
				})).
				Icon(
					m.Icon().Scale(1.5).Color("DocText").Src(&icons.ActionDelete),
				).
				Background("")
		}
	}
	if added {
		D.Ln("clearing editor")
		m.input.editor.SetText("")
		m.input.editor.Focus()
	}
	return m
}

func (m *Multi) Fn(gtx l.Context) l.Dimensions {
	m.UpdateWidgets()
	m.PopulateWidgets()
	addButton := m.IconButton(m.addClickable).Icon(
		m.Icon().Scale(1.5).Color("Primary").Src(&icons.ContentAdd),
	)
	var widgets []l.Widget
	if m.inputLocation > 0 && m.inputLocation < len(*m.lines) {
		m.input.Editor().SetText((*m.lines)[m.inputLocation])
	}
	for i := range *m.lines {
		if m.buttons[i] == nil {
			x := i
			btn := m.ButtonLayout(m.clickables[i].SetClick(
				func() {
					D.Ln("button pressed", (*m.lines)[x], x, m.inputLocation)
					m.inputLocation = x
					m.input.editor.SetText((*m.lines)[x])
					m.input.editor.Focus()
				})).CornerRadius(0).Background("Transparent").
				Embed(
					m.Theme.Flex().
						Rigid(
							m.Inset(0.25,
								m.Body2((*m.lines)[x]).Color("DocText").Fn,
							).Fn,
						).Fn,
				)
			m.buttons[i] = btn
		}
		if i == m.inputLocation {
			m.input.Editor().SetText((*m.lines)[i])
			input := m.Flex().
				Rigid(
					m.removeButtons[i].Fn,
				).
				Flexed(1,
					m.input.Fn,
				).
				Fn
			widgets = append(widgets, input)
		} else {
			x := i
			m.clickables[i].SetClick(
				func() {
					D.Ln("setting", x, m.inputLocation)
					m.inputLocation = x
					m.input.editor.SetText((*m.lines)[x])
					m.input.editor.Focus()
				})
			button := m.Flex().AlignStart().
				Rigid(
					m.removeButtons[i].Fn,
				).
				Flexed(1,
					m.buttons[i].Fn,
				).
				Fn
			widgets = append(widgets, button)
		}
	}
	widgets = append(widgets, addButton.SetClick(func() {
		D.Ln("clicked add")
		*m.lines = append(*m.lines, "")
		m.inputLocation = len(*m.lines) - 1
		D.S([]string(*m.lines))
		m.UpdateWidgets()
		m.PopulateWidgets()
		m.input.editor.SetText("")
		m.input.editor.Focus()
	}).Background("").Fn)
	// m.UpdateWidgets()
	// m.PopulateWidgets()
	// D.Ln(m.inputLocation)
	// if m.inputLocation > 0 {
	// 	m.input.Editor().Focus()
	// }
	out := m.Theme.VFlex()
	for i := range widgets {
		out.Rigid(widgets[i])
	}
	return out.Fn(gtx)
}

func (m *Multi) Widgets() (widgets []l.Widget) {
	m.UpdateWidgets()
	m.PopulateWidgets()
	if m.inputLocation > 0 && m.inputLocation < len(*m.lines) {
		m.input.Editor().SetText((*m.lines)[m.inputLocation])
	}
	focusFunc := func(is bool) {
		mi := m.inputLocation
		D.Ln("editor", "is focused", is)
		// debug.PrintStack()
		if !is {
			m.input.borderColor = m.input.borderColorUnfocused
			// submit the current edit if any
			txt := m.input.editor.Text()
			cur := (*m.lines)[m.inputLocation]
			if txt != cur {
				D.Ln("changed text")
				// run submit hook
				m.input.editor.submitHook(txt)
			} else {
				D.Ln("text not changed")
				// When a new item is added this unfocus event occurs and this makes it behave correctly
				// Normally the editor would not be rendered if not focused so setting it to focus does no harm in the
				// case of switching to another
			}
			// m.inputLocation = -1
		} else {
			m.input.borderColor = m.input.borderColorFocused
			m.inputLocation = mi
			// m.input.editor.Focus()
		}
	}
	m.input.editor.SetFocus(focusFunc)
	for ii := range *m.lines {
		i := ii
		// D.Ln("iterating lines", i, len(*m.lines))
		if m.buttons[i] == nil {
			D.Ln("making new button layout")
			btn := m.ButtonLayout(m.clickables[i].SetClick(
				func() {
					D.Ln("button pressed", (*m.lines)[i], i, m.inputLocation)
					m.UpdateWidgets()
					m.PopulateWidgets()
					m.inputLocation = i
					m.input.editor.SetText((*m.lines)[i])
					m.input.editor.Focus()
				})).CornerRadius(0).Background("").
				Embed(
					func(gtx l.Context) l.Dimensions {
						return m.Theme.Flex().
							Flexed(1,
								m.Inset(0.25,
									m.Body2((*m.lines)[i]).Color("DocText").Fn,
								).Fn,
							).Fn(gtx)
					},
				)
			m.buttons[i] = btn
		}
		if i == m.inputLocation {
			// x := i
			// D.Ln("rendering editor", x)
			
			input := func(gtx l.Context) l.Dimensions {
				return m.Inset(0.25,
					m.Flex().
						Rigid(
							m.removeButtons[i].Fn,
						).
						Flexed(1,
							m.input.Fn,
						).
						Fn,
				).
					Fn(gtx)
			}
			widgets = append(widgets, input)
		} else {
			// D.Ln("rendering button", i)
			m.clickables[i].SetClick(
				func() {
					m.UpdateWidgets()
					m.PopulateWidgets()
					m.inputLocation = i
					m.input.editor.SetText((*m.lines)[i])
					m.input.editor.Focus()
					D.Ln("setting", i, m.inputLocation)
				})
			button := func(gtx l.Context) l.Dimensions {
				return m.Inset(0.25,
					m.Flex().AlignStart().
						Rigid(
							m.removeButtons[i].Fn,
						).
						Rigid(
							m.buttons[i].Fn,
						).
						Flexed(1, EmptyMaxWidth()).
						Fn,
				).Fn(gtx)
			}
			widgets = append(widgets, button)
		}
	}
	// D.Ln("widgets", widgets)
	addButton := func(gtx l.Context) l.Dimensions {
		addb :=
			m.Inset(0.25,
				m.Theme.Flex().AlignStart().
					Rigid(
						m.IconButton(
							m.addClickable).
							Icon(
								m.Icon().Scale(1.5).Color("Primary").Src(&icons.ContentAdd),
							).
							SetClick(func() {
								D.Ln("clicked add")
								m.inputLocation = len(*m.lines)
								*m.lines = append(*m.lines, "")
								m.input.editor.SetText("")
								D.S([]string(*m.lines))
								m.UpdateWidgets()
								m.PopulateWidgets()
								m.input.editor.Focus()
							}).
							Background("Transparent").
							Fn,
					).
					Flexed(1, EmptyMaxWidth()).
					Fn,
			).Fn
		widgets = append(widgets, addb)
		return addb(gtx)
	}
	widgets = append(widgets, addButton)
	return
}
