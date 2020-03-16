package save

import (
	"encoding/json"
	"io/ioutil"

	"github.com/urfave/cli"

	"github.com/p9c/pod/app/apputil"
	log "github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/pod"
)

// Pod saves the configuration to the configured location
func Pod(c *pod.Config) (success bool) {
	log.L.Trace("saving configuration to", *c.ConfigFile)
	var uac cli.StringSlice
	if len(*c.UserAgentComments) > 0 {
		uac = make(cli.StringSlice, len(*c.UserAgentComments))
		copy(uac, *c.UserAgentComments)
		*c.UserAgentComments = uac[1:]
	}
	if yp, e := json.MarshalIndent(c, "", "  "); e == nil {
		apputil.EnsureDir(*c.ConfigFile)
		if e := ioutil.WriteFile(*c.ConfigFile, yp, 0600); e != nil {
			log.L.Error(e)
			success = false
		}
		success = true
	}
	*c.UserAgentComments = uac

	return
}
