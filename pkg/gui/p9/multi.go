package p9

import (
	l "gioui.org/layout"
	"github.com/urfave/cli"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

type Multi struct {
	*Theme
	lines            *cli.StringSlice
	clickables       []*Clickable
	buttons          []*ButtonLayout
	input            *Input
	inputLocation    int
	addClickable     *Clickable
	removeClickables []*Clickable
	removeButtons    []*IconButton
}

func (th *Theme) Multiline(txt *cli.StringSlice, borderColorFocused, borderColorUnfocused string,
	size int, handle func(txt []string)) (m *Multi) {
	if handle == nil {
		handle = func(txt []string) {
			Debug(txt)
		}
	}
	addClickable := th.Clickable()
	m = &Multi{
		Theme:         th,
		lines:         txt,
		inputLocation: -1,
		addClickable:  addClickable,
	}
	handleChange := func(txt string) {
		(*m.lines)[m.inputLocation] = txt
		// after submit clear the editor
		m.inputLocation = -1
		handle(*m.lines)
	}
	m.input = th.Input("", borderColorFocused, borderColorUnfocused, size, handleChange)
	m.clickables = append(m.clickables, (*Clickable)(nil))
	m.buttons = append(m.buttons, (*ButtonLayout)(nil))
	m.removeClickables = append(m.clickables, (*Clickable)(nil))
	m.removeButtons = append(m.removeButtons, (*IconButton)(nil))
	for i, v := range *m.lines {
		Debug("making clickables")
		clickable := m.Theme.Clickable().SetClick(
			func() {
				m.inputLocation = i
			})
		if len(*m.lines) > len(m.clickables) {
			m.clickables = append(m.clickables, clickable)
		} else {
			m.clickables[i] = clickable
		}
		Debug("making button")
		btn := m.Theme.ButtonLayout(clickable).
			Embed(
				m.Theme.Flex().
					Rigid(
						m.Theme.Fill("Primary",
							m.Theme.Inset(0.5,
								m.Theme.Body2(v).Color("DocBg").Fn,
							).Fn,
						).Fn,
					).Fn,
			)
		if len(*m.lines) > len(m.buttons) {
			m.buttons = append(m.buttons, btn)
		} else {
			m.buttons[i] = btn
		}
		Debug("making clickables")
		removeClickable := m.Theme.Clickable()
		if len(*m.lines) > len(m.removeClickables) {
			m.removeClickables = append(m.clickables, removeClickable)
		} else {
			m.removeClickables[i] = removeClickable
		}
		Debug("making remove button")
		removeBtn := m.Theme.IconButton(removeClickable).
			Icon(
				m.Theme.Icon().Scale(1.5).Color("DocText").Src(icons.ActionDelete),
			).
			Background("DocBg").
			SetClick(func() {
				Debug("remove button", i, "clicked")
				m.inputLocation = -1
				if i == len(*m.lines)-1 {
					*m.lines = (*m.lines)[:len(*m.lines)-1]
					m.clickables = m.clickables[:len(m.clickables)-1]
					m.buttons = m.buttons[:len(m.buttons)-1]
					m.removeClickables = m.removeClickables[:len(m.removeClickables)-1]
					m.removeButtons = m.removeButtons[:len(m.removeButtons)-1]
				} else {
					*m.lines = append((*m.lines)[:i], (*m.lines)[i+1:]...)
					m.clickables = append(m.clickables[:i], m.clickables[i+1:]...)
					m.buttons = append(m.buttons[:i], m.buttons[i+1:]...)
					m.removeClickables = append(m.removeClickables[:i], m.removeClickables[i+1:]...)
					m.removeButtons = append(m.removeButtons[:i], m.removeButtons[i+1:]...)
				}
			})
		if len(*m.lines) > len(m.removeButtons) {
			m.removeButtons = append(m.removeButtons, removeBtn)
		} else {
			m.removeButtons[i] = removeBtn
		}
	}
	return m
}

func (m *Multi) Fn(gtx l.Context) l.Dimensions {
	// m.removeClickables = m.removeClickables[:0]
	// m.removeButtons = m.removeButtons[:0]
	// m.clickables = m.clickables[:0]
	// m.buttons = m.buttons[:0]
	if len(m.clickables) < len(*m.lines) {
		Debug("making new clickables")
		m.clickables = append(m.clickables, (*Clickable)(nil))
	}
	if len(m.buttons) < len(*m.lines) {
		Debug("making new buttons")
		m.buttons = append(m.buttons, (*ButtonLayout)(nil))
	}
	if len(m.removeClickables) < len(*m.lines) {
		Debug("making new removeClickables")
		m.removeClickables = append(m.clickables, (*Clickable)(nil))
	}
	if len(m.removeButtons) < len(*m.lines) {
		Debug("making new removeButtons")
		m.removeButtons = append(m.removeButtons, (*IconButton)(nil))
	}
	for i, v := range *m.lines {
		if m.clickables[i] == nil {
			// Debug("making clickables")
			clickable := m.Theme.Clickable().SetClick(
				func() {
					m.inputLocation = i
				})
			if len(*m.lines) > len(m.clickables) {
				m.clickables = append(m.clickables, clickable)
			} else {
				m.clickables[i] = clickable
			}
		}
		if m.buttons[i] == nil {
			// Debug("making button")
			btn := m.Theme.ButtonLayout(m.clickables[i]).
				Embed(
					m.Theme.Flex().
						Rigid(
							m.Theme.Fill("Primary",
								m.Theme.Inset(0.5,
									m.Theme.Body2(v).Color("DocBg").Fn,
								).Fn,
							).Fn,
						).Fn,
				)
			if len(*m.lines) > len(m.buttons) {
				m.buttons = append(m.buttons, btn)
			} else {
				m.buttons[i] = btn
			}
		}
		if m.removeClickables[i] == nil {
			// Debug("making clickables")
			removeClickable := m.Theme.Clickable()
			if len(*m.lines) > len(m.removeClickables) {
				m.removeClickables = append(m.clickables, removeClickable)
			} else {
				m.removeClickables[i] = removeClickable
			}
		}
		if m.removeButtons[i] == nil {
			// Debug("making remove button")
			removeBtn := m.Theme.IconButton(m.removeClickables[i]).
				Icon(
					m.Theme.Icon().Scale(1.5).Color("DocText").Src(icons.ActionDelete),
				).
				Background("DocBg").
				SetClick(func() {
					Debug("remove button", i, "clicked")
					m.inputLocation = -1
					if len(*m.lines)-1 == i {
						*m.lines = (*m.lines)[:len(*m.lines)-1]
					} else {
						*m.lines = append((*m.lines)[:i], (*m.lines)[i+1:]...)
					}
				})
			if len(*m.lines) > len(m.removeButtons) {
				m.removeButtons = append(m.removeButtons, removeBtn)
			} else {
				m.removeButtons[i] = removeBtn
			}
		}
	}
	addButton := m.Theme.IconButton(m.addClickable).Icon(
		m.Theme.Icon().Scale(1.5).Color("Primary").Src(icons.ContentAdd),
	).SetClick(func() {
		*m.lines = append(*m.lines, "0")
		Debug("clicked addButton")
		Debugs([]string(*m.lines))
		m.inputLocation = len(*m.lines) - 1
	})
	var widgets []l.Widget
	if m.inputLocation > 0 && m.inputLocation < len(*m.lines) {
		m.input.Editor().SetText((*m.lines)[m.inputLocation])
	}
	for i := range m.buttons {
		if m.buttons[i] == nil {
			continue
		}
		if i == m.inputLocation {
			input := m.Flex().
				Rigid(
					m.removeButtons[m.inputLocation].Fn,
				).
				Rigid(
					m.input.Fn,
				).
				Fn
			widgets = append(widgets, input)
		} else {
			button := m.Flex().
				Rigid(
					m.removeButtons[i].Fn,
				).
				Rigid(
					m.buttons[i].Fn,
				).
				Fn
			widgets = append(widgets, button)
		}
	}
	widgets = append(widgets, addButton.Background("DocBg").Fn)
	out := m.Theme.VFlex()
	for i := range widgets {
		out.Rigid(widgets[i])
	}
	return out.Fn(gtx)
}
