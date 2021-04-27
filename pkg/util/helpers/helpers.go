// Package helpers provides convenience functions to simplify wallet code. This package is intended for internal wallet
// use only.
package helpers

import (
	"github.com/p9c/matrjoska/pkg/amt"
	"github.com/p9c/matrjoska/pkg/wire"
)

// SumOutputValues sums up the list of TxOuts and returns an Amount.
func SumOutputValues(outputs []*wire.TxOut) (totalOutput amt.Amount) {
	for _, txOut := range outputs {
		totalOutput += amt.Amount(txOut.Value)
	}
	return totalOutput
}

// SumOutputSerializeSizes sums up the serialized size of the supplied outputs.
func SumOutputSerializeSizes(outputs []*wire.TxOut) (serializeSize int) {
	for _, txOut := range outputs {
		serializeSize += txOut.SerializeSize()
	}
	return serializeSize
}
