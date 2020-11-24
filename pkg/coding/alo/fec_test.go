package alo_test

import (
	"crypto/rand"
	"testing"

	"github.com/p9c/pod/pkg/coding/alo"
	"github.com/p9c/pod/pkg/util/logi"
)

func MakeRandomBytes(size int, t *testing.T) (p []byte) {
	p = make([]byte, size)
	var err error
	if _, err = rand.Read(p); alo.Check(err) {
		t.Fail()
	}
	return
}

func TestSegmentBytes(t *testing.T) {
	for dataLen := 256; dataLen < 65536; dataLen += 16 {
		b := MakeRandomBytes(dataLen, t)
		for size := 32; size < 65536; size *= 2 {
			s := alo.SegmentBytes(b, size)
			if len(s) != alo.Pieces(dataLen, size) {
				t.Fatal(dataLen, size, len(s), "segments were not correctly split")
			}
		}
	}
}

func TestGetShards(t *testing.T) {
	logi.L.SetLevel("trace", false, "pod")
	for dataLen := 256; dataLen < 1025; dataLen += 16 {
		red := 300
		b := MakeRandomBytes(dataLen, t)
		segs := alo.GetShards(b, red)
		alo.Debugs(segs)
		var err error
		var p *alo.Partials
		for i := range segs {
			for j := range segs[i] {
				if i == 0 && j == 0 {
					if p, err = alo.NewPacket(segs[i][j]); alo.Check(err) {
						t.Fail()
					}
				} else {
					if err = p.AddShard(segs[i][j]); alo.Check(err) {
						t.Fail()
					}
				}
			}
		}
		// if we got to here we should be able to decode it
		var ob []byte
		if ob, err = p.Decode(); alo.Check(err) {
			t.Fail()
		}
		if string(ob) != string(b) {
			// alo.Error(err)
			alo.Debugs(b)
			alo.Debugs(ob)
			t.Fatal("codec failed to decode encoded content")
		}
	}
}
