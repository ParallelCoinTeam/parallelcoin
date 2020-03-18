package wallettx

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/util/prompt"
	waddrmgr "github.com/p9c/pod/pkg/wallet/addrmgr"
	walletdb "github.com/p9c/pod/pkg/wallet/db"
)

// A bunch of constants
const (
	WalletDbName = "wallet.db"
)

var (
	// ErrLoaded describes the error condition of attempting to load or
	// create a wallet when the loader has already done so.
	ErrLoaded = errors.New("wallet already loaded")
	// ErrNotLoaded describes the error condition of attempting to close a
	// loaded wallet when a wallet has not been loaded.
	ErrNotLoaded = errors.New("wallet is not loaded")
	// ErrExists describes the error condition of attempting to create a new
	// wallet when one exists already.
	ErrExists = errors.New("wallet already exists")
)

// Loader implements the creating of new and opening of existing wallets, while
// providing a callback system for other subsystems to handle the loading of a
// wallet.  This is primarily intended for use by the RPC servers, to enable
// methods and services which require the wallet when the wallet is loaded by
// another subsystem.
//
// Loader is safe for concurrent access.
type Loader struct {
	callbacks      []func(*Wallet)
	chainParams    *netparams.Params
	dbDirPath      string
	recoveryWindow uint32
	wallet         *Wallet
	db             walletdb.DB
	mu             sync.Mutex
}

// NewLoader constructs a Loader with an optional recovery window. If the
// recovery window is non-zero, the wallet will attempt to recovery addresses
// starting from the last SyncedTo height.
func NewLoader(chainParams *netparams.Params, dbDirPath string,
	recoveryWindow uint32) *Loader {
	return &Loader{
		chainParams:    chainParams,
		dbDirPath:      dbDirPath,
		recoveryWindow: recoveryWindow,
	}
}

// onLoaded executes each added callback and prevents loader from loading any
// additional wallets.  Requires mutex to be locked.
func (l *Loader) onLoaded(w *Wallet, db walletdb.DB) {
	for _, fn := range l.callbacks {
		fn(w)
	}
	l.wallet = w
	l.db = db
	l.callbacks = nil // not needed anymore
}

// RunAfterLoad adds a function to be executed when the loader creates or opens
// a wallet.  Functions are executed in a single goroutine in the order they are
// added.
func (l *Loader) RunAfterLoad(fn func(*Wallet)) {
	l.mu.Lock()
	if l.wallet != nil {
		w := l.wallet
		l.mu.Unlock()
		fn(w)
	} else {
		l.callbacks = append(l.callbacks, fn)
		l.mu.Unlock()
	}
}

// CreateNewWallet creates a new wallet using the provided public and private
// passphrases.  The seed is optional.  If non-nil, addresses are derived from
// this seed.  If nil, a secure random seed is generated.
func (l *Loader) CreateNewWallet(pubPassphrase, privPassphrase, seed []byte, bday time.Time) (*Wallet, error) {
	defer l.mu.Unlock()
	l.mu.Lock()
	if l.wallet != nil {
		return nil, ErrLoaded
	}
	dbPath := filepath.Join(l.dbDirPath, WalletDbName)
	exists, err := fileExists(dbPath)
	if err != nil {
		L.Error(err)
		return nil, err
	}
	if exists {
		return nil, errors.New("ERROR: " + dbPath + " already exists")
	}
	// Create the wallet database backed by bolt db.
	err = os.MkdirAll(l.dbDirPath, 0700)
	if err != nil {
		L.Error(err)
		return nil, err
	}
	db, err := walletdb.Create("bdb", dbPath)
	if err != nil {
		L.Error(err)
		return nil, err
	}
	// Initialize the newly created database for the wallet before opening.
	err = Create(
		db, pubPassphrase, privPassphrase, seed, l.chainParams, bday,
	)
	if err != nil {
		L.Error(err)
		return nil, err
	}
	// Open the newly-created wallet.
	w, err := Open(db, pubPassphrase, nil, l.chainParams, l.recoveryWindow)
	if err != nil {
		L.Error(err)
		return nil, err
	}
	w.Start()
	l.onLoaded(w, db)
	return w, nil
}

var errNoConsole = errors.New("db upgrade requires console access for additional input")

func noConsole() ([]byte, error) {
	return nil, errNoConsole
}

// OpenExistingWallet opens the wallet from the loader's wallet database path
// and the public passphrase.  If the loader is being called by a context where
// standard input prompts may be used during wallet upgrades, setting
// canConsolePrompt will enables these prompts.
func (l *Loader) OpenExistingWallet(pubPassphrase []byte, canConsolePrompt bool) (*Wallet, error) {
	defer l.mu.Unlock()
	l.mu.Lock()
	if l.wallet != nil {
		return nil, ErrLoaded
	}
	// Ensure that the network directory exists.
	if err := checkCreateDir(l.dbDirPath); err != nil {
		return nil, err
	}
	// Open the database using the boltdb backend.
	dbPath := filepath.Join(l.dbDirPath, WalletDbName)
	db, err := walletdb.Open("bdb", dbPath)
	if err != nil {
		L.Error(err)
		L.Error(
			"failed to open database:", err)
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
	w, err := Open(db, pubPassphrase, cbs, l.chainParams, l.recoveryWindow)
	if err != nil {
		L.Error(err)
		// If opening the wallet fails (e.g. because of wrong
		// passphrase), we must close the backing database to
		// allow future calls to walletdb.Open().
		e := db.Close()
		if e != nil {
			L.Warn(
				"error closing database:", e)
		}
		return nil, err
	}
	w.Start()
	l.onLoaded(w, db)
	return w, nil
}

// WalletExists returns whether a file exists at the loader's database path.
// This may return an error for unexpected I/O failures.
func (l *Loader) WalletExists() (bool, error) {
	dbPath := filepath.Join(l.dbDirPath, WalletDbName)
	return fileExists(dbPath)
}

// LoadedWallet returns the loaded wallet, if any, and a bool for whether the
// wallet has been loaded or not.  If true, the wallet pointer should be safe to
// dereference.
func (l *Loader) LoadedWallet() (*Wallet, bool) {
	l.mu.Lock()
	w := l.wallet
	l.mu.Unlock()
	return w, w != nil
}

// UnloadWallet stops the loaded wallet, if any, and closes the wallet database.
// This returns ErrNotLoaded if the wallet has not been loaded with
// CreateNewWallet or LoadExistingWallet.  The Loader may be reused if this
// function returns without error.
func (l *Loader) UnloadWallet() error {
	defer l.mu.Unlock()
	l.mu.Lock()
	if l.wallet == nil {
		return ErrNotLoaded
	}
	l.wallet.Stop()
	l.wallet.WaitForShutdown()
	err := l.db.Close()
	if err != nil {
		L.Error(err)
		return err
	}
	l.wallet = nil
	l.db = nil
	return nil
}
func fileExists(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err != nil {
		L.Error(err)
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
