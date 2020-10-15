package jdb

// This example launches an IPFS-Lite peer and fetches a hello-world
// hash from the IPFS network.

import (
	"context"
	"github.com/ipfs/go-cid"

	format "github.com/ipfs/go-ipld-format"
	crypto "github.com/libp2p/go-libp2p-core/crypto"
)

// JavazacDb Structure
type JavazacDB struct {
	ctx       context.Context
	peer      *Peer
	index     map[string]string
	Datastore string
}

func New(datastore string) *JavazacDB {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	crypto.MinRsaKeyBits = 1024
	ds, err := BadgerDatastore(datastore)
	if err != nil {
		panic(err)
	}
	peer := GetPeer(ctx, ds)
	return &JavazacDB{
		ctx:   ctx,
		peer:  peer,
		index: make(map[string]string),
	}
}
func (j *JavazacDB) ReadList(c cid.Cid) []*format.Link {
	rsc, err := j.peer.Get(j.ctx, c)
	checkError(err)
	return rsc.Links()
}
