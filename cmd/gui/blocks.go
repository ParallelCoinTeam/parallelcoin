package gui

import (
	"fmt"
	l "gioui.org/layout"
	"github.com/p9c/pod/pkg/gui/p9"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

type block struct {
	time      string
	data      *btcjson.GetBlockVerboseResult
	clickPrev *p9.Clickable
	clickNext *p9.Clickable
	list      *p9.List
}

func (wg *WalletGUI) blockIitem(label, data string) l.Widget {
	if data != "" {
		return wg.Inset(0.25,
			wg.th.VFlex().
				Rigid(
					wg.Inset(0.0, wg.Fill("PanelBg", wg.Inset(0.2, wg.H6(label).Color("DocText").Fn).Fn).Fn).Fn,
				).
				Rigid(
					wg.Inset(0.0, wg.Fill("DocBg", wg.Inset(0.2, wg.Body1(data).Color("DocText").Font("go regular").Fn).Fn).Fn).Fn,
				).Fn,
		).Fn
	} else {
		return p9.EmptyMaxWidth()
	}
}

func (wg *WalletGUI) blockPage(blockHeight int64) func() {
	b := wg.getBlock(blockHeight)
	blockLayout := []l.Widget{
		//wg.blockIitem("Block Height:", fmt.Sprint(b.data.Height)),
		wg.blockIitem("Hash:", fmt.Sprint(blockHeight)),
		//wg.blockIitem("Confirmations:", fmt.Sprint(b.data.Confirmations)),
		//wg.blockIitem("Stripped Size:", fmt.Sprint(b.data.StrippedSize)),
		//wg.blockIitem("Size:", fmt.Sprint(b.data.Size)),
		//wg.blockIitem("Weight:", fmt.Sprint(b.data.Weight)),
		//wg.blockIitem("Height:", fmt.Sprint(b.data.Height)),
		//wg.blockIitem("Version:", fmt.Sprint(b.data.Version)),
		//wg.blockIitem("Version Hex:", fmt.Sprint(b.data.VersionHex)),
		//wg.blockIitem("Pow Algo ID:", fmt.Sprint(b.data.PowAlgoID)),
		//wg.blockIitem("Pow Algo:", fmt.Sprint(b.data.PowAlgo)),
		//wg.blockIitem("Pow Hash:", fmt.Sprint(b.data.PowHash)),
		//wg.blockIitem("Merkle Root:", fmt.Sprint(b.data.MerkleRoot)),
		//wg.blockIitem("Transactions Number:", fmt.Sprint(b.data.TxNum)),
		//wg.blockIitem("Transaction:", fmt.Sprint(b.data.Tx)),
		//wg.blockIitem("Raw Transaction:", fmt.Sprint(b.data.RawTx)),
		//wg.blockIitem("Time:", fmt.Sprint(b.data.Time)),
		//wg.blockIitem("Nonce:", fmt.Sprint(b.data.Nonce)),
		//wg.blockIitem("Bits:", fmt.Sprint(b.data.Bits)),
		//wg.blockIitem("Difficulty:", fmt.Sprint(b.data.Difficulty)),
		//wg.blockIitem("Previous Hash:", fmt.Sprint(b.data.PreviousHash)),
		//wg.blockIitem("Next Hash:", fmt.Sprint(b.data.NextHash)),
	}
	le := func(gtx l.Context, index int) l.Dimensions {
		return blockLayout[index](gtx)
	}

	return func() {
		wg.newWindow(b.data.Hash, "Block: "+fmt.Sprint(b.data.Height), 600, 800,
			wg.th.VFlex().
				Rigid(
					wg.Inset(0.0, wg.Fill("Primary", wg.Inset(0.5, wg.Caption("Block "+fmt.Sprint(blockHeight)).Color("DocBg").Fn).Fn).Fn).Fn,
				).
				Flexed(1,
					wg.Inset(0,
						func(gtx l.Context) l.Dimensions {
							return b.list.Vertical().Length(len(blockLayout)).ListElement(le).Fn(gtx)
						},
					).Fn,
				).
				Rigid(
					wg.th.Flex().
						Flexed(0.5,
							wg.Button(
								b.clickPrev.SetClick(func() {
									//wg.w[wg.State.txs[i].data.TxID].Window.Close()
								})).
								CornerRadius(0).
								Background("Primary").
								Color("Dark").
								Font("bariol bold").
								TextScale(1).
								Text("< previous block").
								Inset(0.5).
								Fn,
						).
						Flexed(0.5,
							wg.Button(
								b.clickNext.SetClick(func() {
									//wg.w[wg.State.txs[i].data.TxID].Window.Close()
								})).
								CornerRadius(0).
								Background("Primary").
								Color("Dark").
								Font("bariol bold").
								TextScale(1).
								Text("next block >").
								Inset(0.5).
								Fn,
						).Fn,
				).Fn,
		)
	}
}

func (wg *WalletGUI) getBlock(blockHeight int64) *block {
	chainClient, err := wg.chainClient()
	if err != nil {
	}
	blockHash, err := chainClient.GetBlockHash(blockHeight)
	if err != nil {
	}
	data, err := chainClient.GetBlockVerbose(blockHash)
	if err != nil {
	}
	fmt.Println("dadad", data)
	return &block{
		data:      data,
		clickPrev: wg.th.Clickable(),
		clickNext: wg.th.Clickable(),
		list:      wg.th.List(),
	}
}
