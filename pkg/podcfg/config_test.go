package podcfg

import (
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"testing"
)

func TestForEach(t *testing.T) {
	c, _ := EmptyConfig()
	c.ForEach(func(ifc interface{}) bool {
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
	})
}

func TestMarshalUnmarshal(t *testing.T) {
	c, _ := EmptyConfig()
	b, e := json.MarshalIndent(&c,"","    ")
	if e != nil {
		t.Fatal(e)
	}
	t.Log(string(b))
}