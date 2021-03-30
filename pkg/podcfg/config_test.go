package podcfg

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/p9c/pod/pkg/opts"
	bool2 "github.com/p9c/pod/pkg/opts/binary"
	"github.com/p9c/pod/pkg/opts/duration"
	"github.com/p9c/pod/pkg/opts/float"
	"github.com/p9c/pod/pkg/opts/integer"
	"github.com/p9c/pod/pkg/opts/list"
	"github.com/p9c/pod/pkg/opts/text"
	"math/rand"
	"testing"
)

func TestForEach(t *testing.T) {
	c := GetDefaultConfig()
	c.ForEach(
		func(ifc opts.Option) bool {
			switch ii := ifc.(type) {
			case *bool2.Opt:
				t.Log("case *Opt")
				t.Log(spew.Sdump(ii.Metadata))
			case *list.Opt:
				t.Log("case *Opt")
				t.Log(spew.Sdump(ii.Metadata))
			case *float.Opt:
				t.Log("case *Opt")
				t.Log(spew.Sdump(ii.Metadata))
			case *integer.Opt:
				t.Log("case *Opt")
				t.Log(spew.Sdump(ii.Metadata))
			case *text.Opt:
				t.Log("case *Opt")
				t.Log(spew.Sdump(ii.Metadata))
			case *duration.Opt:
				t.Log("case *Opt")
				t.Log(spew.Sdump(ii.Metadata))
			default:
				// t.Log(spew.Sdump(ii))
			}
			return true
		},
	)
}

func TestMarshalUnmarshal(t *testing.T) {
	c := GetDefaultConfig()
	// d := GetDefaultConfig()
	// c.ShowAll = true
	// I.S(c)
	c.MinRelayTxFee.Set(0.3352)
	c.UUID.Set(int(rand.Int63()))
	c.Username = c.Username.Set("aoeuaoeu")
	b, e := json.MarshalIndent(c, "", "    ")
	if e != nil {
		t.Fatal(e)
	}
	t.Log("\n" + string(b))
	// c.MinRelayTxFee.Set(0.99999)
	// // t.Log("\n" + string(b))
	// c.UUID.Set(int(rand.Int63()))
	// if e = json.Unmarshal(b, d); E.Chk(e) {
	// }
	// c.UUID.Set(int(rand.Int63()))
	// c.NodeOff.Set(true)
	// c.Username = c.Username.Set("qwertyuiop")
	// c.AddPeers.Set([]string{"a", "b", "c"})
	// b, e = json.MarshalIndent(c, "", "    ")
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// t.Log("\n" + string(b))
	// c.AddPeers.Set([]string{"hello", "world"})
	// c.BanDuration.Set(69*time.Microsecond + time.Hour*5)
	// c.MinRelayTxFee.Set(0.03)
	// c.Username = c.Username.Set("123412341234")
	// // t.Log("\n" + string(b))
	// if e = json.Unmarshal(b, c); E.Chk(e) {
	// }
	// // I.S(c)
	// // d, _ := EmptyConfig()
	// c.Username = c.Username.Set("testingtesting")
	// c.NodeOff.Set(true)
	// c.MinRelayTxFee.Set(1.1)
	// b, e = json.MarshalIndent(d, "", "    ")
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// t.Log("\n" + string(b))
	// if e = json.Unmarshal(b, d); E.Chk(e) {
	// }
	// c.AddPeers.Set([]string{"one", "two", "three"})
	// c.BanDuration.Set(69*time.Millisecond + time.Second*5)
	// b, e = json.MarshalIndent(d, "", "    ")
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// t.Log("\n" + string(b))
	// b, e = json.MarshalIndent(c, "", "    ")
	// if e != nil {
	// 	t.Fatal(e)
	// }
	// c.UUID.Set(int(rand.Int63()))
	// if e = json.Unmarshal(b, c); E.Chk(e) {
	// }
	// // I.Ln("b")
	// I.S(c)
	// if e = json.Unmarshal(b, d); E.Chk(e) {
	// }
}

func TestDefaultConfig(t *testing.T) {
	c := GetDefaultConfig()
	// I.S(c)
	var e error
	var cm *opts.Command
	var depth, dist int
	var found bool
	if found, depth, dist, cm, e = c.Commands.Find("drophistory", depth, dist); !E.Chk(e) || found {
		I.F("found %v depth %d dist %d \n%s\n%v", found, depth, dist, spew.Sdump(cm), e)
	}
}
