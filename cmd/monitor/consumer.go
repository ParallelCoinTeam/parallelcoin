package monitor

//
//func (s *State) Consume() {
//	sc := consume.Log(s.Ctx.KillAll, func(ent *logi.Entry) (err error) {
//		Debugf("from child: %s '%s'", ent.Level, ent.Text)
//		return
//	})
//	consume.Start(sc)
//out:
//	for {
//		select {
//		case <-s.Ctx.KillAll:
//			consume.Stop(sc)
//			break out
//		}
//	}
//}
