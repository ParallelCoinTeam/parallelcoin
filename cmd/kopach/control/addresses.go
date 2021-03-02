package control

import (
	"errors"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/util"
	"github.com/urfave/cli"
	"math/rand"
	"time"
)

// GetNewAddressFromWallet gets a new address from the wallet if it is
// connected, or returns an error
func (c *Controller) GetNewAddressFromWallet() (addr util.Address, err error) {
	if c.walletClient != nil {
		if !c.walletClient.Disconnected() {
			Debug("have access to a wallet, generating address")
			if addr, err = c.walletClient.GetNewAddress("default"); Check(err) {
			} else {
				Debug("-------- found address", addr)
			}
		}
	} else {
		err = errors.New("no wallet available for new address")
		Debug(err)
	}
	return
}

// GetNewAddressFromMiningAddrs tries to get an address from the mining
// addresses list in the configuration file
func (c *Controller) GetNewAddressFromMiningAddrs() (addr util.Address, err error) {
	if c.cx.Config.MiningAddrs == nil {
		err = errors.New("mining addresses is nil")
		Debug(err)
		return
	}
	if len(*c.cx.Config.MiningAddrs) < 1 {
		err = errors.New("no mining addresses")
		Debug(err)
		return
	}
	// Choose a payment address at random.
	rand.Seed(time.Now().UnixNano())
	p2a := rand.Intn(len(*c.cx.Config.MiningAddrs))
	addr = c.cx.StateCfg.ActiveMiningAddrs[p2a]
	// remove the address from the state
	if p2a == 0 {
		c.cx.StateCfg.ActiveMiningAddrs = c.cx.StateCfg.ActiveMiningAddrs[1:]
	} else {
		c.cx.StateCfg.ActiveMiningAddrs = append(
			c.cx.StateCfg.ActiveMiningAddrs[:p2a],
			c.cx.StateCfg.ActiveMiningAddrs[p2a+1:]...,
		)
	}
	// update the config
	var ma cli.StringSlice
	for i := range c.cx.StateCfg.ActiveMiningAddrs {
		ma = append(ma, c.cx.StateCfg.ActiveMiningAddrs[i].String())
	}
	*c.cx.Config.MiningAddrs = ma
	save.Pod(c.cx.Config)
	return
}
