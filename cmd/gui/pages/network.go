package pages

import (
	"fmt"
	"gioui.org/layout"
	"gioui.org/unit"
	"github.com/p9c/pod/cmd/gui/component"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gelook"
)

var (
//itemValue = &gel.DuoUIcounter{
//	Value:        11,
//	OperateValue: 1,
//	From:         0,
//	To:           15,
//	CounterInput: &gel.Editor{
//		Alignment:  text.Middle,
//		SingleLine: true,
//	},
//	CounterIncrease: new(gel.Button),
//	//CounterDecrease: new(controller.Button),
//	CounterReset: new(gel.Button),
//}
)

func Network(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) *gelook.DuoUIpage {
	page := gelook.DuoUIpage{
		Title:         "NETWORK",
		TxColor:       "",
		Command:       rc.GetPeerInfo(),
		Border:        4,
		Header:        networkHeader(rc, gtx, th),
		HeaderBgColor: "",
		HeaderPadding: 0,
		Body:          networkBody(rc, gtx, th),
		BodyBgColor:   "",
		BodyPadding:   0,
		Footer:        func() {},
		FooterBgColor: "",
		FooterPadding: 0,
	}
	return th.DuoUIpage(page)
}

func networkHeader(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.Flex{
			Spacing: layout.SpaceBetween,
		}.Layout(gtx,
			//layout.Rigid(component.TransactionsFilter(rc, gtx, th)),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetPeerInfo()).Layout(gtx, rc.Network.PerPage, "Peers per page: ", fmt.Sprint(rc.Network.PerPage.Value))
			}),
			layout.Rigid(func() {
				th.DuoUIcounter(rc.GetPeerInfo()).Layout(gtx, rc.Network.Page, "Peers page: ", fmt.Sprint(rc.Network.Page.Value))
			}),
		)
	}
}

func networkBody(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		layout.UniformInset(unit.Dp(16)).Layout(gtx, func() {
			layout.Flex{
				Axis: layout.Vertical,
			}.Layout(gtx,
				layout.Rigid(component.PeersList(rc, gtx, th)))
		})
	}
}
