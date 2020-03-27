// Package Entry is a message type for logi log entries
package Entry

import (
	"crypto/rand"
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/logi"
	"github.com/p9c/pod/pkg/simplebuffer/String"
	"io"
	"time"

	"github.com/p9c/pod/pkg/simplebuffer"
	"github.com/p9c/pod/pkg/simplebuffer/Time"
)

var EntryMagic = []byte{'e', 'n', 't', 'r'}

type Container struct {
	simplebuffer.Container
}

func Get(ent *logi.Entry) Container {
	nonce := make([]byte, 4)
	if _, err := io.ReadFull(rand.Reader, nonce); L.Check(err) {
	}
	return Container{*simplebuffer.Serializers{
		Time.New().Put(ent.Time),
		String.New().Put(ent.Package),
		String.New().Put(ent.Level),
		String.New().Put(ent.CodeLocation),
		String.New().Put(ent.Text),
	}.CreateContainer(EntryMagic)}
}

// LoadContainer takes a message byte slice payload and loads it into a container
// ready to be decoded
func LoadContainer(b []byte) (out Container) {
	out.Data = b
	return
}

func (c *Container) GetTime() time.Time {
	return Time.New().DecodeOne(c.Get(0)).Get()
}

func (c *Container) GetPackage() string {
	return String.New().DecodeOne(c.Get(1)).Get()
}
func (c *Container) GetLevel() string {
	return String.New().DecodeOne(c.Get(2)).Get()
}
func (c *Container) GetCodeLocation() string {
	return String.New().DecodeOne(c.Get(3)).Get()
}
func (c *Container) GetText() string {
	return String.New().DecodeOne(c.Get(4)).Get()
}

func (c *Container) String() (s string) {
	spew.Sdump(*c.Struct())
	return
}

// Struct deserializes the data all in one go by calling the field deserializing
// functions into a structure containing the fields.
// The height is given in this report as it is part of the job message
// and makes it faster for clients to look up the algorithm name according to the
// block height, which can change between hard fork versions
func (c *Container) Struct() (out *logi.Entry) {
	out = &logi.Entry{
		Time:         c.GetTime(),
		Package:      c.GetPackage(),
		Level:        c.GetLevel(),
		CodeLocation: c.GetCodeLocation(),
		Text:         c.GetText(),
	}
	return
}
