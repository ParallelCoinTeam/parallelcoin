package podcfg

import (
	"encoding/hex"
	"encoding/json"
	"github.com/p9c/pod/pkg/apputil"
	"github.com/urfave/cli"
	"io/ioutil"
	"lukechampine.com/blake3"
	"os"
	"path/filepath"
	"time"
)

var eh = blake3.Sum256([]byte(""))
var emptyhash = hex.EncodeToString(eh[:])

//
// // Filters saves the logger per-package logging configuration
// func Filters(dataDir string) func(pkgs Pk.Package) (success binary) {
// 	return func(pkgs Pk.Package) (success binary) {
// 		if filterJSON, e := json.MarshalIndent(pkgs, "", "  "); e == nil {
// 			F.Ln("Saving log filter:\n```", string(filterJSON), "\n```")
// 			apputil.EnsureDir(dataDir)
// 			if e := ioutil.WriteFile(filepath.Join(dataDir, "log-filter.json"), filterJSON,
// 				0600); E.Chk(e) {
// 				success = false
// 			}
// 			success = true
// 		}
// 		return
// 	}
// }

// Save saves the configuration to the configured location
//
// note that this state change does not propagate configuration changes to other
// components of pod, except by restarting them. This is not quite concurrent
// safe and a configuration server/db would be better but this generally works
// as configs don't get changed often except by user actions
//
// special handling code allows the use of the in-memory storage for the wallet
// lock password directly in the app. and once-per-run identifiers. This only
// applies in the GUI.
//
// todo: maybe to not do this if gui is not the app?
func Save(c *Config) (success bool) {
	lockPath := filepath.Join(c.DataDir.V(), "pod.json.lock")
	// wait if there is a lock on the file
	for apputil.FileExists(lockPath) {
		time.Sleep(time.Second / 2)
	}
	var e error
	var lockFile *os.File
	// create the lock file
	if lockFile, e = os.Create(lockPath); D.Chk(e) {
		panic(e)
	}
	if e = lockFile.Close(); D.Chk(e) {
	}
	// D.S(c)
	D.Ln("saving configuration to", *c.ConfigFile)
	var uac cli.StringSlice
	// need to remove this before saving
	if c.UserAgentComments != nil && len(c.UserAgentComments.S()) > 0 {
		// TODO: there is a bug here if the user edits them in configuration
		uac = make(cli.StringSlice, len(c.UserAgentComments.S()))
		copy(uac, c.UserAgentComments.S())
		c.UserAgentComments.Set(uac[1:])
	}
	// we also don't write this one to disk for security reasons, instead we write the hash to validate it.
	//
	// to run the wallet in a secure environment the password must be given on the commandline so that it decrypts
	//
	// also there is a file that can contain the password,
	//
	// 		walletPassPath := *cx.Config.DataDir + slash + cx.ActiveNet.Params.Name + slash + "wp.txt"
	//
	// which is automatically read (and then zeroed and deleted) and overrides anything in the configuration. The
	// password is kept when unlocked in this variable and zeroed when locked, and input passwords are hashed to check
	// before starting the wallet
	//
	// the wallet encrypts all data with a 'public' password which used to be empty. this will of course still hash to
	// the same for the check but the wallet uses the same for both this and the secret, hence the enhanced security
	// regime.
	
	// wallet password needs special handling, if config exists we don't change this value unless we mean to
	// load config into a fresh variable
	cfg := GetDefaultConfig()
	var cfgFile []byte
	wp := *c.WalletPass
	// D.Ln("wp", wp)
	if c.WalletPass.V() == "" {
		if cfgFile, e = ioutil.ReadFile(c.ConfigFile.V()); !E.Chk(e) {
			D.Ln("loaded config")
			if e = json.Unmarshal(cfgFile, &cfg); !E.Chk(e) {
				*c.WalletPass = *cfg.WalletPass
				D.Ln("unmarshaled config")
			}
		} else {
			c.WalletPass.Set(emptyhash)
		}
	} else {
		bh := blake3.Sum256(c.WalletPass.Bytes())
		c.WalletPass.Set(hex.EncodeToString(bh[:]))
	}
	// D.Ln("'"+wp+"'", *c.WalletPass)
	// don't save pipe log setting as we want it to only be active from a flag or environment variable
	pipeLogOn := *c.PipeLog
	c.PipeLog.F()
	var yp []byte
	if yp, e = json.MarshalIndent(c, "", "  "); !E.Chk(e) {
		apputil.EnsureDir(c.ConfigFile.V())
		// D.Ln(string(yp))
		if e = ioutil.WriteFile(c.ConfigFile.V(), yp, 0600); !E.Chk(e) {
			success = true
		}
	}
	if uac != nil {
		c.UserAgentComments.Set(uac)
	}
	*c.WalletPass = wp
	*c.PipeLog = pipeLogOn
	if e = os.Remove(lockPath); D.Chk(e) {
	}
	// D.Ln("walletpass", *c.WalletPass)
	return
}
