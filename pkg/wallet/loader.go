package wallet

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	`github.com/p9c/pod/pkg/chain/config/netparams`
	"github.com/p9c/pod/pkg/util/cl"
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

// CreateNewWallet creates a new wallet using the provided public and private passphrases.  The seed is optional.  If non-nil, addresses are derived from this seed.  If nil, a secure random seed is generated.
func (l *Loader) CreateNewWallet(pubPassphrase, privPassphrase, seed []byte, bday time.Time) (*Wallet, error) {
	defer l.Mutex.Unlock()
	l.Mutex.Lock()
	if l.Loaded {
		return nil, ErrLoaded
	}
	dbPath := filepath.Join(l.DDDirPath, WalletDbName)
	exists, err := fileExists(dbPath)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("Wallet ERROR: " + dbPath + " already exists")
	}
	// Create the wallet database backed by bolt db.
	err = os.MkdirAll(l.DDDirPath, 0700)
	if err != nil {
		return nil, err
	}
	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		return nil, err
	}
	// Initialize the newly created database for the wallet before opening.
	err = Create(db, pubPassphrase, privPassphrase, seed, l.ChainParams, bday)
	if err != nil {
		return nil, err
	}
	// Open the newly-created wallet.
	w, err := Open(db, pubPassphrase, nil, l.ChainParams, l.RecoveryWindow)
	if err != nil {
		return nil, err
	}
	w.Start()
	l.onLoaded(db)
	return w, nil
}

// LoadedWallet returns the loaded wallet, if any, and a bool for whether the
// wallet has been loaded or not.  If true, the wallet pointer should be safe to
// dereference.
func (l *Loader) LoadedWallet() (*Wallet, bool) {
	l.Mutex.Lock()
	w := l.Wallet
	l.Mutex.Unlock()
	return w, w != nil
}

// OpenExistingWallet opens the wallet from the loader's wallet database path and the public passphrase.  If the loader is being called by a context where standard input prompts may be used during wallet upgrades, setting canConsolePrompt will enables these prompts.
func (l *Loader) OpenExistingWallet(pubPassphrase []byte, canConsolePrompt bool) (*Wallet, error) {
	defer l.Mutex.Unlock()
	l.Mutex.Lock()
	// log <- cl.Info{"opening existing wallet", l.DDDirPath, cl.Ine()}
	if l.Loaded {
		log <- cl.Info{"already loaded wallet"}
		return nil, ErrLoaded
	}
	// Ensure that the network directory exists.
	if err := checkCreateDir(l.DDDirPath); err != nil {
		log <- cl.Error{"cannot create directory", l.DDDirPath, cl.Ine()}
		return nil, err
	}
	// log <- cl.Info{"directory exists", cl.Ine()}
	// Open the database using the boltdb backend.
	dbPath := filepath.Join(l.DDDirPath, WalletDbName)
	// log <- cl.Info{"opening database", dbPath, cl.Ine()}
	db, err := walletdb.Open("bdb", dbPath)
	if err != nil {
		log <- cl.Error{"failed to open database '", l.DDDirPath, "':", err, cl.Ine()}
		return nil, err
	}
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
	log <- cl.Trace{"opening wallet", cl.Ine()}
	w, err := Open(db, pubPassphrase, cbs, l.ChainParams, l.RecoveryWindow)
	if err != nil {
		log <- cl.Info{"failed to open wallet", err, cl.Ine()}
		// If opening the wallet fails (e.g. because of wrong
		// passphrase), we must close the backing database to
		// allow future calls to walletdb.Open().
		e := db.Close()
		if e != nil {
			log <- cl.Warn{"error closing database:", e}
		}
		return nil, err
	}
	l.Wallet = w
	log <- cl.Trace{"starting wallet", w != nil, cl.Ine()}
	w.Start()
	log <- cl.Trace{"waiting for load", db != nil, cl.Ine()}
	l.onLoaded(db)
	log <- cl.Trace{"wallet opened successfully", w != nil, cl.Ine()}
	return w, nil
}

// RunAfterLoad adds a function to be executed when the loader creates or opens
// a wallet.  Functions are executed in a single goroutine in the order they
// are added.
func (l *Loader) RunAfterLoad(fn func(*Wallet)) {
	l.Mutex.Lock()
	if l.Loaded {
		// w := l.Wallet
		l.Mutex.Unlock()
		fn(l.Wallet)
	} else {
		l.Callbacks = append(l.Callbacks, fn)
		l.Mutex.Unlock()
	}
}

// UnloadWallet stops the loaded wallet, if any, and closes the wallet database.
// This returns ErrNotLoaded if the wallet has not been loaded with
// CreateNewWallet or LoadExistingWallet.  The Loader may be reused if this
// function returns without error.
func (l *Loader) UnloadWallet() error {
	log <- cl.Trace{"unloading wallet", cl.Ine()}
	defer l.Mutex.Unlock()
	l.Mutex.Lock()
	if l.Wallet == nil {
		log <- cl.Debug{"wallet not loaded"}
		return ErrNotLoaded
	}
	log <- cl.Trace{"wallet stopping", cl.Ine()}
	l.Wallet.Stop()
	log <- cl.Trace{"waiting for wallet shutdown", cl.Ine()}
	l.Wallet.WaitForShutdown()
	if l.DB == nil {
		log <- cl.Debug{"there was no database", cl.Ine()}
		return ErrNotLoaded
	}
	log <- cl.Trace{"wallet stopped", cl.Ine()}
	err := l.DB.Close()
	if err != nil {
		log <- cl.Debug{"error closing database", err, cl.Ine()}
		return err
	}
	log <- cl.Trace{"database closed", cl.Ine()}
	time.Sleep(time.Second / 4)
	l.Loaded = false
	l.DB = nil
	return nil
}

// WalletExists returns whether a file exists at the loader's database path.
// This may return an error for unexpected I/O failures.
func (l *Loader) WalletExists() (bool, error) {
	dbPath := filepath.Join(l.DDDirPath, WalletDbName)
	return fileExists(dbPath)
}

// onLoaded executes each added callback and prevents loader from loading any
// additional wallets.  Requires mutex to be locked.
func (l *Loader) onLoaded(db walletdb.DB) {
	log <- cl.Trace{"wallet loader callbacks running ", l.Wallet != nil,
		cl.Ine()}
	for _, fn := range l.Callbacks {
		fn(l.Wallet)
	}
	log <- cl.Trace{"wallet loader callbacks finished", cl.Ine()}
	l.Loaded = true
	l.DB = db
	l.Callbacks = nil // not needed anymore
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
