package wallet

import (
	"github.com/p9c/matrjoska/pkg/btcaddr"
	"github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/walletdb"
	"github.com/p9c/matrjoska/pkg/wire"
)

// OutputSelectionPolicy describes the rules for selecting an output from the wallet.
type OutputSelectionPolicy struct {
	Account               uint32
	RequiredConfirmations int32
}

func (p *OutputSelectionPolicy) meetsRequiredConfs(txHeight, curHeight int32) bool {
	return confirmed(p.RequiredConfirmations, txHeight, curHeight)
}

// UnspentOutputs fetches all unspent outputs from the wallet that match rules described in the passed policy.
func (w *Wallet) UnspentOutputs(policy OutputSelectionPolicy) ([]*TransactionOutput, error) {
	var outputResults []*TransactionOutput
	e := walletdb.View(
		w.db, func(tx walletdb.ReadTx) (e error) {
			addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
			txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
			syncBlock := w.Manager.SyncedTo()
			// TODO: actually stream outputs from the db instead of fetching all of them at once.
			outputs, e := w.TxStore.UnspentOutputs(txmgrNs)
			if e != nil {
				return e
			}
			for _, output := range outputs {
				// Ignore outputs that haven't reached the required number of confirmations.
				if !policy.meetsRequiredConfs(output.Height, syncBlock.Height) {
					continue
				}
				// Ignore outputs that are not controlled by the account.
				var addrs []btcaddr.Address
				_, addrs, _, e = txscript.ExtractPkScriptAddrs(
					output.PkScript,
					w.chainParams,
				)
				if e != nil || len(addrs) == 0 {
					// Cannot determine which account this belongs to without a valid address.
					//
					// TODO: Fix this by saving outputs per account, or accounts per output.
					continue
				}
				var outputAcct uint32
				_, outputAcct, e = w.Manager.AddrAccount(addrmgrNs, addrs[0])
				if e != nil {
					return e
				}
				if outputAcct != policy.Account {
					continue
				}
				// Stakebase isn't exposed by wtxmgr so those will be OutputKindNormal for now.
				outputSource := OutputKindNormal
				if output.FromCoinBase {
					outputSource = OutputKindCoinbase
				}
				result := &TransactionOutput{
					OutPoint: output.OutPoint,
					Output: wire.TxOut{
						Value:    int64(output.Amount),
						PkScript: output.PkScript,
					},
					OutputKind:      outputSource,
					ContainingBlock: BlockIdentity(output.Block),
					ReceiveTime:     output.Received,
				}
				outputResults = append(outputResults, result)
			}
			return nil
		},
	)
	return outputResults, e
}
