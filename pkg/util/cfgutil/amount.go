package cfgutil

import (
	"github.com/p9c/matrjoska/pkg/amt"
	"strconv"
	"strings"
)

// AmountFlag embeds a util.Amount and implements the flags.Marshaler and Unmarshaler interfaces so it can be used as a
// config struct field.
type AmountFlag struct {
	amt.Amount
}

// NewAmountFlag creates an AmountFlag with a default util.Amount.
func NewAmountFlag(defaultValue amt.Amount) *AmountFlag {
	return &AmountFlag{defaultValue}
}

// MarshalFlag satisifes the flags.Marshaler interface.
func (a *AmountFlag) MarshalFlag() (string, error) {
	return a.Amount.String(), nil
}

// UnmarshalFlag satisifes the flags.Unmarshaler interface.
func (a *AmountFlag) UnmarshalFlag(value string) (e error) {
	value = strings.TrimSuffix(value, " DUO")
	valueF64, e := strconv.ParseFloat(value, 64)
	if e != nil {
		E.Ln(e)
		return e
	}
	amount, e := amt.NewAmount(valueF64)
	if e != nil {
		E.Ln(e)
		return e
	}
	a.Amount = amount
	return nil
}
