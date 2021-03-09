package alo_test

import (
	"crypto/rand"
	"github.com/p9c/pod/pkg/logg"
	"testing"

	"github.com/p9c/pod/pkg/coding/alo"
	"github.com/p9c/pod/pkg/util/logi"
)

func MakeRandomBytes(size int, t *testing.T) (p []byte) {
	p = make([]byte, size)
	var e error
	if _, e = rand.Read(p); alo.err.Chk(e) {
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
				t.ftl.Ln(dataLen, size, len(s), "segments were not correctly split")
			}
		}
	}
}

func TestGetShards(t *testing.T) {
	logg.SetLogLevel("trace")
	for dataLen := 256; dataLen < 1025; dataLen += 16 {
		red := 300
		b := MakeRandomBytes(dataLen, t)
		segs := alo.GetShards(b, red)
		alo.dbg.S(segs)
		var e error
		var p *alo.Partials
		for i := range segs {
			for j := range segs[i] {
				if i == 0 && j == 0 {
					if p, e = alo.NewPacket(segs[i][j]); alo.err.Chk(e) {
						t.Fail()
					}
				} else {
					if e = p.AddShard(segs[i][j]); alo.err.Chk(e) {
						t.Fail()
					}
				}
			}
		}
		// if we got to here we should be able to decode it
		var ob []byte
		if ob, e = p.Decode(); alo.err.Chk(e) {
			t.Fail()
		}
		if string(ob) != string(b) {
			// alo.			alo.dbg.S(b)
			alo.dbg.S(ob)
			t.ftl.Ln("codec failed to decode encoded content")
		}
	}
}
