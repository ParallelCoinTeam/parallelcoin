package component

//func TransactionsList(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {

// func ContentHeader(gtx *layout.Context, th *gelook.DuoUItheme, b func()) func() {
//	return func() {
//		hmin := gtx.Constraints.Width.Min
//		vmin := gtx.Constraints.Height.Min
//		layout.Stack{Alignment: layout.Center}.Layout(gtx,
//			layout.Expanded(func() {
//				clip.Rect{
//					Rect: f32.Rectangle{Max: f32.Point{
//						X: float32(gtx.Constraints.Width.Min),
//						Y: float32(gtx.Constraints.Height.Min),
//					}},
//				}.Op(gtx.Ops).Add(gtx.Ops)
//				fill(gtx, gelook.HexARGB(th.Colors["Primary"]))
//			}),
//			layout.Stacked(func() {
//				gtx.Constraints.Width.Min = hmin
//				gtx.Constraints.Height.Min = vmin
//				layout.UniformInset(unit.Dp(0)).Layout(gtx, b)
//			}),
//		)
//	}
// }

// func InputField(rc *rcd.RcVar, gtx *layout.Context, th *theme.DuoUItheme, f *Field) func() {
//	return func() {
//		e := th.DuoUIEditor(f.Field.Label)
//		e.Font.Typeface = th.Font.Primary
//		e.Font.Style = text.Italic
//		//e.Layout(gtx, (rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor))
//		//(ui.rc.Settings.Daemon.Widgets[f.Field.Label]).(*controller.Editor).SetText(f.Field.Value.(reflect.Value).String())
//		//Info(f.Field.Value.(reflect.Value).String())
//		//for _, e := range lineEditor.Events(ui.ly.Context) {
//		//	if _, ok := e.(controller.SubmitEvent); ok {
//		//		//topLabel = e.Text
//		//		lineEditor.SetText(f.Field.Value.(reflect.Value).String())
//		//		Info(f.Field.Value.(reflect.Value).String())
//		//	}
//		//}
//	}
// }

//func navMenuLine(gtx *layout.Context, th *gelook.DuoUItheme) func() {
//	return func() {
//		gelook.DuoUIdrawRectangle(gtx, int(navItemWidth), 1,
//		th.Colors["LightGrayIII"], [4]float32{0, 0, 0, 0}, [4]float32{0, 0, 0, 0})
//	}
//}

//
// var typeRegistry = make(map[string]reflect.Type)
//
// func makeInstance(name string) interface{} {
//	v := reflect.New(typeRegistry["cx.DuoUIconfigurationig."+name]).Elem()
//	// Maybe fill in fields here if necessary
//	return v.Interface()
// }

// func DuoUIqrCode(pubAddr string) {
//	//qr, err := qrcode.New(strings.ToUpper(pubAddr), qrcode.Medium)
//	//if err != nil {
//	//	Fatal(err)
//	//}
//	//qr.BackgroundColor = rgb(0xe8f5e9)
//	//addrQR := paint.NewImageOp(qr.Image(256))
//	return
// }

// func NewQrCode(pubAddr string) *model.DuoUIqrCode {
//	//qr, err := qrcode.New(strings.ToUpper(pubAddr), qrcode.Medium)
//	//if err != nil {
//	//	Fatal(err)
//	//}
//	//Info(pubAddr)
//	//qr.BackgroundColor = theme.HexARGB("ff3030cf")
//	//return &model.DuoUIqrCode{
//	//	AddrQR:  paint.NewImageOp(qr.Image(256)),
//	//	PubAddr: pubAddr,
//	//}
// }

//	return func() {
//		transactionsPanel := th.DuoUIPanel()
//		transactionsPanel.PanelObject = rc.History.Txs.Txs
//		transactionsPanel.ScrollBar = th.ScrollBar()
//		transactionsPanelElement.PanelObjectsNumber = len(rc.History.Txs.Txs)
//		transactionsPanel.Layout(gtx, transactionsPanelElement, func(i int, in interface{}) {
//			txs := in.([]model.DuoUItransactionExcerpt)
//			t := txs[i]
//			th.DuoUILine(gtx, 0, 0, 1, th.Colors["Hint"])()
//			for t.Link.Clicked(gtx) {
//				rc.ShowPage = fmt.Sprintf("TRANSACTION %s", t.TxID)
//				rc.GetSingleTx(t.TxID)()
//				//SetPage(rc, txPage(rc, gtx, th, t.TxID))
//			}
//			width := gtx.Constraints.Width.Max
//			button := th.DuoUIbutton("", "", "", "", "", "", "", "", 0, 0, 0, 0, 0, 0, 0, 0)
//			button.InsideLayout(gtx, t.Link, func() {
//				gtx.Constraints.Width.Min = width
//				layout.Flex{
//					Spacing: layout.SpaceBetween,
//				}.Layout(gtx,
//					layout.Rigid(txsDetails(gtx, th, i, &t)),
//					layout.Rigid(Label(gtx, th, th.Fonts["Mono"], 12, th.Colors["Secondary"], fmt.Sprintf("%0.8f", t.Amount))))
//			})
//		})
//	}
//}
