package votingpool

import (
	"github.com/parallelcointeam/parallelcoin/pkg/chain/wire"
	"github.com/parallelcointeam/parallelcoin/pkg/util/hdkeychain"
	waddrmgr "github.com/parallelcointeam/parallelcoin/pkg/wallet/addrmgr"
	walletdb "github.com/parallelcointeam/parallelcoin/pkg/wallet/db"
)

var TstLastErr = lastErr

// const TstEligibleInputMinConfirmations = eligibleInputMinConfirmations

// TstPutSeries transparently wraps the voting pool putSeries method.
func (p *Pool) TstPutSeries(ns walletdb.ReadWriteBucket, version, seriesID, reqSigs uint32, inRawPubKeys []string) error {
	return p.putSeries(ns, version, seriesID, reqSigs, inRawPubKeys)
}

var TstBranchOrder = branchOrder

// TstExistsSeries checks whether a series is stored in the database.
func (p *Pool) TstExistsSeries(dbtx walletdb.ReadTx, seriesID uint32) (bool, error) {
	ns, _ := TstRNamespaces(dbtx)
	poolBucket := ns.NestedReadBucket(p.ID)
	if poolBucket == nil {
		return false, nil
	}
	bucket := poolBucket.NestedReadBucket(seriesBucketName)
	if bucket == nil {
		return false, nil
	}
	return bucket.Get(uint32ToBytes(seriesID)) != nil, nil
}

// TstGetRawPublicKeys gets a series public keys in string format.
func (s *SeriesData) TstGetRawPublicKeys() []string {
	rawKeys := make([]string, len(s.publicKeys))
	for i, key := range s.publicKeys {
		rawKeys[i] = key.String()
	}
	return rawKeys
}

// TstGetRawPrivateKeys gets a series private keys in string format.
func (s *SeriesData) TstGetRawPrivateKeys() []string {
	rawKeys := make([]string, len(s.privateKeys))
	for i, key := range s.privateKeys {
		if key != nil {
			rawKeys[i] = key.String()
		}
	}
	return rawKeys
}

// TstGetReqSigs expose the series reqSigs attribute.
func (s *SeriesData) TstGetReqSigs() uint32 {
	return s.reqSigs
}

// TstEmptySeriesLookup empties the voting pool seriesLookup attribute.
func (p *Pool) TstEmptySeriesLookup() {
	p.seriesLookup = make(map[uint32]*SeriesData)
}

// TstDecryptExtendedKey expose the decryptExtendedKey method.
func (p *Pool) TstDecryptExtendedKey(keyType waddrmgr.CryptoKeyType, encrypted []byte) (*hdkeychain.ExtendedKey, error) {
	return p.decryptExtendedKey(keyType, encrypted)
}

// TstGetMsgTx returns a copy of the withdrawal transaction with the given
// ntxid.
func (s *WithdrawalStatus) TstGetMsgTx(ntxid Ntxid) *wire.MsgTx {
	return s.transactions[ntxid].MsgTx.Copy()
}
