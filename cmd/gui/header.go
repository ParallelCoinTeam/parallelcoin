package gwallet

//
//var (
//	submitBtn         = new(widget.Clickable)
//	tourBtn           = new(widget.Clickable)
//	welcomeBtn        = new(widget.Clickable)
//	navBtn            = new(widget.Clickable)
//	browseBtn         = new(widget.Clickable)
//	headerSearchInput = &widget.Editor{
//		SingleLine: true,
//		Submit:     true,
//	}
//
//
//	coinLogoTextImageOp = paint.ImageOp{}
//	coinLogoTextImage   image.Image
//	coinLogoImageOp     = paint.ImageOp{}
//	coinLogoImage       image.Image
//	currentPage         string
//	)
//
//func (g *GioWallet) header() func(gtx C) D {
//	return ContainerLayout(th.Colors["Info"], 0, 0, 0, 0, func(gtx C) D {
//		gtx.Constraints.Min.X = gtx.Constraints.Max.X
//		return lyt.Format(gtx, "hflex(middle,r(inset(0dp0dp0dp6dp,_)),r(inset(20dp30dp20dp3dp,_)))",
//			g.headerSearch(),
//			g.headerMenu(),
//		)
//	})
//}
//func (g *GioWallet) headerMenu() func(gtx C) D {
//	return func(gtx C) D {
//		for tourBtn.Clicked() {
//
//			currentPage = "Welcome"
//		}
//		return lyt.Format(gtx, "hflex(middle,r(_),r(_),r(_))",
//			g.pageButton(tourBtn, func() {}, "StrokeCase", ""),
//			helper.DuoUIline(true, 0, 2, 2, th.Colors["DarkGrayI"]),
//			//g.pageButton(tourBtn, func() {}, "GlyphPencil", "Welcome"),
//			func(gtx C) D {
//				btn := material.IconButton(th.T, welcomeBtn, th.Icons["GlyphPencil"])
//				btn.Inset = layout.Inset{unit.Dp(2), unit.Dp(2), unit.Dp(2), unit.Dp(2)}
//				btn.Size = unit.Dp(21)
//				btn.Background = helper.HexARGB(th.Colors["Secondary"])
//				for welcomeBtn.Clicked() {
//					//f()
//					currentPage = "Welcome"
//				}
//				return btn.Layout(gtx)
//			},
//		)
//	}
//
//}
//
//func (g *GioWallet) headerSearch() func(gtx C) D {
//	return func(gtx C) D {
//		return lyt.Format(gtx, "hflex(middle,r(inset(20dp0dp20dp30dp,_)),r(_))",
//			ContainerLayout(th.Colors["White"], 8, 8, 8, 8, func(gtx C) D {
//				gtx.Constraints.Min.X = 430
//				e := material.Editor(th.T, headerSearchInput, "QmHash")
//				return e.Layout(gtx)
//			}),
//			func(gtx C) D {
//				e := material.Button(th.T, browseBtn, "Browse")
//				e.Inset = layout.Inset{
//					Top:    unit.Dp(4),
//					Right:  unit.Dp(4),
//					Bottom: unit.Dp(4),
//					Left:   unit.Dp(4),
//				}
//				e.CornerRadius = unit.Dp(4)
//				return e.Layout(gtx)
//			},
//		)
//	}
//}
