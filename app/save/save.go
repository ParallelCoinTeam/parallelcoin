package save

import (
   "io/ioutil"
   
   "github.com/pelletier/go-toml"
   
   "git.parallelcoin.io/dev/pod/app/apputil"
   "git.parallelcoin.io/dev/pod/pkg/pod"
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
