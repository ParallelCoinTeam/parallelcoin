package util_test

import (
	"github.com/p9c/pod/pkg/log"
	"math"

	"github.com/p9c/pod/pkg/util"
)

func ExampleAmount() {
	a := util.Amount(0)
	log.Println("Zero Satoshi:", a)
	a = util.Amount(1e8)
	log.Println("100,000,000 Satoshis:", a)
	a = util.Amount(1e5)
	log.Println("100,000 Satoshis:", a)
	// Output:
	// Zero Satoshi: 0 DUO
	// 100,000,000 Satoshis: 1 DUO
	// 100,000 Satoshis: 0.001 DUO
}
func ExampleNewAmount() {
	amountOne, err := util.NewAmount(1)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(amountOne) //Output 1
	amountFraction, err := util.NewAmount(0.01234567)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(amountFraction) //Output 2
	amountZero, err := util.NewAmount(0)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(amountZero) //Output 3
	amountNaN, err := util.NewAmount(math.NaN())
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(amountNaN) //Output 4
	// Output: 1 DUO
	// 0.01234567 DUO
	// 0 DUO
	// invalid bitcoin amount
}
func ExampleAmount_unitConversions() {
	amount := util.Amount(44433322211100)
	log.Println("Satoshi to kDUO:", amount.Format(util.AmountKiloDUO))
	log.Println("Satoshi to DUO:", amount)
	log.Println("Satoshi to MilliDUO:", amount.Format(util.AmountMilliDUO))
	log.Println("Satoshi to MicroDUO:", amount.Format(util.AmountMicroDUO))
	log.Println("Satoshi to Satoshi:", amount.Format(util.AmountSatoshi))
	// Output:
	// Satoshi to kDUO: 444.333222111 kDUO
	// Satoshi to DUO: 444333.222111 DUO
	// Satoshi to MilliDUO: 444333222.111 mDUO
	// Satoshi to MicroDUO: 444333222111 μDUO
	// Satoshi to Satoshi: 44433322211100 Satoshi
}
