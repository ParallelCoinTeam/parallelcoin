package controller_test

import (
	"encoding/binary"
	"encoding/hex"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/controller"
	"net"
	"testing"
)

func TestPort(t *testing.T) {
	var portNumber uint16 = 11047
	port := controller.NewPort()
	port.Put(portNumber)
	port2 := controller.NewPort()
	port2.Decode(port.Encode())
	if port2.Get() != port.Get() {
		t.Fail()
	}
}

func TestIP(t *testing.T) {
	var ipa = net.ParseIP("127.0.0.1")
	ip := controller.NewIP()
	ip.Put(&ipa)
	t.Log(ip.Get().MarshalText())
	ip2 := controller.NewIP()
	ip2.Decode(ip.Encode())
	if ip.Get().Equal(*ip2.Get()) {
		t.Fail()
	}
}

func TestIPs(t *testing.T) {
	var ipa1 = net.ParseIP("127.0.0.1")
	var ipa2 = net.ParseIP("fe80::6382:2df5:7014:e156")
	ips := controller.NewIPs()
	ips.Put([]*net.IP{&ipa1, &ipa2})
	ips2 := controller.NewIPs()
	ips2.Decode(ips.Encode())
	dec := ips.Get()
	dec2 := ips2.Get()
	for i := range dec {
		if !dec[i].Equal(*dec2[i]) {
			t.Fail()
		}
	}
}

func TestInt32(t *testing.T) {
	by, err := hex.DecodeString("deadbeef")
	if err != nil {
		panic(err)
	}
	bits := binary.BigEndian.Uint32(by)
	bt := controller.NewInt32()
	bt.Put(int32(bits))
	bt2 := controller.NewInt32()
	bt2.Decode(bt.Encode())
	if bt.Get() != bt2.Get() {
		t.Fail()
	}
}

func TestHash(t *testing.T) {
	by, err := hex.DecodeString(
		"00c44981699c4b621fe89b32148a64fc11fb680fa484ab1abe0e6fba4fcca462")
	var bhash chainhash.Hash
	err = bhash.SetBytes(by)
	if err != nil {
		panic(err)
	}
	h := controller.NewHash()
	h.Put(bhash)
	h2 := controller.NewHash()
	h2.Decode(h.Encode())
	if !h.Get().IsEqual(h2.Get()) {
		t.Fail()
	}
}

func TestTransaction(t *testing.T) {
	//txI := wire.NewMsgTx(wire.TxVersion)
	//txx :=
}

// standardCoinbaseScript returns a standard script suitable for use as the
// signature script of the coinbase transaction of a new block.  In particular,
// it starts with the block height that is required by version 2 blocks and
// adds the extra nonce as well as additional coinbase flags.
func standardCoinbaseScript(nextBlockHeight int32, extraNonce uint64) ([]byte, error) {
	return txscript.NewScriptBuilder().AddInt64(int64(nextBlockHeight)).
		AddInt64(int64(extraNonce)).AddData([]byte(mining.CoinbaseFlags)).
		Script()
}
