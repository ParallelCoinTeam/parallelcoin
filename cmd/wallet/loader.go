package wallet

import (
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/p9c/qu"

	"github.com/p9c/matrjoska/pkg/chaincfg"
	"github.com/p9c/matrjoska/pkg/util/prompt"
	"github.com/p9c/matrjoska/pkg/waddrmgr"
	"github.com/p9c/matrjoska/pkg/walletdb"
	"github.com/p9c/matrjoska/pod/config"
)

// Loader implements the creating of new and opening of existing wallets, while providing a callback system for other
// subsystems to handle the loading of a wallet. This is primarily intended for use by the RPC servers, to enable
// methods and services which require the wallet when the wallet is loaded by another subsystem.
//
// Loader is safe for concurrent access.
type Loader struct {
	Callbacks      []func(*Wallet)
	ChainParams    *chaincfg.Params
	DDDirPath      string
	RecoveryWindow uint32
	Wallet         *Wallet
	Loaded         bool
	DB             walletdb.DB
	Mutex          sync.Mutex
}

const ()

var (
	// ErrExists describes the error condition of attempting to create a new wallet when one exists already.
	ErrExists = errors.New("wallet already exists")
	// ErrLoaded describes the error condition of attempting to load or create a wallet when the loader has already done
	// so.
	ErrLoaded = errors.New("wallet already loaded")
	// ErrNotLoaded describes the error condition of attempting to close a loaded wallet when a wallet has not been
	// loaded.
	ErrNotLoaded = errors.New("wallet is not loaded")
	errNoConsole = errors.New("db upgrade requires console access for additional input")
)

// CreateNewWallet creates a new wallet using the provided public and private passphrases. The seed is optional. If
// non-nil, addresses are derived from this seed. If nil, a secure random seed is generated.
func (ld *Loader) CreateNewWallet(
	pubPassphrase, privPassphrase, seed []byte,
	bday time.Time,
	noStart bool,
	podConfig *config.Config,
	quit qu.C,
) (w *Wallet, e error) {
	ld.Mutex.Lock()
	defer ld.Mutex.Unlock()
	if ld.Loaded {
		return nil, ErrLoaded
	}
	// dbPath := filepath.Join(ld.DDDirPath, WalletDbName)
	var exists bool
	if exists, e = fileExists(ld.DDDirPath); E.Chk(e) {
		return nil, e
	}
	if exists {
		return nil, errors.New("Wallet ERROR: " + ld.DDDirPath + " already exists")
	}
	// Create the wallet database backed by bolt db.
	p := filepath.Dir(ld.DDDirPath)
	if e = os.MkdirAll(p, 0700); E.Chk(e) {
		return nil, e
	}
	var db walletdb.DB
	if db, e = walletdb.Create("bdb", ld.DDDirPath); E.Chk(e) {
		return nil, e
	}
	// Initialize the newly created database for the wallet before opening.
	if e = Create(db, pubPassphrase, privPassphrase, seed, ld.ChainParams,
		bday); E.Chk(e) {
		return nil, e
	}
	// Open the newly-created wallet.
	if w, e = Open(db, pubPassphrase, nil, ld.ChainParams, ld.RecoveryWindow,
		podConfig, quit); E.Chk(e) {
		return nil, e
	}
	if !noStart {
		w.Start()
		ld.onLoaded(db)
	} else {
		if e = w.db.Close(); E.Chk(e) {
		}
	}
	return w, nil
}

// LoadedWallet returns the loaded wallet, if any, and a bool for whether the wallet has been loaded or not. If true,
// the wallet pointer should be safe to dereference.
func (ld *Loader) LoadedWallet() (*Wallet, bool) {
	ld.Mutex.Lock()
	w := ld.Wallet
	ld.Mutex.Unlock()
	return w, w != nil
}

// OpenExistingWallet opens the wallet from the loader's wallet database path and the public passphrase. If the loader
// is being called by a context where standard input prompts may be used during wallet upgrades, setting
// canConsolePrompt will enables these prompts.
func (ld *Loader) OpenExistingWallet(
	pubPassphrase []byte,
	canConsolePrompt bool,
	podConfig *config.Config,
	quit qu.C,
) (w *Wallet, e error) {
	defer ld.Mutex.Unlock()
	ld.Mutex.Lock()
	I.Ln("opening existing wallet", ld.DDDirPath)
	if ld.Loaded {
		I.Ln("already loaded wallet")
		return nil, ErrLoaded
	}
	// Ensure that the network directory exists.
	if e = checkCreateDir(filepath.Dir(ld.DDDirPath)); E.Chk(e) {
		E.Ln("cannot create directory", ld.DDDirPath)
		return nil, e
	}
	D.Ln("directory exists")
	// Open the database using the boltdb backend.
	dbPath := ld.DDDirPath
	I.Ln("opening database", dbPath)
	var db walletdb.DB
	if db, e = walletdb.Open("bdb", dbPath); E.Chk(e) {
		E.Ln("failed to open database '", ld.DDDirPath)
		return nil, e
	}
	I.Ln("opened wallet database")
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
	D.Ln("opening wallet '" + string(pubPassphrase) + "'")
	if w, e = Open(
		db,
		pubPassphrase,
		cbs,
		ld.ChainParams,
		ld.RecoveryWindow,
		podConfig,
		quit,
	); E.Chk(e) {
		E.Ln("failed to open wallet", e)
		// If opening the wallet fails (e.g. because of wrong passphrase), we must close the backing database to allow
		// future calls to walletdb.Open().
		if e = db.Close(); E.Chk(e) {
			W.Ln("error closing database:", e)
		}
		return nil, e
	}
	ld.Wallet = w
	D.Ln("starting wallet", w != nil)
	w.Start()
	D.Ln("waiting for load", db != nil)
	ld.onLoaded(db)
	D.Ln("wallet opened successfully", w != nil)
	return w, nil
}

// RunAfterLoad adds a function to be executed when the loader creates or opens a wallet. Functions are executed in a
// single goroutine in the order they are added.
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

// UnloadWallet stops the loaded wallet, if any, and closes the wallet database. This returns ErrNotLoaded if the wallet
// has not been loaded with CreateNewWallet or LoadExistingWallet. The Loader may be reused if this function returns
// without error.
func (ld *Loader) UnloadWallet() (e error) {
	F.Ln("unloading wallet")
	defer ld.Mutex.Unlock()
	ld.Mutex.Lock()
	if ld.Wallet == nil {
		D.Ln("wallet not loaded")
		return ErrNotLoaded
	}
	F.Ln("wallet stopping")
	ld.Wallet.Stop()
	F.Ln("waiting for wallet shutdown")
	ld.Wallet.WaitForShutdown()
	if ld.DB == nil {
		D.Ln("there was no database")
		return ErrNotLoaded
	}
	F.Ln("wallet stopped")
	e = ld.DB.Close()
	if e != nil {
		D.Ln("error closing database", e)
		return e
	}
	F.Ln("database closed")
	ld.Loaded = false
	ld.DB = nil
	return nil
}

// WalletExists returns whether a file exists at the loader's database path. This may return an error for unexpected I/O
// failures.
func (ld *Loader) WalletExists() (bool, error) {
	return fileExists(ld.DDDirPath)
}

// onLoaded executes each added callback and prevents loader from loading any additional wallets. Requires mutex to be
// locked.
func (ld *Loader) onLoaded(db walletdb.DB) {
	D.Ln("wallet loader callbacks running ", ld.Wallet != nil)
	for i, fn := range ld.Callbacks {
		D.Ln("running wallet loader callback", i)
		fn(ld.Wallet)
	}
	D.Ln("wallet loader callbacks finished")
	ld.Loaded = true
	ld.DB = db
	ld.Callbacks = nil // not needed anymore
}

// NewLoader constructs a Loader with an optional recovery window. If the recovery window is non-zero, the wallet will
// attempt to recovery addresses starting from the last SyncedTo height.
func NewLoader(
	chainParams *chaincfg.Params, dbDirPath string, recoveryWindow uint32,
) *Loader {
	l := &Loader{
		ChainParams:    chainParams,
		DDDirPath:      dbDirPath,
		RecoveryWindow: recoveryWindow,
	}
	return l
}
func fileExists(filePath string) (bool, error) {
	_, e := os.Stat(filePath)
	if e != nil {
		if os.IsNotExist(e) {
			return false, nil
		}
		return false, e
	}
	return true, nil
}
func noConsole() ([]byte, error) {
	return nil, errNoConsole
}
