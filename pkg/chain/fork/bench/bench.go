package main

import (
	"fmt"
	"time"

	"git.parallelcoin.io/dev/pod/pkg/chain/fork"
)

func main() {
	h := fork.Hash([]byte{}, "blake14lr", fork.List[1].ActivationHeight)
	for i := range fork.P9Algos {
		fmt.Print(`"`, i, `": {, FirstPowLimitBits, , `)
		now := time.Now().UnixNano()
		var samples int64 = 50
		for j := int64(0); j < samples; j++ {

			h = fork.Hash(h.CloneBytes(), i, fork.List[1].ActivationHeight+1)
		}
		fmt.Println((time.Now().UnixNano()-now)/samples, "},")
	}
}
