package rpctest

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
	
	rpc "github.com/p9c/matrjoska/pkg/rpcclient"
	"github.com/p9c/matrjoska/pkg/util"
)

// nodeConfig contains all the args and data required to launch a pod process and connect the rpc client to it.
type nodeConfig struct {
	rpcUser      string
	rpcPass      string
	listen       string
	rpcListen    string
	rpcConnect   string
	dataDir      string
	logDir       string
	profile      string
	debugLevel   string
	extra        []string
	prefix       string
	exe          string
	endpoint     string
	certFile     string
	keyFile      string
	certificates []byte
}

// newConfig returns a newConfig with all default values.
func newConfig(prefix, certFile, keyFile string, extra []string) (*nodeConfig, error) {
	podPath, e := podExecutablePath()
	if e != nil {
		podPath = "pod"
	}
	a := &nodeConfig{
		listen:    "127.0.0.1:41047",
		rpcListen: "127.0.0.1:41048",
		rpcUser:   "user",
		rpcPass:   "pass",
		extra:     extra,
		prefix:    prefix,
		exe:       podPath,
		endpoint:  "ws",
		certFile:  certFile,
		keyFile:   keyFile,
	}
	if e := a.setDefaults(); E.Chk(e) {
		return nil, e
	}
	return a, nil
}

// setDefaults sets the default values of the config. It also creates the temporary data, and log directories which must
// be cleaned up with a call to cleanup().
func (n *nodeConfig) setDefaults() (e error) {
	datadir, e := ioutil.TempDir("", n.prefix+"-data")
	if e != nil {
		return e
	}
	n.dataDir = datadir
	logdir, e := ioutil.TempDir("", n.prefix+"-logs")
	if e != nil {
		return e
	}
	n.logDir = logdir
	cert, e := ioutil.ReadFile(n.certFile)
	if e != nil {
		return e
	}
	n.certificates = cert
	return nil
}

// arguments returns an array of arguments that be used to launch the pod process.
func (n *nodeConfig) arguments() []string {
	args := []string{}
	if n.rpcUser != "" {
		// --rpcuser
		args = append(args, fmt.Sprintf("--rpcuser=%s", n.rpcUser))
	}
	if n.rpcPass != "" {
		// --rpcpass
		args = append(args, fmt.Sprintf("--rpcpass=%s", n.rpcPass))
	}
	if n.listen != "" {
		// --listen
		args = append(args, fmt.Sprintf("--listen=%s", n.listen))
	}
	if n.rpcListen != "" {
		// --rpclisten
		args = append(args, fmt.Sprintf("--rpclisten=%s", n.rpcListen))
	}
	if n.rpcConnect != "" {
		// --rpcconnect
		args = append(args, fmt.Sprintf("--rpcconnect=%s", n.rpcConnect))
	}
	// --rpccert
	args = append(args, fmt.Sprintf("--rpccert=%s", n.certFile))
	// --rpckey
	args = append(args, fmt.Sprintf("--rpckey=%s", n.keyFile))
	if n.dataDir != "" {
		// --datadir
		args = append(args, fmt.Sprintf("--datadir=%s", n.dataDir))
	}
	if n.logDir != "" {
		// --logdir
		args = append(args, fmt.Sprintf("--logdir=%s", n.logDir))
	}
	if n.profile != "" {
		// --profile
		args = append(args, fmt.Sprintf("--profile=%s", n.profile))
	}
	if n.debugLevel != "" {
		// --debuglevel
		args = append(args, fmt.Sprintf("--debuglevel=%s", n.debugLevel))
	}
	args = append(args, n.extra...)
	return args
}

// command returns the exec.Cmd which will be used to start the pod process.
func (n *nodeConfig) command() *exec.Cmd {
	return exec.Command(n.exe, n.arguments()...)
}

// rpcConnConfig returns the rpc connection config that can be used to connect to the pod process that is launched via
// Start().
func (n *nodeConfig) rpcConnConfig() rpc.ConnConfig {
	return rpc.ConnConfig{
		Host:                 n.rpcListen,
		Endpoint:             n.endpoint,
		User:                 n.rpcUser,
		Pass:                 n.rpcPass,
		Certificates:         n.certificates,
		DisableAutoReconnect: true,
	}
}

// String returns the string representation of this nodeConfig.
func (n *nodeConfig) String() string {
	return n.prefix
}

// cleanup removes the tmp data and log directories.
func (n *nodeConfig) cleanup() (e error) {
	dirs := []string{
		n.logDir,
		n.dataDir,
	}
	for _, dir := range dirs {
		if e = os.RemoveAll(dir); E.Chk(e) {
			E.F("Cannot remove dir %s: %v", dir, e)
		}
	}
	return e
}

// node houses the necessary state required to configure, launch, and manage a pod process.
type node struct {
	config  *nodeConfig
	cmd     *exec.Cmd
	pidFile string
	dataDir string
}

// newNode creates a new node instance according to the passed config. dataDir will be used to hold a file recording the
// pid of the launched process, and as the base for the log and data directories for pod.
func newNode(config *nodeConfig, dataDir string) (*node, error) {
	return &node{
		config:  config,
		dataDir: dataDir,
		cmd:     config.command(),
	}, nil
}

// start creates a new pod process and writes its pid in a file reserved for recording the pid of the launched process.
// This file can be used to terminate the process in case of a hang or panic. In the case of a failing test case, or
// panic, it is important that the process be stopped via stop( ) otherwise it will persist unless explicitly killed.
func (n *node) start() (e error) {
	if e = n.cmd.Start(); E.Chk(e) {
		return e
	}
	pid, e := os.Create(
		filepath.Join(
			n.dataDir,
			fmt.Sprintf("%s.pid", n.config),
		),
	)
	if e != nil {
		return e
	}
	n.pidFile = pid.Name()
	if _, e = fmt.Fprintf(pid, "%d\n", n.cmd.Process.Pid); E.Chk(e) {
		return e
	}
	if e = pid.Close(); E.Chk(e) {
		return e
	}
	return nil
}

// stop interrupts the running pod process process, and waits until it exits properly. On windows, interrupt is not
// supported so a kill signal is used instead
func (n *node) stop() (e error) {
	if n.cmd == nil || n.cmd.Process == nil {
		// return if not properly initialized or error starting the process
		return nil
	}
	defer func() {
		e := n.cmd.Wait()
		if e != nil {
		}
	}()
	if runtime.GOOS == "windows" {
		return n.cmd.Process.Signal(os.Kill)
	}
	return n.cmd.Process.Signal(os.Interrupt)
}

// cleanup cleanups process and args files. The file housing the pid of the created process will be deleted as well as
// any directories created by the process.
func (n *node) cleanup() (e error) {
	if n.pidFile != "" {
		if e := os.Remove(n.pidFile); E.Chk(e) {
			E.F("unable to remove file %s: %v", n.pidFile, e)
		}
	}
	return n.config.cleanup()
}

// shutdown terminates the running pod process and cleans up all file/directories created by node.
func (n *node) shutdown() (e error) {
	if e := n.stop(); E.Chk(e) {
		return e
	}
	if e := n.cleanup(); E.Chk(e) {
		return e
	}
	return nil
}

// genCertPair generates a key/cert pair to the paths provided.
func genCertPair(certFile, keyFile string) (e error) {
	org := "rpctest autogenerated cert"
	validUntil := time.Now().Add(10 * 365 * 24 * time.Hour)
	var key []byte
	var cert []byte
	cert, key, e = util.NewTLSCertPair(org, validUntil, nil)
	if e != nil {
		return e
	}
	// Write cert and key files.
	if e = ioutil.WriteFile(certFile, cert, 0666); E.Chk(e) {
		return e
	}
	if e = ioutil.WriteFile(keyFile, key, 0600); E.Chk(e) {
		defer func() {
			if e = os.Remove(certFile); E.Chk(e) {
			}
		}()
		return e
	}
	return nil
}
