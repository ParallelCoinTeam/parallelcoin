package control

import (
	"fmt"
	"github.com/VividCortex/ewma"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/cmd/kopach/control/job"
	"github.com/p9c/pod/cmd/kopach/control/p2padvt"
	"github.com/p9c/pod/cmd/walletmain"
	blockchain "github.com/p9c/pod/pkg/chain"
	chainhash "github.com/p9c/pod/pkg/chain/hash"
	"github.com/p9c/pod/pkg/comm/transport"
	rpcclient "github.com/p9c/pod/pkg/rpc/client"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/routeable"
	"github.com/urfave/cli"
	"net"
	"time"
)

func (c *Controller) walletRPCWatcher() {
	Debug("starting wallet rpc connection watcher for mining addresses")
	var err error
	backoffTime := time.Second
	certs := walletmain.ReadCAFile(c.cx.Config)
totalOut:
	for {
	trying:
		for {
			select {
			case <-c.cx.KillAll.Wait():
				break totalOut
			default:
			}
			Debug("trying to connect to wallet for mining addresses...")
			// If we can reach the wallet configured in the same datadir we can mine
			if c.walletClient, err = rpcclient.New(
				&rpcclient.ConnConfig{
					Host:         *c.cx.Config.WalletServer,
					Endpoint:     "ws",
					User:         *c.cx.Config.Username,
					Pass:         *c.cx.Config.Password,
					TLS:          *c.cx.Config.TLS,
					Certificates: certs,
				}, nil, c.cx.KillAll,
			); Check(err) {
				Debug("failed, will try again")
				c.isMining.Store(false)
				select {
				case <-time.After(backoffTime):
				case <-c.quit.Wait():
					c.isMining.Store(false)
					break totalOut
				}
				if backoffTime <= time.Second*5 {
					backoffTime += time.Second
				}
				continue
			} else {
				Debug("<<<controller has wallet connection>>>")
				c.isMining.Store(true)
				backoffTime = time.Second
				break trying
			}
		}
		Debug("<<<connected to wallet>>>")
		retryTicker := time.NewTicker(time.Second)
	connected:
		for {
			select {
			case <-retryTicker.C:
				if c.walletClient.Disconnected() {
					c.isMining.Store(false)
					break connected
				}
			case <-c.quit.Wait():
				c.isMining.Store(false)
				break totalOut
			}
		}
		Debug("disconnected from wallet")
	}
}

func (c *Controller) advertiserAndRebroadcaster() {
	if !c.active.Load() {
		Info("ready to send out jobs!")
		c.active.Store(true)
	}
	factor := 2
	ticker := time.NewTicker(time.Second * time.Duration(factor))
	const countTick = 10
	counter := countTick / 2
	once := false
	var err error
out:
	for {
		select {
		case <-ticker.C:
			c.height.Store(c.cx.RPCServer.Cfg.Chain.BestSnapshot().Height)
			if c.isMining.Load() {
				if !once {
					c.cx.RealNode.Chain.Subscribe(c.chainNotifier())
					once = true
					c.active.Store(true)
				}
			}
			if counter%countTick == 0 {
				j := p2padvt.GetAdvt(c.cx)
				if *c.cx.Config.AutoListen {
					*c.cx.Config.P2PConnect = cli.StringSlice{}
					_, addresses := routeable.GetAllInterfacesAndAddresses()
					Traces(addresses)
					for i := range addresses {
						addrS := net.JoinHostPort(addresses[i].IP.String(), fmt.Sprint(j.P2P))
						*c.cx.Config.P2PConnect = append(*c.cx.Config.P2PConnect, addrS)
					}
					save.Pod(c.cx.Config)
				}
			}
			counter++
			// send out advertisment
			if err = c.multiConn.SendMany(p2padvt.Magic, transport.GetShards(p2padvt.Get(c.cx))); Check(err) {
			}
			if c.isMining.Load() {
				if err = c.updateAndSendWork(); Check(err) {
				}
			}
		case <-c.quit.Wait():
			Debug("quitting on close quit channel")
			break out
		case <-c.cx.NodeKill.Wait():
			Debug("quitting on NodeKill")
			c.quit.Q()
			break out
		case <-c.cx.KillAll.Wait():
			Debug("quitting on KillAll")
			c.quit.Q()
			break out
		}
	}
	c.active.Store(false)
	Debug("controller exiting")
}

func (c *Controller) chainNotifier() func(n *blockchain.Notification) {
	return func(n *blockchain.Notification) {
		// First to arrive locks out any others while processing
		switch n.Type {
		case blockchain.NTBlockConnected:
			Trace("received new chain notification")
			// construct work message
			if _, ok := n.Data.(*util.Block); !ok {
				Warn("chain accepted notification is not a block")
				break
			}
			if err := c.updateAndSendWork(); Check(err) {
			}
		}
	}
}

func (c *Controller) updateAndSendWork() (err error) {
	var getNew bool
	// The current block is stale if the best block has changed.
	best := c.blockTemplateGenerator.BestSnapshot()
	oB, ok := c.oldBlocks.Load().([][]byte)
	switch {
	case len(oB) == 0:
		Trace("cached template is zero length")
		getNew = true
		fallthrough
	case !ok:
		Trace("cached template is nil")
		getNew = true
		fallthrough
	case !c.prevHash.Load().(*chainhash.Hash).IsEqual(&best.Hash):
		Debug("new best block hash")
		getNew = true
		fallthrough
	case c.lastTxUpdate.Load() != c.blockTemplateGenerator.GetTxSource().LastUpdated() && time.Now().
		After(time.Unix(0, c.lastGenerated.Load().(int64)+int64(time.Minute))):
		Trace("block is stale, regenerating")
		getNew = true
	}
	if getNew {
		if oB, err = c.GetTemplateMessageShards(); Check(err) {
			return
		}
		c.oldBlocks.Store(oB)
	}
	if err = c.SendShards(job.Magic, oB); Check(err) {
	}
	return
}

// GetTemplateMessageShards gets a new address, template message and returns FEC
// shards for the template, and saves the template
func (c *Controller) GetTemplateMessageShards() (o [][]byte, err error) {
	var addr util.Address
	if addr, err = c.GetNewAddressFromMiningAddrs(); Check(err) {
		if addr, err = c.GetNewAddressFromWallet(); Check(err) {
			return
		}
	}
	if c.msgBlockTemplate, err = c.GetMsgBlockTemplate(addr); !Check(err) {
		o = transport.GetShards(c.msgBlockTemplate.Serialize())
	}
	return
}

func (c *Controller) SendShards(magic []byte, data [][]byte) (err error) {
	if err = c.multiConn.SendMany(magic, data); Check(err) {
	}
	return
}

func (c *Controller) hashReport() float64 {
	c.hashSampleBuf.Add(c.hashCount.Load())
	av := ewma.NewMovingAverage()
	var i int
	var prev uint64
	if err := c.hashSampleBuf.ForEach(
		func(v uint64) error {
			if i < 1 {
				prev = v
			} else {
				interval := v - prev
				av.Add(float64(interval))
				prev = v
			}
			i++
			return nil
		},
	); Check(err) {
	}
	return av.Value()
}
