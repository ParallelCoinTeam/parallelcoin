package gwallet

import (
	"gioui.org/text"
	"github.com/gioapp/gel/helper"
	"github.com/p9c/pod/pkg/gui/wallet/dap/box"
	"github.com/p9c/pod/pkg/gui/wallet/lyt"
	"github.com/p9c/pod/pkg/gui/wallet/theme"
)

func (g *GioWallet) GetExplore() {
	//f, err := g.sh.ID()
	//checkError(err)
	//fv, ft, err := g.sh.Version()
	//checkError(err)
	//g.Status = Status{
	//	Title: "string",
	//	HostingSize: "uint",
	//PeerId:  f.ID,
	//Version: fv + " " + ft,
	//Gateway:   ,
	//Api:       "string",
	//Addresses: f.Addresses,
	//Pub:       f.PublicKey,
	//}
}

func (g *GioWallet) exploreHeader() func(gtx C) D {
	return func(gtx C) D {
		//return ContainerLayout(th.Colors["PanelBg"], 10, 10, 10, 10, func(gtx C) D {
		//
		//})
		return D{}
	}
}

func (g *GioWallet) exploreBody() func(gtx C) D {
	//return g.twoPanels(10, 0, g.exploreLeft(), g.exploreRight())
	return noReturn
}

//func (g *GioWallet) exploreBody() []func(gtx C) D {
//	retu//		ContainerLayout(th.Colors["Silver"], 1, 1, 1, 1, ContainerLayout(th.Colors["PanelBg"], 30, 30, 30, 30, func(gtx C) D {
//			gtx.Constraints.Min.X = gtx.Constraints.Max.X
//			return lyt.Format(gtx, "vflex(middle,f(0.6,inset(5dp0dp5dp0dp,_)),f(0.4, inset(5dp0dp5dp0dp,_))))",
//				func(gtx C) D {
//					gtx.Constraints.Max.Y = gtx.Constraints.Min.Y
//					return g.6UI.logoIpld.Layout(gtx)
//				},
//				func(gtx C) D {
//					title := theme.Body(g.ui.Theme, "IPLD is the data model of the content-addressable web. It allows us to treat all hash-linked data structures as subsets of a unified information space, unifying all data models that link data with hashes as instances of IPLD.\n\nContent addressing through hashes has become a widely-used means of connecting data in distributed systems, from the blockchains that run your favorite cryptocurrencies, to the commits that back your code, to the web’s content at large. Yet, whilst all of these tools rely on some common primitives, their specific underlying data structures are not interoperable.\n\nEnter IPLD: a single namespace for all hash-inspired protocols. Through IPLD, links can be traversed across protocols, allowing you to explore data regardless of the underlying protocol.")
//					title.Alignment = text.Start
//					return title.Layout(gtx)
//				})
//		})))

//}
//}
func (g *GioWallet) exploreLeft() func(gtx C) D {
	return func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return lyt.Format(gtx, "vflex(start,r(inset(0dp0dp30dp0dp,_)),r(inset(0dp0dp80dp0dp,_))))",
			func(gtx C) D {
				title := theme.H3(g.ui.Theme, "Explore the blockchain")
				title.Alignment = text.Start
				title.Color = helper.HexARGB(g.ui.Theme.Colors["Success"])
				return title.Layout(gtx)
			},
			func(gtx C) D {
				title := theme.Body(g.ui.Theme, "Paste a address the box to fetch.")
				title.Alignment = text.Start
				return title.Layout(gtx)
			})
	}
}

func (g *GioWallet) exploreRight() func(gtx C) D {
	return box.BoxPanel(g.ui.Theme, func(gtx C) D {

		//return container.C().OutsideColor(th.Colors["PanelBg"],).Layout(func(gtx C) D {
		gtx.Constraints.Min.X = gtx.Constraints.Max.X
		return lyt.Format(gtx, "vflex(start,r(inset(0dp0dp0dp0dp,_)),r(inset(0dp0dp0dp0dp,_)))",
			func(gtx C) D {
				title := theme.H3(g.ui.Theme, "What is ParallelCoin?")
				title.Alignment = text.Start
				return title.Layout(gtx)
			},
			func(gtx C) D {
				title := theme.Body(g.ui.Theme, "We are very proud of the work so far and we think that people will love the new Parallelcoin: \n9 complex multi-algorithm orthogonal complexity proof of work hash functions\n4 part averaging algorithm\nstrong resistance to timing and rhythm attacks by the use of multiple competing averagers\ndefault odds-following weighted randomising work schedulers\n9 second blocks with difficulty reduction damper to reduce coincidence of very short intervals between blocks.")
				title.Alignment = text.Start
				return title.Layout(gtx)
			})
	})
}
