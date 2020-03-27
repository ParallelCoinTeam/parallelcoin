package monitor

func (s *State) Consume() {
out:
	for {
		select {
		case <-s.Ctx.KillAll:
			break out
		}
	}
}
