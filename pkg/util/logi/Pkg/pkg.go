// Package Package is a message type for logi package filtering
package Pkg

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/coding/simplebuffer"
	"github.com/p9c/pod/pkg/coding/simplebuffer/Byte"
	"github.com/p9c/pod/pkg/coding/simplebuffer/String"
	"github.com/p9c/pod/pkg/util/logi/Pkg/Pk"
)

var PackageMagic = []byte{'p', 'k', 'g', 's'}

type Container struct {
	simplebuffer.Container
}

func Get(pkgs Pk.Package) Container {
	c := simplebuffer.Serializers{}
	for i := range pkgs {
		c = append(c, String.New().Put(i))
		x := 0
		if pkgs[i] {
			x = 1
		}
		c = append(c, Byte.New().Put(byte(x)))
	}
	return Container{*c.CreateContainer(PackageMagic)}
}

// LoadContainer takes a message byte slice payload and loads it into a container
// ready to be decoded
func LoadContainer(b []byte) (out *Container) {
	out = &Container{simplebuffer.Container{Data: b}}
	return
}

func (c *Container) GetPackages() (out Pk.Package) {
	out = make(Pk.Package)
	for i := 0; i < int(c.Count()/2); i++ {
		si := uint16(i * 2)
		s := String.New().DecodeOne(c.Get(si)).Get()
		b := Byte.New().DecodeOne(c.Get(si + 1)).Get()
		boo := false
		if b != 0 {
			boo = true
		}
		out[s] = boo
	}
	return
}

func (c *Container) String() (s string) {
	s = spew.Sdump(c.GetPackages())
	return
}
