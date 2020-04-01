package wallet

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util/prompt"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	walletdb "github.com/p9c/pod/pkg/wallet/db"
)

// Loader implements the creating of new and opening of existing wallets, while
// providing a callback system for other subsystems to handle the loading of a
// wallet.  This is primarily intended for use by the RPC servers, to enable
// methods and services which require the wallet when the wallet is loaded by
// another subsystem.
//
// Loader is safe for concurrent access.
type Loader struct {
	Callbacks      []func(*Wallet)
	ChainParams    *netparams.Params
	DDDirPath      string
	RecoveryWindow uint32
	Wallet         *Wallet
	Loaded         bool
	DB             walletdb.DB
	Mutex          sync.Mutex
}

const (
	// WalletDbName is
	WalletDbName = "wallet.db"
)

var (
	// ErrExists describes the error condition of attempting to create a new
	// wallet when one exists already.
	ErrExists = errors.New("wallet already exists")
	// ErrLoaded describes the error condition of attempting to load or
	// create a wallet when the loader has already done so.
	ErrLoaded = errors.New("wallet already loaded")
	// ErrNotLoaded describes the error condition of attempting to close a
	// loaded wallet when a wallet has not been loaded.
	ErrNotLoaded = errors.New("wallet is not loaded")
	errNoConsole = errors.New("db upgrade requires console access for additional input")
)

// CreateNewWallet creates a new wallet using the provided public and private passphrases.  The seed is optional.  If
// non-nil, addresses are derived from this seed.  If nil, a secure random seed is generated.
func (ld *Loader) CreateNewWallet(pubPassphrase, privPassphrase, seed []byte, bday time.Time, noStart bool,
	podConfig *pod.Config) (*Wallet, error) {
	defer ld.Mutex.Unlock()
	ld.Mutex.Lock()
	if ld.Loaded {
		return nil, ErrLoaded
	}
	// dbPath := filepath.Join(ld.DDDirPath, WalletDbName)
	exists, err := fileExists(ld.DDDirPath)
	if err != nil {
		Error(err)
		return nil, err
	}
	if exists {
		return nil, errors.New("Wallet ERROR: " + ld.DDDirPath + " already exists")
	}
	// Create the wallet database backed by bolt db.
	p := filepath.Dir(ld.DDDirPath)
	err = os.MkdirAll(p, 0700)
	if err != nil {
		Error(err)
		return nil, err
	}
	db, err := walletdb.Create("bdb", ld.DDDirPath)
	if err != nil {
		Error(err)
		return nil, err
	}
	// Initialize the newly created database for the wallet before opening.
	err = Create(db, pubPassphrase, privPassphrase, seed, ld.ChainParams, bday)
	if err != nil {
		Error(err)
		return nil, err
	}
	// Open the newly-created wallet.
	w, err := Open(db, pubPassphrase, nil, ld.ChainParams, ld.RecoveryWindow, podConfig)
	if err != nil {
		Error(err)
		return nil, err
	}
	if !noStart {
		w.Start()
		ld.onLoaded(db)
	} else {
		w.db.Close()
	}
	return w, nil
}

// LoadedWallet returns the loaded wallet, if any, and a bool for whether the
// wallet has been loaded or not.  If true, the wallet pointer should be safe to
// dereference.
func (ld *Loader) LoadedWallet() (*Wallet, bool) {
	ld.Mutex.Lock()
	w := ld.Wallet
	ld.Mutex.Unlock()
	return w, w != nil
}

// OpenExistingWallet opens the wallet from the loader's wallet database path and the public passphrase.  If the loader is being called by a context where standard input prompts may be used during wallet upgrades, setting canConsolePrompt will enables these prompts.
func (ld *Loader) OpenExistingWallet(pubPassphrase []byte, canConsolePrompt bool, podConfig *pod.Config) (*Wallet, error) {
	defer ld.Mutex.Unlock()
	ld.Mutex.Lock()
	// INFO("opening existing wallet", ld.DDDirPath}
	if ld.Loaded {
		Info("already loaded wallet")
		return nil, ErrLoaded
	}
	// Ensure that the network directory exists.
	if err := checkCreateDir(filepath.Dir(ld.DDDirPath)); err != nil {
		Error("cannot create directory", ld.DDDirPath)
		return nil, err
	}
	Info("directory exists")
	// Open the database using the boltdb backend.
	dbPath := ld.DDDirPath
	Info("opening database", dbPath)
	db, err := walletdb.Open("bdb", dbPath)
	if err != nil {
		Error(err)
		Error("failed to open database '", ld.DDDirPath, "':", err)
		return nil, err
	}
	Info("opened wallet database")
	var cbs *waddrmgr.OpenCallbacks
	if canConsolePrompt {
		cbs = &waddrmgr.OpenCallbacks{
			ObtainSeed:        prompt.ProvideSeed,
			ObtainPrivatePass: prompt.ProvidePrivPassphrase,
		}
	} else {
		cbs = &waddrmgr.OpenCallbacks{
			ObtainSeed:        noConsole,
			ObtainPrivatePass: noConsole,
		}
	}
	Trace("opening wallet")
	w, err := Open(db, pubPassphrase, cbs, ld.ChainParams, ld.RecoveryWindow, podConfig)
	if err != nil {
		Error(err)
		Info("failed to open wallet", err)
		// If opening the wallet fails (e.g. because of wrong
		// passphrase), we must close the backing database to
		// allow future calls to walletdb.Open().
		e := db.Close()
		if e != nil {
			Warn("error closing database:", e)
		}
		return nil, err
	}
	ld.Wallet = w
	Trace("starting wallet", w != nil)
	w.Start()
	Trace("waiting for load", db != nil)
	ld.onLoaded(db)
	Trace("wallet opened successfully", w != nil)
	return w, nil
}

// RunAfterLoad adds a function to be executed when the loader creates or opens
// a wallet.  Functions are executed in a single goroutine in the order they
// are added.
func (ld *Loader) RunAfterLoad(fn func(*Wallet)) {
	ld.Mutex.Lock()
	if ld.Loaded {
		// w := ld.Wallet
		ld.Mutex.Unlock()
		fn(ld.Wallet)
	} else {
		ld.Callbacks = append(ld.Callbacks, fn)
		ld.Mutex.Unlock()
	}
}

// UnloadWallet stops the loaded wallet, if any, and closes the wallet database.
// This returns ErrNotLoaded if the wallet has not been loaded with
// CreateNewWallet or LoadExistingWallet.  The Loader may be reused if this
// function returns without error.
func (ld *Loader) UnloadWallet() error {
	Trace("unloading wallet")
	defer ld.Mutex.Unlock()
	ld.Mutex.Lock()
	if ld.Wallet == nil {
		Debug("wallet not loaded")
		return ErrNotLoaded
	}
	Trace("wallet stopping")
	ld.Wallet.Stop()
	Trace("waiting for wallet shutdown")
	ld.Wallet.WaitForShutdown()
	if ld.DB == nil {
		Debug("there was no database")
		return ErrNotLoaded
	}
	Trace("wallet stopped")
	err := ld.DB.Close()
	if err != nil {
		Error(err)
		Debug("error closing database", err)
		return err
	}
	Trace("database closed")
	time.Sleep(time.Second / 4)
	ld.Loaded = false
	ld.DB = nil
	return nil
}

// WalletExists returns whether a file exists at the loader's database path.
// This may return an error for unexpected I/O failures.
func (ld *Loader) WalletExists() (bool, error) {
	return fileExists(ld.DDDirPath)
}

// onLoaded executes each added callback and prevents loader from loading any
// additional wallets.  Requires mutex to be locked.
func (ld *Loader) onLoaded(db walletdb.DB) {
	Trace("wallet loader callbacks running ", ld.Wallet != nil)
	for _, fn := range ld.Callbacks {
		fn(ld.Wallet)
	}
	Trace("wallet loader callbacks finished")
	ld.Loaded = true
	ld.DB = db
	ld.Callbacks = nil // not needed anymore
}

// NewLoader constructs a Loader with an optional recovery window. If the
// recovery window is non-zero, the wallet will attempt to recovery addresses
// starting from the last SyncedTo height.
func NewLoader(chainParams *netparams.Params, dbDirPath string, recoveryWindow uint32) *Loader {
	l := &Loader{
		ChainParams:    chainParams,
		DDDirPath:      dbDirPath,
		RecoveryWindow: recoveryWindow,
	}
	return l
}
func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		Error(err)
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
func noConsole() ([]byte, error) {
	return nil, errNoConsole
}
