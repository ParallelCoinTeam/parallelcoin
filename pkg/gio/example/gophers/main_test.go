// SPDX-License-Identifier: Unlicense OR MIT

package main

import (
	"image"
	"testing"

	"github.com/p9c/pod/pkg/gio/layout"
)

func BenchmarkUI(b *testing.B) {
	fetch := func(_ string) {}
	u := newUI(fetch)
	gtx := new(layout.Context)
	for i := 0; i < b.N; i++ {
		gtx.Reset(nil, image.Point{800, 600})
		u.Layout(gtx)
	}
}
