package wallet

import (
	"github.com/p9c/pkg/app/slog"

	wtxmgr "github.com/p9c/pod/pkg/chain/tx/mgr"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/db/walletdb"
	"github.com/p9c/pod/pkg/util"
)

// OutputSelectionPolicy describes the rules for selecting an output from the
// wallet.
type OutputSelectionPolicy struct {
	Account               uint32
	RequiredConfirmations int32
}

func (p *OutputSelectionPolicy) meetsRequiredConfs(txHeight, curHeight int32) bool {
	return confirmed(p.RequiredConfirmations, txHeight, curHeight)
}

// UnspentOutputs fetches all unspent outputs from the wallet that match rules
// described in the passed policy.
func (w *Wallet) UnspentOutputs(policy OutputSelectionPolicy) (outputResults []*TransactionOutput, err error) {
	if err = walletdb.View(w.db, func(tx walletdb.ReadTx) (err error) {
		addrmgrNs := tx.ReadBucket(waddrmgrNamespaceKey)
		txmgrNs := tx.ReadBucket(wtxmgrNamespaceKey)
		syncBlock := w.Manager.SyncedTo()
		// TODO: actually stream outputs from the db instead of fetching  all of them at once.
		var outputs []wtxmgr.Credit
		if outputs, err = w.TxStore.UnspentOutputs(txmgrNs); slog.Check(err) {
			return
		}
		var addrs []util.Address
		for _, output := range outputs {
			// Ignore outputs that haven't reached the required number of confirmations.
			if !policy.meetsRequiredConfs(output.Height, syncBlock.Height) {
				continue
			}
			// Ignore outputs that are not controlled by the account.
			if _, addrs, _, err = txscript.ExtractPkScriptAddrs(output.PkScript, w.chainParams); slog.Check(err) || len(addrs) == 0 {
				// Cannot determine which account this belongs to without a valid address.
				// TODO: Fix this by saving outputs per account, or accounts per output.
				continue
			}
			var outputAcct uint32
			if _, outputAcct, err = w.Manager.AddrAccount(addrmgrNs, addrs[0]); slog.Check(err) {
				return
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
		return
	}); slog.Check(err) {
	}
	return
}
