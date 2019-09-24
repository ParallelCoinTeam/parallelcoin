package save

import (
   "io/ioutil"
   
   "github.com/pelletier/go-toml"
   
   "github.com/parallelcointeam/parallelcoin/app/apputil"
   "github.com/parallelcointeam/parallelcoin/pkg/pod"
)

func Pod(c *pod.Config) (success bool) {
	if yp, e := toml.Marshal(c); e == nil {
		apputil.EnsureDir(*c.ConfigFile)
		if e := ioutil.WriteFile(*c.ConfigFile, yp, 0600); e != nil {
			return
		}
		return true
	}
	return
}
