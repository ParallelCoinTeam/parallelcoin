package walletmain

import (
	"bufio"
	"os"
	"path/filepath"
	"time"
	
	"github.com/p9c/pod/pkg/blockchain/chaincfg/netparams"
	"github.com/p9c/pod/pkg/blockchain/wire"
	"github.com/p9c/pod/pkg/database/walletdb"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/legacy/keystore"
	"github.com/p9c/pod/pkg/util/prompt"
	"github.com/p9c/pod/pkg/wallet"
	waddrmgr "github.com/p9c/pod/pkg/wallet/waddrmgr"
	
	// This initializes the bdb driver
	_ "github.com/p9c/pod/pkg/database/walletdb/bdb"
)

const slash = string(os.PathSeparator)

// CreateSimulationWallet is intended to be called from the rpcclient and used
// to create a wallet for actors involved in simulations.
func CreateSimulationWallet(activenet *netparams.Params, cfg *Config) (e error) {
	// Simulation wallet password is 'password'.
	privPass := []byte("password")
	// Public passphrase is the default.
	pubPass := []byte(wallet.InsecurePubPassphrase)
	netDir := NetworkDir(*cfg.AppDataDir, activenet)
	// Create the wallet.
	dbPath := filepath.Join(netDir, WalletDbName)
	inf.Ln("Creating the wallet...")
	// Create the wallet database backed by bolt db.
	db, e := walletdb.Create("bdb", dbPath)
	if e != nil {
		return e
	}
	defer func() {
		if e := db.Close(); err.Chk(e) {
		}
	}()
	// Create the wallet.
	e = wallet.Create(db, pubPass, privPass, nil, activenet, time.Now())
	if e != nil {
		return e
	}
	inf.Ln("The wallet has been created successfully.")
	return nil
}

// CreateWallet prompts the user for information needed to generate a new wallet and generates the wallet accordingly.
// The new wallet will reside at the provided path.
func CreateWallet(activenet *netparams.Params, config *pod.Config) (e error) {
	dbDir := *config.WalletFile
	loader := wallet.NewLoader(activenet, dbDir, 250)
	dbg.Ln("WalletPage", loader.ChainParams.Name)
	// When there is a legacy keystore, open it now to ensure any errors don't end up exiting the process after the user
	// has spent time entering a bunch of information.
	netDir := NetworkDir(*config.DataDir, activenet)
	keystorePath := filepath.Join(netDir, keystore.Filename)
	var legacyKeyStore *keystore.Store
	_, e = os.Stat(keystorePath)
	if e != nil && !os.IsNotExist(e) {
		// A stat error not due to a non-existant file should be returned to the caller.
		return e
	} else if e == nil {
		// Keystore file exists.
		legacyKeyStore, e = keystore.OpenDir(netDir)
		if e != nil {
			return e
		}
	}
	// Start by prompting for the private passphrase. When there is an existing keystore, the user will be promped for
	// that passphrase, otherwise they will be prompted for a new one.
	reader := bufio.NewReader(os.Stdin)
	privPass, e := prompt.PrivatePass(reader, legacyKeyStore)
	if e != nil {
		dbg.Ln(e)
		time.Sleep(time.Second * 3)
		return e
	}
	// When there exists a legacy keystore, unlock it now and set up a callback to import all keystore keys into the new
	// walletdb wallet
	if legacyKeyStore != nil {
		e = legacyKeyStore.Unlock(privPass)
		if e != nil {
			return e
		}
		// Import the addresses in the legacy keystore to the new wallet if any exist, locking each wallet again when
		// finished.
		loader.RunAfterLoad(
			func(w *wallet.Wallet) {
				defer func() {
					e := legacyKeyStore.Lock()
					if e != nil {
						dbg.Ln(e)
					}
				}()
				inf.Ln("Importing addresses from existing wallet...")
				lockChan := make(chan time.Time, 1)
				defer func() {
					lockChan <- time.Time{}
				}()
				e := w.Unlock(privPass, lockChan)
				if e != nil {
					err.F(
						"ERR: Failed to unlock new wallet "+
							"during old wallet key import: %v", e,
					)
					return
				}
				e = convertLegacyKeystore(legacyKeyStore, w)
				if e != nil {
					err.F(
						"ERR: Failed to import keys from old "+
							"wallet format: %v %s", e,
					)
					return
				}
				// Remove the legacy key store.
				e = os.Remove(keystorePath)
				if e != nil {
					err.Ln(
						"WARN: Failed to remove legacy wallet "+
							"from'%s'\n", keystorePath,
					)
				}
			},
		)
	}
	// Ascertain the public passphrase. This will either be a value specified by the user or the default hard-coded
	// public passphrase if the user does not want the additional public data encryption.
	pubPass, e := prompt.PublicPass(reader, privPass, []byte(""), []byte(*config.WalletPass))
	if e != nil {
		dbg.Ln(e)
		time.Sleep(time.Second * 5)
		return e
	}
	// Ascertain the wallet generation seed. This will either be an automatically generated value the user has already
	// confirmed or a value the user has entered which has already been validated.
	seed, e := prompt.Seed(reader)
	if e != nil {
		dbg.Ln(e)
		time.Sleep(time.Second * 5)
		return e
	}
	dbg.Ln("Creating the wallet")
	w, e := loader.CreateNewWallet(pubPass, privPass, seed, time.Now(), false, config, nil)
	if e != nil {
		dbg.Ln(e)
		time.Sleep(time.Second * 5)
		return e
	}
	w.Manager.Close()
	dbg.Ln("The wallet has been created successfully.")
	return nil
}

// NetworkDir returns the directory name of a network directory to hold wallet files.
func NetworkDir(dataDir string, chainParams *netparams.Params) string {
	netname := chainParams.Name
	// For now, we must always name the testnet data directory as "testnet" and not "testnet3" or any other version, as
	// the chaincfg testnet3 paramaters will likely be switched to being named "testnet3" in the future. This is done to
	// future proof that change, and an upgrade plan to move the testnet3 data directory can be worked out later.
	if chainParams.Net == wire.TestNet3 {
		netname = "testnet"
	}
	return filepath.Join(dataDir, netname)
}

// // checkCreateDir checks that the path exists and is a directory.
// // If path does not exist, it is created.
// func checkCreateDir(// 	path string) (e error) {
// 	if fi, e := os.Stat(path); err.Chk(e) {
// 		if os.IsNotExist(e) {
// 			// Attempt data directory creation
// 			if e = os.MkdirAll(path, 0700); err.Chk(e) {
// 				return fmt.Errorf("cannot create directory: %s", e)
// 			}
// 		} else {
// 			return fmt.Errorf("error checking directory: %s", e)
// 		}
// 	} else {
// 		if !fi.IsDir() {
// 			return fmt.Errorf("path '%s' is not a directory", path)
// 		}
// 	}
// 	return nil
// }

// convertLegacyKeystore converts all of the addresses in the passed legacy key store to the new waddrmgr.Manager
// format. Both the legacy keystore and the new manager must be unlocked.
func convertLegacyKeystore(legacyKeyStore *keystore.Store, w *wallet.Wallet) (e error) {
	netParams := legacyKeyStore.Net()
	blockStamp := waddrmgr.BlockStamp{
		Height: 0,
		Hash:   *netParams.GenesisHash,
	}
	for _, walletAddr := range legacyKeyStore.ActiveAddresses() {
		switch addr := walletAddr.(type) {
		case keystore.PubKeyAddress:
			privKey, e := addr.PrivKey()
			if e != nil {
				wrn.F(
					"Failed to obtain private key "+
						"for address %v: %v", addr.Address(),
					err,
				)
				continue
			}
			wif, e := util.NewWIF(
				privKey,
				netParams, addr.Compressed(),
			)
			if e != nil {
				err.Ln(
					"Failed to create wallet "+
						"import format for address %v: %v",
					addr.Address(), e,
				)
				continue
			}
			_, e = w.ImportPrivateKey(
				waddrmgr.KeyScopeBIP0044,
				wif, &blockStamp, false,
			)
			if e != nil {
				wrn.F(
					"WARN: Failed to import private "+
						"key for address %v: %v",
					addr.Address(), e,
				)
				continue
			}
		case keystore.ScriptAddress:
			_, e := w.ImportP2SHRedeemScript(addr.Script())
			if e != nil {
				wrn.F(
					"WARN: Failed to import "+
						"pay-to-script-hash script for "+
						"address %v: %v\n", addr.Address(), e,
				)
				continue
			}
		default:
			wrn.F(
				"WARN: Skipping unrecognized legacy "+
					"keystore type: %T\n", addr,
			)
			continue
		}
	}
	return nil
}
