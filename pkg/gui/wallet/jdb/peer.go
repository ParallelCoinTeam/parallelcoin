package jdb

import (
	"context"
	"github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	provider "github.com/ipfs/go-ipfs-provider"
	ipld "github.com/ipfs/go-ipld-format"
	"github.com/libp2p/go-libp2p-core/crypto"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/routing"
	"github.com/multiformats/go-multiaddr"
)

// Peer is an IPFS-Lite peer. It provides a DAG service that can fetch and put
// blocks from/to the IPFS network.
type Peer struct {
	ctx context.Context

	cfg *Config

	host  host.Host
	dht   routing.Routing
	store datastore.Batching

	ipld.DAGService // become a DAG service
	bstore          blockstore.Blockstore
	bserv           blockservice.BlockService
	reprovider      provider.System
}

func GetPeer(ctx context.Context, ds datastore.Batching) *Peer {
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}
	listen, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")
	//ob, _ := pub.Bytes()
	//pp := hex.EncodeToString(ob)
	//fmt.Println("Pub", pp)
	h, dht, err := SetupLibp2p(
		ctx,
		priv,
		nil,
		[]multiaddr.Multiaddr{listen},
		ds,
		Libp2pOptionsExtra...,
	)
	if err != nil {
		panic(err)
	}
	p, err := NewPeer(ctx, ds, h, dht, nil)
	if err != nil {
		panic(err)
	}
	p.Bootstrap(DefaultBootstrapPeers())
	return p
}

// New creates an IPFS-Lite Peer. It uses the given datastore, libp2p Host and
// Routing (usuall the DHT). The Host and the Routing may be nil if
// config.Offline is set to true, as they are not used in that case. Peer
// implements the ipld.DAGService interface.
func NewPeer(ctx context.Context, store datastore.Batching, host host.Host, dht routing.Routing, cfg *Config) (*Peer, error) {
	if cfg == nil {
		cfg = &Config{}
	}
	cfg.setDefaults()

	p := &Peer{
		ctx:   ctx,
		cfg:   cfg,
		host:  host,
		dht:   dht,
		store: store,
	}

	err := p.setupBlockstore()
	if err != nil {
		return nil, err
	}
	err = p.setupBlockService()
	if err != nil {
		return nil, err
	}
	err = p.setupDAGService()
	if err != nil {
		p.bserv.Close()
		return nil, err
	}
	err = p.setupReprovider()
	if err != nil {
		p.bserv.Close()
		return nil, err
	}

	go p.autoclose()

	return p, nil
}
