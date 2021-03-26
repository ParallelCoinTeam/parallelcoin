package podcfg

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"math/rand"
	"testing"
	"time"
)

func TestForEach(t *testing.T) {
	c := New()
	c.ForEach(
		func(ifc interface{}) bool {
			switch ii := ifc.(type) {
			case *Bool:
				t.Log("case *Bool")
				t.Log(spew.Sdump(ii.metadata))
			case *Strings:
				t.Log("case *Strings")
				t.Log(spew.Sdump(ii.metadata))
			case *Float:
				t.Log("case *Float")
				t.Log(spew.Sdump(ii.metadata))
			case *Int:
				t.Log("case *Int")
				t.Log(spew.Sdump(ii.metadata))
			case *String:
				t.Log("case *String")
				t.Log(spew.Sdump(ii.metadata))
			case *Duration:
				t.Log("case *Duration")
				t.Log(spew.Sdump(ii.metadata))
			default:
				t.Log(spew.Sdump(ii))
			}
			return true
		},
	)
}

func TestMarshalUnmarshal(t *testing.T) {
	c := New()
	d := New()
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
	c.MinRelayTxFee.Set(0.99999)
	// t.Log("\n" + string(b))
	c.UUID.Set(int(rand.Int63()))
	if e = json.Unmarshal(b, d); E.Chk(e) {
	}
	c.UUID.Set(int(rand.Int63()))
	c.NodeOff.Set(true)
	c.Username = c.Username.Set("qwertyuiop")
	c.AddPeers.Set([]string{"a", "b", "c"})
	b, e = json.MarshalIndent(c, "", "    ")
	if e != nil {
		t.Fatal(e)
	}
	t.Log("\n" + string(b))
	c.AddPeers.Set([]string{"hello","world"})
	c.BanDuration.Set(69*time.Microsecond + time.Hour*5)
	c.MinRelayTxFee.Set(0.03)
	c.Username = c.Username.Set("123412341234")
	// t.Log("\n" + string(b))
	if e = json.Unmarshal(b, c); E.Chk(e) {
	}
	// I.S(c)
	// d, _ := EmptyConfig()
	c.Username = c.Username.Set("testingtesting")
	c.NodeOff.Set(true)
	c.MinRelayTxFee.Set(1.1)
	b, e = json.MarshalIndent(d, "", "    ")
	if e != nil {
		t.Fatal(e)
	}
	t.Log("\n" + string(b))
	if e = json.Unmarshal(b, d); E.Chk(e) {
	}
	c.AddPeers.Set([]string{"one", "two", "three"})
	c.BanDuration.Set(69*time.Millisecond + time.Second*5)
	b, e = json.MarshalIndent(d, "", "    ")
	if e != nil {
		t.Fatal(e)
	}
	t.Log("\n" + string(b))
	b, e = json.MarshalIndent(c, "", "    ")
	if e != nil {
		t.Fatal(e)
	}
	c.UUID.Set(int(rand.Int63()))
	if e = json.Unmarshal(b, c); E.Chk(e) {
	}
	// I.Ln("b")
	I.S(c)
	// if e = json.Unmarshal(b, d); E.Chk(e) {
	// }
}
