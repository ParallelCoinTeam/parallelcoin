package app

import (
	"github.com/p9c/log"
	"github.com/p9c/pod/cmd/node/state"
	"github.com/p9c/pod/pkg/chainrpc"
	"github.com/p9c/pod/pkg/opts"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/qu"
	"go.uber.org/atomic"
	"math/rand"
	"reflect"
	"time"
)

// Main is the entrypoint for the pod suite
func Main() int {
	log.SetLogLevel("trace")
	var e error
	var cx *pod.State
	if cx, e = GetNewContext(GetDefaultConfig()); F.Chk(e) {
		return 1
	}
	if e = cx.Config.Initialize(); E.Chk(e) {
	}
	T.Ln("running command", cx.
		Config.
		RunningCommand.
		Name,
	)
	if e = cx.Config.RunningCommand.Entrypoint(cx); E.Chk(e) {
		return 1
	}
	return 0
}

// GetDefaultConfig returns a Config struct pristine factory fresh
func GetDefaultConfig() (c *opts.Config) {
	c = &opts.Config{
		Commands: GetCommands(),
		Map:      GetConfigs(),
	}
	c.RunningCommand = c.Commands[0]
	// I.S(c.Commands[0])
	// I.S(c.Map)
	t := reflect.ValueOf(c)
	t = t.Elem()
	for i := range c.Map {
		tf := t.FieldByName(i)
		if tf.IsValid() && tf.CanSet() && tf.CanAddr() {
			val := reflect.ValueOf(c.Map[i])
			tf.Set(val)
		}
	}
	// I.S(c)
	return
}

// GetNewContext returns a fresh new context
func GetNewContext(config *opts.Config) (s *pod.State, e error) {
	if e = config.Initialize(); E.Chk(e) {
		// return
		panic(e)
	}
	config.RunningCommand = config.Commands[0]
	chainClientReady := qu.T()
	rand.Seed(time.Now().UnixNano())
	rand.Seed(rand.Int63())
	s = &pod.State{
		ChainClientReady: chainClientReady,
		KillAll:          qu.T(),
		// App:              cli.NewApp(),
		Config:    config,
		ConfigMap: config.Map,
		StateCfg:  new(state.Config),
		// Language:         lang.ExportLanguage(appLang),
		// DataDir:          appdata.Dir(appName, false),
		NodeChan: make(chan *chainrpc.Server),
		Syncing:  atomic.NewBool(true),
	}
	return
}
