package explorer

import (
	"fmt"

	l "gioui.org/layout"

	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/gui/p9"
)

func (ex *Explorer) Blocks() l.Widget {
	listPageSize := 10

	fmt.Println("Best Block Height:", ex.State.bestBlockHeight)
	le := func(gtx l.Context, index int) l.Dimensions {
		b := l.Dimensions{}
		if ex.State.bestBlockHeight > 0 {
			blockHash, err := ex.ChainClient.GetBlockHash(int64(ex.State.bestBlockHeight - index))
			if err != nil {
			}
			block, err := ex.ChainClient.GetBlock(blockHash)
			if err != nil {
			}
			fmt.Println("Block:", block)
			b = ex.singleBlock(gtx, block)
		}
		return b
	}
	return func(gtx l.Context) l.Dimensions {
		return ex.th.Responsive(*ex.App.Size, p9.Widgets{
			{
				Widget: ex.th.VFlex().
					Flexed(1,
						// ex.Inset(0.25,
						func(gtx l.Context) l.Dimensions {
							return ex.lists["blocks"].Vertical().Length(listPageSize).ListElement(le).Fn(gtx)
							// ).
							// Rigid(
							// ex.sendFooter(),
							// ).Fn
							//		).Fn},
						}).Fn,
			},
		}).Fn(gtx)
	}
}

func (ex *Explorer) singleBlock(gtx l.Context, block *wire.MsgBlock) l.Dimensions {
	return ex.th.Inset(0.25,
		ex.th.Fill("DocBg", ex.th.Inset(0.25,
			ex.th.VFlex().
				Rigid(
					ex.th.Inset(0.2,
						ex.th.Caption(block.Header.BlockHash().String()).
							Color("PanelText").
							Fn,
					).Fn,
				).
				Rigid(
					ex.th.Inset(0.2,
						ex.th.Caption(block.Header.Timestamp.Format("")).
							Color("PanelText").
							Fn,
					).Fn,
				).Fn,
		).Fn, l.Center).Fn,
	).Fn(gtx)
}
