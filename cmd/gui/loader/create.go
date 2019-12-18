package loader

import (
	"encoding/hex"
	"github.com/p9c/pod/app/save"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/util/hdkeychain"
	"github.com/p9c/pod/pkg/wallet"
	"time"
)

func CreateWallet(ldr *DuoUIload, privPassphrase, duoSeed, pubPassphrase, walletDir string) {
	var err error
	var seed []byte
	if walletDir == "" {
		walletDir = *ldr.cx.Config.WalletFile
	}
	l := wallet.NewLoader(ldr.cx.ActiveNet, *ldr.cx.Config.WalletFile, 250)

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

	_, err = l.CreateNewWallet([]byte(pubPassphrase), []byte(privPassphrase), seed, time.Now(), true)
	if err != nil {
		log.ERROR(err)
		panic(err)
	}

	//duo.Boot.IsFirstRun = false
	*ldr.cx.Config.WalletPass = pubPassphrase
	*ldr.cx.Config.WalletFile = walletDir

	save.Pod(ldr.cx.Config)
	//log.INFO(rc)
}
