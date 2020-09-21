package walletmain

import (
	"bufio"
	"github.com/stalker-loki/app/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/stalker-loki/pod/pkg/chain/config/netparams"
	"github.com/stalker-loki/pod/pkg/chain/wire"
	ec "github.com/stalker-loki/pod/pkg/coding/elliptic"
	"github.com/stalker-loki/pod/pkg/db/walletdb"
	"github.com/stalker-loki/pod/pkg/pod"
	"github.com/stalker-loki/pod/pkg/util"
	"github.com/stalker-loki/pod/pkg/util/legacy/keystore"
	"github.com/stalker-loki/pod/pkg/util/prompt"
	"github.com/stalker-loki/pod/pkg/wallet"
	waddrmgr "github.com/stalker-loki/pod/pkg/wallet/addrmgr"

	// This initializes the bdb driver
	_ "github.com/stalker-loki/pod/pkg/db/walletdb/bdb"
)

const slash = string(os.PathSeparator)

// CreateSimulationWallet is intended to be called from the rpcclient
// and used to create a wallet for actors involved in simulations.
func CreateSimulationWallet(activenet *netparams.Params, cfg *Config) error {
	// Simulation wallet password is 'password'.
	privPass := []byte("password")
	// Public passphrase is the default.
	pubPass := []byte(wallet.InsecurePubPassphrase)
	netDir := NetworkDir(*cfg.AppDataDir, activenet)
	// Create the wallet.
	dbPath := filepath.Join(netDir, WalletDbName)
	slog.Info("Creating the wallet...")
	// Create the wallet database backed by bolt db.
	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		slog.Error(err)
		return err
	}
	defer db.Close()
	// Create the wallet.
	err = wallet.Create(db, pubPass, privPass, nil, activenet, time.Now())
	if err != nil {
		slog.Error(err)
		return err
	}
	slog.Info("The wallet has been created successfully.")
	return nil
}

// CreateWallet prompts the user for information needed to generate a new
// wallet and generates the wallet accordingly.
// The new wallet will reside at the provided path.
func CreateWallet(activenet *netparams.Params, config *pod.Config) error {
	dbDir := *config.WalletFile
	loader := wallet.NewLoader(activenet, dbDir, 250)
	// When there is a legacy keystore, open it now to ensure any errors
	// don't end up exiting the process after the user has spent time
	// entering a bunch of information.
	netDir := NetworkDir(*config.DataDir, activenet)
	keystorePath := filepath.Join(netDir, keystore.Filename)
	var legacyKeyStore *keystore.Store
	_, err := os.Stat(keystorePath)
	if err != nil && !os.IsNotExist(err) {
		// A stat error not due to a non-existant file should be
		// returned to the caller.
		return err
	} else if err == nil {
		// Keystore file exists.
		legacyKeyStore, err = keystore.OpenDir(netDir)
		if err != nil {
			slog.Error(err)
			return err
		}
	}
	// Start by prompting for the private passphrase.  When there is an existing
	// keystore, the user will be promped for that passphrase, otherwise they
	// will be prompted for a new one.
	reader := bufio.NewReader(os.Stdin)
	privPass, err := prompt.PrivatePass(reader, legacyKeyStore)
	if err != nil {
		slog.Error(err)
		slog.Debug(err)
		time.Sleep(time.Second * 3)
		return err
	}
	// When there exists a legacy keystore, unlock it now and set up a callback
	// to import all keystore keys into the new walletdb wallet
	if legacyKeyStore != nil {
		err = legacyKeyStore.Unlock(privPass)
		if err != nil {
			slog.Error(err)
			return err
		}
		// Import the addresses in the legacy keystore to the new wallet if any
		// exist, locking each wallet again when finished.
		loader.RunAfterLoad(func(w *wallet.Wallet) {
			defer func() {
				err := legacyKeyStore.Lock()
				if err != nil {
					slog.Error(err)
					slog.Debug(err)
				}
			}()
			slog.Info("Importing addresses from existing wallet...")
			lockChan := make(chan time.Time, 1)
			defer func() {
				lockChan <- time.Time{}
			}()
			err := w.Unlock(privPass, lockChan)
			if err != nil {
				slog.Errorf("ERR: Failed to unlock new wallet "+
					"during old wallet key import: %v", err)
				return
			}
			err = convertLegacyKeystore(legacyKeyStore, w)
			if err != nil {
				slog.Errorf("ERR: Failed to import keys from old "+
					"wallet format: %v %s", err)
				return
			}
			// Remove the legacy key store.
			err = os.Remove(keystorePath)
			if err != nil {
				slog.Error("WARN: Failed to remove legacy wallet "+
					"from'%s'\n", keystorePath)
			}
		})
	}
	// Ascertain the public passphrase.  This will either be a value specified
	// by the user or the default hard-coded public passphrase if the user does
	// not want the additional public data encryption.
	pubPass, err := prompt.PublicPass(reader, privPass, []byte(""),
		[]byte(*config.WalletPass))
	if err != nil {
		slog.Error(err)
		slog.Debug(err)
		time.Sleep(time.Second * 5)
		return err
	}
	// Ascertain the wallet generation seed.  This will either be an
	// automatically generated value the user has already confirmed or a
	// value the user has entered which has already been validated.
	seed, err := prompt.Seed(reader)
	if err != nil {
		slog.Error(err)
		slog.Debug(err)
		time.Sleep(time.Second * 5)
		return err
	}
	slog.Debug("Creating the wallet")
	w, err := loader.CreateNewWallet(pubPass, privPass, seed, time.Now(), false, config)
	if err != nil {
		slog.Error(err)
		slog.Debug(err)
		time.Sleep(time.Second * 5)
		return err
	}
	w.Manager.Close()
	slog.Debug("The wallet has been created successfully.")
	return nil
}

// NetworkDir returns the directory name of a network directory to hold wallet files.
func NetworkDir(dataDir string, chainParams *netparams.Params) string {
	netname := chainParams.Name
	// For now, we must always name the testnet data directory as "testnet" and not
	// "testnet3" or any other version, as the chaincfg testnet3 paramaters will
	// likely be switched to being named "testnet3" in the future.  This is done to
	// future proof that change, and an upgrade plan to move the testnet3 data
	// directory can be worked out later.
	if chainParams.Net == wire.TestNet3 {
		netname = "testnet"
	}
	return filepath.Join(dataDir, netname)
}

// // checkCreateDir checks that the path exists and is a directory.
// // If path does not exist, it is created.
// func checkCreateDir(// 	path string) error {
// 	if fi, err := os.Stat(path); err != nil {
// 		if os.IsNotExist(err) {
// 			// Attempt data directory creation
// 			if err = os.MkdirAll(path, 0700); err != nil {
// 				return fmt.Errorf("cannot create directory: %s", err)
// 			}
// 		} else {
// 			return fmt.Errorf("error checking directory: %s", err)
// 		}
// 	} else {
// 		if !fi.IsDir() {
// 			return fmt.Errorf("path '%s' is not a directory", path)
// 		}
// 	}
// 	return nil
// }

// convertLegacyKeystore converts all of the addresses in the passed legacy
// key store to the new waddrmgr.Manager format.
// Both the legacy keystore and the new manager must be unlocked.
func convertLegacyKeystore(legacyKeyStore *keystore.Store, w *wallet.Wallet) error {
	netParams := legacyKeyStore.Net()
	blockStamp := waddrmgr.BlockStamp{
		Height: 0,
		Hash:   *netParams.GenesisHash,
	}
	for _, walletAddr := range legacyKeyStore.ActiveAddresses() {
		switch addr := walletAddr.(type) {
		case keystore.PubKeyAddress:
			privKey, err := addr.PrivKey()
			if err != nil {
				slog.Warnf("Failed to obtain private key "+
					"for address %v: %v", addr.Address(),
					err)
				continue
			}
			wif, err := util.NewWIF((*ec.PrivateKey)(privKey),
				netParams, addr.Compressed())
			if err != nil {
				slog.Warn("Failed to create wallet "+
					"import format for address %v: %v",
					addr.Address(), err)
				continue
			}
			_, err = w.ImportPrivateKey(waddrmgr.KeyScopeBIP0044,
				wif, &blockStamp, false)
			if err != nil {
				slog.Warnf("WARN: Failed to import private "+
					"key for address %v: %v",
					addr.Address(), err)
				continue
			}
		case keystore.ScriptAddress:
			_, err := w.ImportP2SHRedeemScript(addr.Script())
			if err != nil {
				slog.Warnf("WARN: Failed to import "+
					"pay-to-script-hash script for "+
					"address %v: %v\n", addr.Address(), err)
				continue
			}
		default:
			slog.Warnf("WARN: Skipping unrecognized legacy "+
				"keystore type: %T\n", addr)
			continue
		}
	}
	return nil
}
