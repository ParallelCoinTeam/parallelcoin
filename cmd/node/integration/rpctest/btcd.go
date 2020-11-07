package rpctest

import (
	"fmt"
	"go/build"
	"os/exec"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/p9c/pod/pkg/util/gobin"
)

var (
	// compileMtx guards access to the executable path so that the project is only compiled once.
	compileMtx sync.Mutex
	// executablePath is the path to the compiled executable. This is the empty string until pod is compiled. This
	// should not be accessed directly; instead use the function podExecutablePath().
	executablePath string
)

// podExecutablePath returns a path to the pod executable to be used by rpctests. To ensure the code tests against the
// most up-to-date version of pod, this method compiles pod the first time it is called. After that, the generated
// binary is used for subsequent test harnesses. The executable file is not cleaned up, but since it lives at a static
// path in a temp directory, it is not a big deal.
func podExecutablePath() (string, error) {
	compileMtx.Lock()
	defer compileMtx.Unlock()
	// If pod has already been compiled, just use that.
	if len(executablePath) != 0 {
		return executablePath, nil
	}
	testDir, err := baseDir()
	if err != nil {
		Error(err)
		return "", err
	}
	// Determine import path of this package. Not necessarily pod if this is a forked repo.
	_, rpctestDir, _, ok := runtime.Caller(1)
	if !ok {
		return "", fmt.Errorf("Cannot get path to pod source code")
	}
	podPkgPath := filepath.Join(rpctestDir, "..", "..", "..")
	podPkg, err := build.ImportDir(podPkgPath, build.FindOnly)
	if err != nil {
		Error(err)
		return "", fmt.Errorf("Failed to build pod: %v", err)
	}
	// Build pod and output an executable in a static temp path.
	outputPath := filepath.Join(testDir, "pod")
	if runtime.GOOS == "windows" {
		outputPath += ".exe"
	}
	var gb string
	if gb, err = gobin.Get(); Check(err) {
		return "", err
	}
	cmd := exec.Command(gb, "build", "-o", outputPath, podPkg.ImportPath)
	err = cmd.Run()
	if err != nil {
		Error(err)
		return "", fmt.Errorf("Failed to build pod: %v", err)
	}
	// Save executable path so future calls do not recompile.
	executablePath = outputPath
	return executablePath, nil
}
