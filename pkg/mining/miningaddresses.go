package mining

import (
	"github.com/p9c/pod/pkg/podcfg"
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/wallet"
	wm "github.com/p9c/pod/pkg/waddrmgr"
)

// RefillMiningAddresses adds new addresses to the mining address pool for the miner
// todo: make this remove ones that have been used or received a payment or mined
func RefillMiningAddresses(w *wallet.Wallet, cfg *podcfg.Config, stateCfg *state.Config) {
	if w == nil {
		D.Ln("trying to refill without a wallet")
		return
	}
	if cfg == nil {
		D.Ln("config is empty")
		return
	}
	var miningAddressLen int
	if cfg.MiningAddrs != nil {
		D.Ln("miningAddressLen", len(*cfg.MiningAddrs))
		miningAddressLen = len(*cfg.MiningAddrs)
	} else {
		D.Ln("miningaddrs slice is missing")
		cfg.MiningAddrs = new(cli.StringSlice)
	}
	toMake := 99 - miningAddressLen
	if miningAddressLen >= 99 {
		toMake = 0
	}
	if toMake < 1 {
		D.Ln("not making any new addresses")
		return
	}
	D.Ln("refilling mining addresses")
	account, e := w.AccountNumber(
		wm.KeyScopeBIP0044,
		"default",
	)
	if e != nil {
		E.Ln("error getting account number ", e)
	}
	for i := 0; i < toMake; i++ {
		addr, e := w.NewAddress(
			account, wm.KeyScopeBIP0044,
			true,
		)
		if e == nil {
			// add them to the configuration to be saved
			*cfg.MiningAddrs = append(*cfg.MiningAddrs, addr.EncodeAddress())
			// add them to the active mining address list so they
			// are ready to use
			stateCfg.ActiveMiningAddrs = append(stateCfg.ActiveMiningAddrs, addr)
		} else {
			E.Ln("error adding new address ", e)
		}
	}
	if podcfg.Save(cfg) {
		D.Ln("saved config with new addresses")
	} else {
		E.Ln("error adding new addresses", e)
	}
}
