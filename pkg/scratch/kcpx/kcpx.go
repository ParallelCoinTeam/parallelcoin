// Package kcpx is some tidy wrappers around rpcx/kcp RPC communication
package kcpx

import (
	"context"
	"crypto/sha1"
	"net"

	"github.com/p9c/kcp"
	"github.com/p9c/rpcx/client"
	"github.com/p9c/rpcx/server"
	"golang.org/x/crypto/pbkdf2"

	"github.com/p9c/pod/pkg/log"
)

// NewXClient returns a new XClient configured for parallelcoin rpcx service
// for a given service name
func NewXClient(address, serviceName, password string) (x client.XClient) {
	serviceDiscovery := client.NewPeer2PeerDiscovery("kcp@"+address, "")
	passwordBytes := pbkdf2.Key(reverse([]byte(password)), []byte(password),
		4096, 32,
		sha1.New)
	bc, _ := kcp.NewAESBlockCrypt(passwordBytes)
	option := client.DefaultOption
	option.Block = bc
	x = client.NewXClient(serviceName, client.Failtry,
		client.RoundRobin, serviceDiscovery, option)
	cs := &ConfigUDPSession{}
	pc := client.NewPluginContainer()
	pc.Add(cs)
	x.SetPlugins(pc)
	return
}

// Serve serves up an RPC service that can be contacted over kcp with the
// same password as used for NewXClient
func Serve(address, serviceName, password string,
	service interface{}) (srv *server.Server, shutdown func() <-chan struct{}) {
	passwordBytes := pbkdf2.Key(reverse([]byte(password)), []byte(password),
		4096, 32,
		sha1.New)
	bc, _ := kcp.NewAESBlockCrypt(passwordBytes)
	srv = server.NewServer(server.WithBlockCrypt(bc))
	err := srv.RegisterName(serviceName, service, "")
	if err != nil {
		log.ERROR("error registering interface ", serviceName, " ", err)
		return
	}
	cs := &ConfigUDPSession{}
	srv.Plugins.Add(cs)
	ctx := context.Background()
	shutdown = func() <-chan struct{} {
		err := srv.Shutdown(ctx)
		if err != nil {
			log.ERROR("error shutting down server ", err)
		}
		return ctx.Done()
	}
	err = srv.Serve("kcp", address)
	if err != nil {
		log.ERROR("error serving ", serviceName, " ", err)
	}
	return
}

type ConfigUDPSession struct{}

func (p *ConfigUDPSession) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	session, ok := conn.(*kcp.UDPSession)
	if !ok {
		return conn, true
	}

	session.SetACKNoDelay(true)
	session.SetStreamMode(true)
	return conn, true
}

func reverse(b []byte) []byte {
	out := make([]byte, len(b))
	for i := range b {
		out[i] = b[len(b)-1]
	}
	return out
}
