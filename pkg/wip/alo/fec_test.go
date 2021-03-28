package alo

import (
	"crypto/rand"
	"github.com/p9c/log"
	"testing"
	
)

func MakeRandomBytes(size int, t *testing.T) (p []byte) {
	p = make([]byte, size)
	var e error
	if _, e = rand.Read(p); E.Chk(e) {
		t.Fail()
	}
	return
}

func TestSegmentBytes(t *testing.T) {
	for dataLen := 256; dataLen < 65536; dataLen += 16 {
		b := MakeRandomBytes(dataLen, t)
		for size := 32; size < 65536; size *= 2 {
			s := SegmentBytes(b, size)
			if len(s) != Pieces(dataLen, size) {
				t.Fatal(dataLen, size, len(s), "segments were not correctly split")
			}
		}
	}
}

func TestGetShards(t *testing.T) {
	log.SetLogLevel("trace")
	for dataLen := 256; dataLen < 1025; dataLen += 16 {
		red := 300
		b := MakeRandomBytes(dataLen, t)
		segs := GetShards(b, red)
		// alo.D.S(segs)
		var e error
		var p *Partials
		for i := range segs {
			for j := range segs[i] {
				if i == 0 && j == 0 {
					if p, e = NewPacket(segs[i][j]); E.Chk(e) {
						t.Fail()
					}
				} else {
					if e = p.AddShard(segs[i][j]); E.Chk(e) {
						t.Fail()
					}
				}
			}
		}
		// if we got to here we should be able to decode it
		var ob []byte
		if ob, e = p.Decode(); E.Chk(e) {
			t.Fail()
		}
		if string(ob) != string(b) {
			// alo.			alo.D.S(b)
			D.S(ob)
			t.Fatal("codec failed to decode encoded content")
		}
	}
}
