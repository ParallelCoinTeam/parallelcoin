package rcd

import (
	"encoding/hex"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
	"time"
)

func (r *RcVar) CreateWallet(privPassphrase, duoSeed, pubPassphrase, walletDir string) {
	var err error
	var seed []byte
	if walletDir == "" {
		walletDir = *r.cx.Config.WalletFile
	}
	l := wallet.NewLoader(r.cx.ActiveNet, *r.cx.Config.WalletFile, 250)

	if duoSeed == "" {
		seed, err = hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
		if err != nil {
			log.ERROR(err)
			panic(err)
		}
	} else {
		seed, err = hex.DecodeString(duoSeed)
		if err != nil {
			// Need to make JS invocation to embed
			log.ERROR(err)
		}
	}

	_, err = l.CreateNewWallet([]byte(pubPassphrase), []byte(privPassphrase), seed, time.Now(), true, r.cx.Config)
	if err != nil {
		log.ERROR(err)
		panic(err)
	}

	r.Boot.IsFirstRun = false
	*r.cx.Config.WalletPass = pubPassphrase
	*r.cx.Config.WalletFile = walletDir

	save.Pod(r.cx.Config)
	//log.INFO(rc)
}
