package save

import (
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	
	"lukechampine.com/blake3"
	
	"github.com/urfave/cli"
	
	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/pkg/pod"
)

var eh = blake3.Sum256([]byte(""))
var emptyhash = hex.EncodeToString(eh[:])

// Pod saves the configuration to the configured location
func Pod(c *pod.Config) (success bool) {
	c.Lock()
	defer c.Unlock()
	// Debugs(c)
	Debug("saving configuration to", *c.ConfigFile)
	var uac cli.StringSlice
	// need to remove this before saving
	if c.UserAgentComments != nil && len(*c.UserAgentComments) > 0 {
		// TODO: there is a bug here if the user edits them in configuration
		uac = make(cli.StringSlice, len(*c.UserAgentComments))
		copy(uac, *c.UserAgentComments)
		*c.UserAgentComments = uac[1:]
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
	cfg, _ := pod.EmptyConfig()
	var cfgFile []byte
	var err error
	wp := *c.WalletPass
	// Debug("wp", wp)
	if *c.WalletPass == "" {
		if cfgFile, err = ioutil.ReadFile(*c.ConfigFile); !Check(err) {
			Debug("loaded config")
			if err = json.Unmarshal(cfgFile, &cfg); !Check(err) {
				*c.WalletPass = *cfg.WalletPass
				// Debug("unmarshaled config", wp, *c.WalletPass)
			}
		} else {
			*c.WalletPass = emptyhash
		}
	} else {
		bh := blake3.Sum256([]byte(*c.WalletPass))
		*c.WalletPass = hex.EncodeToString(bh[:])
	}
	// Debug("'"+wp+"'", *c.WalletPass)
	// don't save pipe log setting as we want it to only be active from a flag or environment variable
	pipeLogOn := *c.PipeLog
	*c.PipeLog = false
	if yp, e := json.MarshalIndent(c, "", "  "); e == nil {
		apputil.EnsureDir(*c.ConfigFile)
		Debug(string(yp))
		if e := ioutil.WriteFile(*c.ConfigFile, yp, 0600); e != nil {
			Error(e)
			success = false
		}
		success = true
	}
	if c.UserAgentComments != nil {
		*c.UserAgentComments = uac
	}
	*c.WalletPass = wp
	*c.PipeLog = pipeLogOn
	Debug("walletpass", *c.WalletPass)
	return
}

//
// // Filters saves the logger per-package logging configuration
// func Filters(dataDir string) func(pkgs Pk.Package) (success bool) {
// 	return func(pkgs Pk.Package) (success bool) {
// 		if filterJSON, e := json.MarshalIndent(pkgs, "", "  "); e == nil {
// 			Trace("Saving log filter:\n```", string(filterJSON), "\n```")
// 			apputil.EnsureDir(dataDir)
// 			if e := ioutil.WriteFile(filepath.Join(dataDir, "log-filter.json"), filterJSON,
// 				0600); Check(e) {
// 				success = false
// 			}
// 			success = true
// 		}
// 		return
// 	}
// }
