package cfgutil

import (
	"strconv"
	"strings"

	"github.com/p9c/pod/pkg/util"
)

// AmountFlag embeds a util.Amount and implements the flags.Marshaler and
// Unmarshaler interfaces so it can be used as a config struct field.
type AmountFlag struct {
	util.Amount
}

// NewAmountFlag creates an AmountFlag with a default util.Amount.
func NewAmountFlag(	defaultValue util.Amount) *AmountFlag {
	return &AmountFlag{defaultValue}
}

// MarshalFlag satisifes the flags.Marshaler interface.
func (a *AmountFlag) MarshalFlag() (string, error) {
	return a.Amount.String(), nil
}

// UnmarshalFlag satisifes the flags.Unmarshaler interface.
func (a *AmountFlag) UnmarshalFlag(value string) error {
	value = strings.TrimSuffix(value, " DUO")
	valueF64, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return err
	}
	amount, err := util.NewAmount(valueF64)
	if err != nil {
		return err
	}
	a.Amount = amount
	return nil
}
