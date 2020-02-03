// SPDX-License-Identifier: Unlicense OR MIT

package ops

import (
	"encoding/binary"
	"math"

	"github.com/p9c/pod/pkg/gui/f32"
	"github.com/p9c/pod/pkg/gui/internal/opconst"
	"github.com/p9c/pod/pkg/gui/op"
)

func DecodeTransformOp(d []byte) op.TransformOp {
	bo := binary.LittleEndian
	if opconst.OpType(d[0]) != opconst.TypeTransform {
		panic("invalid op")
	}
	return op.TransformOp{}.Offset(f32.Point{
		X: math.Float32frombits(bo.Uint32(d[1:])),
		Y: math.Float32frombits(bo.Uint32(d[5:])),
	})
}
