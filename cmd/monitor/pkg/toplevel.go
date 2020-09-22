package pkg

import "github.com/p9c/pod/pkg/gui"

func (s *State) TopLevelLayout() {
	s.FlexV(
		s.Header(),
		gui.Flexed(1, func() {
			s.FlexHStart(
				s.LogViewer(),
				s.Sidebar(),
			)
		}),
		s.BottomBar(),
	)
}
