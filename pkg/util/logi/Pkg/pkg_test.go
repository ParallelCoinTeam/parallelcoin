package Pkg

import (
	"testing"

	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
)

func TestPackage(t *testing.T) {
	testPkgs := Pk.Package{
		"testing1": false,
		"testing2": true,
		"testing3": true,
	}
	d := Get(testPkgs).Data
	c := LoadContainer(d)
	t.Log(c.String())
}
