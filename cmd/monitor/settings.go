package monitor

import (
	"gioui.org/layout"

	"github.com/p9c/pod/cmd/gui/pages"
)

func (m *State) SettingsButtons() layout.FlexChild {
	return Flexed(1, func() {
		m.FlexH(Rigid(func() {
			bg, fg := "PanelBg", "PanelText"
			if m.Config.SettingsOpen {
				bg, fg = "DocBg", "DocText"
			}
			m.TextButton("Settings", "Secondary",
				23, fg, bg, m.SettingsFoldButton)
			for m.SettingsFoldButton.Clicked(m.Gtx) {
				L.Debug("settings folder clicked")
				switch {
				case !m.Config.SettingsOpen:
					m.Config.BuildOpen = false
					m.Config.SettingsOpen = true
				case m.Config.SettingsOpen:
					m.Config.SettingsOpen = false
				}
				m.SaveConfig()
			}
		}),
		)
	})
}

func (m *State) SettingsPage() layout.FlexChild {
	if !m.Config.SettingsOpen {
		return Flexed(0, func() {})
	}
	var weight float32 = 0.5
	var settingsInset = 0
	switch {
	case m.WindowWidth < 1024 && m.WindowHeight > 1024:
		// weight = 0.333
	case m.WindowHeight < 1024 && m.WindowWidth < 1024:
		weight = 1
	case m.WindowHeight < 600 && m.WindowWidth > 1024:
		weight = 1
	}
	return Flexed(weight, func() {
		m.Inset(settingsInset, func() {
			cs := m.Gtx.Constraints
			m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			// m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
			m.FlexV(Rigid(func() {
				cs := m.Gtx.Constraints
				m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
				m.Inset(4, func() {})
			}), Rigid(func() {
				m.FlexH(Rigid(func() {
					m.TextButton("Run Settings", "Secondary",
						23, "DocText", "DocBg",
						m.SettingsTitleCloseButton)
					for m.SettingsTitleCloseButton.Clicked(m.Gtx) {
						L.Debug("settings panel title close button clicked")
						m.Config.SettingsOpen = false
						m.SaveConfig()
					}
				}), Spacer(), Rigid(func() {
					if m.WindowWidth > 640 {
						m.SettingsHeader()()
					}
				}), Spacer(), Rigid(func() {
					m.IconButton("minimize", "DocText", "DocBg",
						m.SettingsCloseButton)
					for m.SettingsCloseButton.Clicked(m.Gtx) {
						L.Debug("settings panel close button clicked")
						m.Config.SettingsOpen = false
						m.SaveConfig()
					}
				}),
				)
			}),
				Rigid(func() {
					if m.WindowWidth < 640 {
						cs := m.Gtx.Constraints
						m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
						m.SettingsHeader()()
					}
				}),
				Flexed(1, func() {
					m.Inset(settingsInset, func() {
						cs := m.Gtx.Constraints
						m.Rectangle(cs.Width.Max, cs.Height.Max, "PanelBg")
						pages.SettingsBody(m.Rc, m.Gtx, m.Theme)
					})
				}),
				Rigid(func() {
					cs := m.Gtx.Constraints
					m.Rectangle(cs.Width.Max, cs.Height.Max, "DocBg")
					m.Inset(4, func() {})
				}),
			)
		})
	})
}
