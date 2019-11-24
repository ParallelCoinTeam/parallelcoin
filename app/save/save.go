package save

import (
	"encoding/json"
	"io/ioutil"

	"github.com/p9c/pod/app/apputil"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/pod"
)

func Pod(c *pod.Config) (success bool) {
	log.TRACE("saving configuration to", *c.ConfigFile)
	if yp, e := json.MarshalIndent(c, "", "  "); e == nil {
		apputil.EnsureDir(*c.ConfigFile)
		if e := ioutil.WriteFile(*c.ConfigFile, yp, 0600); e != nil {
			log.ERROR(e)
			return
		}
		return true
	}
	return
}
