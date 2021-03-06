package votingpool

import (
	"bytes"
	"fmt"
	"github.com/p9c/pod/pkg/amt"
	"github.com/p9c/pod/pkg/btcaddr"
	"sort"
	
	"github.com/p9c/pod/pkg/chaincfg"
	"github.com/p9c/pod/pkg/txscript"
	"github.com/p9c/pod/pkg/walletdb"
	"github.com/p9c/pod/pkg/wtxmgr"
)

const eligibleInputMinConfirmations = 100

// Credit is an abstraction over wtxmgr.Credit used in the construction of voting pool withdrawal transactions.
type Credit struct {
	wtxmgr.Credit
	addr WithdrawalAddress
}

// NewCredit is
func NewCredit(c wtxmgr.Credit, addr WithdrawalAddress) Credit {
	return Credit{Credit: c, addr: addr}
}

// String is
func (c *Credit) String() string {
	return fmt.Sprintf("credit of %v locked to %v", c.Amount, c.addr)
}

// byAddress defines the methods needed to satisfy sort.Interface to txsort a slice of credits by their address.
type byAddress []Credit

func (c byAddress) Len() int { return len(c) }
func (c byAddress) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Less returns true if the element at positions i is smaller than the element at position j. The 'smaller-than'
// relation is defined to be the lexicographic ordering defined on the tuple (SeriesID, Index, Branch, TxSha,
// OutputIndex).
func (c byAddress) Less(i, j int) bool {
	iAddr := c[i].addr
	jAddr := c[j].addr
	if iAddr.seriesID < jAddr.seriesID {
		return true
	}
	if iAddr.seriesID > jAddr.seriesID {
		return false
	}
	// The seriesID are equal, so compare index.
	if iAddr.index < jAddr.index {
		return true
	}
	if iAddr.index > jAddr.index {
		return false
	}
	// The seriesID and index are equal, so compare branch.
	if iAddr.branch < jAddr.branch {
		return true
	}
	if iAddr.branch > jAddr.branch {
		return false
	}
	// The seriesID, index, and branch are equal, so compare hash.
	txidComparison := bytes.Compare(c[i].OutPoint.Hash[:], c[j].OutPoint.Hash[:])
	if txidComparison < 0 {
		return true
	}
	if txidComparison > 0 {
		return false
	}
	// The seriesID, index, branch, and hash are equal, so compare output index.
	return c[i].OutPoint.Index < c[j].OutPoint.Index
}

// getEligibleInputs returns eligible inputs with addresses between startAddress and the last used address of
// lastSeriesID. They're reverse ordered based on their address.
func (p *Pool) getEligibleInputs(
	ns, addrmgrNs walletdb.ReadBucket, store *wtxmgr.Store, txmgrNs walletdb.ReadBucket, startAddress WithdrawalAddress,
	lastSeriesID uint32, dustThreshold amt.Amount, chainHeight int32,
	minConf int,
) ([]Credit, error) {
	if p.Series(lastSeriesID) == nil {
		str := fmt.Sprintf("lastSeriesID (%d) does not exist", lastSeriesID)
		return nil, newError(ErrSeriesNotExists, str, nil)
	}
	unspents, e := store.UnspentOutputs(txmgrNs)
	if e != nil {
		return nil, newError(ErrInputSelection, "failed to get unspent outputs", e)
	}
	addrMap, e := groupCreditsByAddr(unspents, p.manager.ChainParams())
	if e != nil {
		return nil, e
	}
	var inputs []Credit
	address := startAddress
	for {
		D.C(
			func() string {
				return "looking for eligible inputs at address" + address.addrIdentifier()
			},
		)
		if candidates, ok := addrMap[address.addr.EncodeAddress()]; ok {
			var eligibles []Credit
			for _, c := range candidates {
				candidate := NewCredit(c, address)
				if p.isCreditEligible(candidate, minConf, chainHeight, dustThreshold) {
					eligibles = append(eligibles, candidate)
				}
			}
			inputs = append(inputs, eligibles...)
		}
		nAddr, e := nextAddr(p, ns, addrmgrNs, address.seriesID, address.branch, address.index, lastSeriesID+1)
		if e != nil {
			return nil, newError(ErrInputSelection, "failed to get next withdrawal address", e)
		} else if nAddr == nil {
			D.Ln("getEligibleInputs: reached last addr, stopping")
			break
		}
		address = *nAddr
	}
	sort.Sort(sort.Reverse(byAddress(inputs)))
	return inputs, nil
}

// nextAddr returns the next WithdrawalAddress according to the input selection rules:
// http://opentransactions.org/wiki/index.php/Input_Selection_Algorithm_(voting_pools) It returns nil if the new
// address' seriesID is >= stopSeriesID.
func nextAddr(
	p *Pool,
	ns, addrmgrNs walletdb.ReadBucket,
	seriesID uint32,
	branch Branch,
	index Index,
	stopSeriesID uint32,
) (
	*WithdrawalAddress, error,
) {
	series := p.Series(seriesID)
	if series == nil {
		return nil, newError(ErrSeriesNotExists, fmt.Sprintf("unknown seriesID: %d", seriesID), nil)
	}
	branch++
	if int(branch) > len(series.publicKeys) {
		highestIdx, e := p.highestUsedSeriesIndex(ns, seriesID)
		if e != nil {
			return nil, e
		}
		if index > highestIdx {
			seriesID++
			D.F(
				"nextAddr(): reached last branch (%d) and highest used index (%d), "+"moving on to next series (%d) %s",
				branch, index, seriesID,
			)
			index = 0
		} else {
			index++
		}
		branch = 0
	}
	if seriesID >= stopSeriesID {
		return nil, nil
	}
	addr, e := p.WithdrawalAddress(ns, addrmgrNs, seriesID, branch, index)
	if e != nil && e.(VPError).ErrorCode == ErrWithdrawFromUnusedAddr {
		// The used indices will vary between branches so sometimes we'll try to get a WithdrawalAddress that hasn't
		// been used before, and in such cases we just need to move on to the next one.
		D.F(
			"nextAddr(): skipping addr (series #%d, branch #%d, index #%d) "+
				"as it hasn't been used before %s", seriesID, branch, index,
		)
		return nextAddr(p, ns, addrmgrNs, seriesID, branch, index, stopSeriesID)
	}
	return addr, e
}

// highestUsedSeriesIndex returns the highest index among all of this Pool's used addresses for the given seriesID. It
// returns 0 if there are no used addresses with the given seriesID.
func (p *Pool) highestUsedSeriesIndex(ns walletdb.ReadBucket, seriesID uint32) (Index, error) {
	maxIdx := Index(0)
	series := p.Series(seriesID)
	if series == nil {
		return maxIdx,
			newError(ErrSeriesNotExists, fmt.Sprintf("unknown seriesID: %d", seriesID), nil)
	}
	for i := range series.publicKeys {
		idx, e := p.highestUsedIndexFor(ns, seriesID, Branch(i))
		if e != nil {
			return Index(0), e
		}
		if idx > maxIdx {
			maxIdx = idx
		}
	}
	return maxIdx, nil
}

// groupCreditsByAddr converts a slice of credits to a map from the string representation of an encoded address to the
// unspent outputs associated with that address.
func groupCreditsByAddr(credits []wtxmgr.Credit, chainParams *chaincfg.Params) (
	addrMap map[string][]wtxmgr.Credit, e error,
) {
	addrMap = make(map[string][]wtxmgr.Credit)
	for _, c := range credits {
		var addrs []btcaddr.Address
		_, addrs, _, e = txscript.ExtractPkScriptAddrs(c.PkScript, chainParams)
		if e != nil {
			return nil, newError(ErrInputSelection, "failed to obtain input address", e)
		}
		// As our credits are all P2SH we should never have more than one address per credit, so let's error out if that
		// assumption is violated.
		if len(addrs) != 1 {
			return nil, newError(ErrInputSelection, "input doesn't have exactly one address", nil)
		}
		encAddr := addrs[0].EncodeAddress()
		if v, ok := addrMap[encAddr]; ok {
			addrMap[encAddr] = append(v, c)
		} else {
			addrMap[encAddr] = []wtxmgr.Credit{c}
		}
	}
	return addrMap, nil
}

// isCreditEligible tests a given credit for eligibility with respect to number of confirmations, the dust threshold and
// that it is not the charter output.
func (p *Pool) isCreditEligible(
	c Credit, minConf int, chainHeight int32,
	dustThreshold amt.Amount,
) bool {
	if c.Amount < dustThreshold {
		return false
	}
	if confirms(c.BlockMeta.Block.Height, chainHeight) < int32(minConf) {
		return false
	}
	if p.isCharterOutput(c) {
		return false
	}
	return true
}

// isCharterOutput -
//
// TODO: In order to determine this, we need the txid and the output index of the current charter output, which we don't
//  have yet.
func (p *Pool) isCharterOutput(c Credit) bool {
	return false
}

// confirms returns the number of confirmations for a transaction in a block at height txHeight (or -1 for an
// unconfirmed tx) given the chain height curHeight.
func confirms(txHeight, curHeight int32) int32 {
	switch {
	case txHeight == -1, txHeight > curHeight:
		return 0
	default:
		return curHeight - txHeight + 1
	}
}
