package component

import (
	"gioui.org/layout"
	"github.com/p9c/pod/cmd/gui/rcd"
	"github.com/p9c/pod/pkg/gelook"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

var ()

func PeersList(rc *rcd.RcVar, gtx *layout.Context, th *gelook.DuoUItheme) func() {
	return func() {
		rc.Network.PeersList.Layout(gtx, len(rc.Network.Peers), func(i int) {
			t := rc.Network.Peers[i]
			HorizontalLine(gtx, 1, th.Colors["Hint"])()
			layout.Flex{
				Spacing: layout.SpaceBetween,
			}.Layout(gtx,
				layout.Rigid(peerDetails(gtx, th, i, t)),
				layout.Rigid(Label(gtx, th, th.Fonts["Mono"], 14, th.Colors["Dark"], t.Addr)))
		})
	}
}

func peerDetails(gtx *layout.Context, th *gelook.DuoUItheme, i int, t *btcjson.GetPeerInfoResult) func() {
	return func() {
		//layout.Flex{
		//	Axis: layout.Horizontal,
		//}.Layout(gtx,
		//	layout.Rigid(func(){th.DuoUIitem(th.Colors["Hint"]).Layout(gtx, Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(i)))}),
		//	layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.ID))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.AddrLocal)),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.Services)),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.RelayTxes))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.LastSend))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.LastRecv))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.BytesSent))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.BytesRecv))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.ConnTime))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.TimeOffset))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.PingTime))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.PingWait))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.Version))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], t.SubVer)),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.Inbound))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.StartingHeight))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.CurrentHeight))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.BanScore))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.FeeFilter))),
		//layout.Rigid(Label(gtx, th, th.Fonts["Primary"], 12, th.Colors["Dark"], fmt.Sprint(t.SyncNode))))
	}
}
