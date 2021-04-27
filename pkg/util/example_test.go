package util_test

import (
	"fmt"
	"github.com/p9c/matrjoska/pkg/amt"
	"math"
)

func ExampleAmount() {
	a := amt.Amount(0)
	fmt.Println("Zero Satoshi:", a)
	a = amt.Amount(1e8)
	fmt.Println("100,000,000 Satoshis:", a)
	a = amt.Amount(1e5)
	fmt.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Satoshi: 0 DUO
	// 100,000,000 Satoshis: 1 DUO
	// 100,000 Satoshis: 0.001 DUO
}
func ExampleNewAmount() {
	amountOne, e := amt.NewAmount(1)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(amountOne) // Output 1
	amountFraction, e := amt.NewAmount(0.01234567)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(amountFraction) // Output 2
	amountZero, e := amt.NewAmount(0)
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(amountZero) // Output 3
	amountNaN, e := amt.NewAmount(math.NaN())
	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println(amountNaN) // Output 4
	// Output: 1 DUO
	// 0.01234567 DUO
	// 0 DUO
	// invalid bitcoin amount
}
func ExampleAmount_unitConversions() {
	amount := amt.Amount(44433322211100)
	fmt.Println("Satoshi to kDUO:", amount.Format(amt.KiloDUO))
	fmt.Println("Satoshi to DUO:", amount)
	fmt.Println("Satoshi to MilliDUO:", amount.Format(amt.MilliDUO))
	fmt.Println("Satoshi to MicroDUO:", amount.Format(amt.MicroDUO))
	fmt.Println("Satoshi to Satoshi:", amount.Format(amt.Satoshi))
	// Output:
	// Satoshi to kDUO: 444.333222111 kDUO
	// Satoshi to DUO: 444333.222111 DUO
	// Satoshi to MilliDUO: 444333222.111 mDUO
	// Satoshi to MicroDUO: 444333222111 μDUO
	// Satoshi to Satoshi: 44433322211100 Satoshi
}
