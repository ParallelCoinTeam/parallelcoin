package kopach

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/parallelcointeam/parallelcoin/pkg/chain/fork"
	"github.com/parallelcointeam/parallelcoin/pkg/util/cl"
)

// Bench is an item in a benchmark
type Bench struct {
	Name string
	Ops  int
}

// Benches is a collection of benchmarks
type Benches []Bench

// Len implements the Sorter Len method
func (b Benches) Len() int {
	return len(b)
}

// Less implements the Sorter Less method
func (b Benches) Less(i, j int) bool {
	if b[i].Ops < b[j].Ops {
		return true
	}
	return false
}

// Swap implements the Sorter Swap method
func (b Benches) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func (b Benches) getOps(a string) int {
	for i := range b {
		if b[i].Name == a {
			return b[i].Ops
		}
	}
	return 0
}

// Benchmark generates a benchmark for kopach
func Benchmark(conf string) {
	log <- cl.Info{"generating benchmarks...", cl.Ine()}
	b := [80]byte{}
	block := b[:]
	log <- cl.Warn{"initial short rough benchmark (64 rounds):"}
	initial := Benches{}
	av := fork.List[1].AlgoVers
	for i := range av {
		k := 1 << 6
		tn := time.Now()
		for j := 0; j < k; j++ {
			rand.Seed(time.Now().UnixNano())
			rand.Read(block)
			_ = fork.Hash(block, av[i], fork.List[1].ActivationHeight+2)
		}
		tnn := time.Now()
		an := fork.List[1].AlgoVers[i]
		ops := int(tnn.Sub(tn)) / k
		initial = append(initial, Bench{av[i], ops})
		pad := 15 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		log <- cl.Info{an, rightJustify(fmt.Sprint(ops/1000), 15), " us/op"}
	}
	sort.Sort(initial)
	log <- cl.Info{"initial ", initial}
	k := 1 << 9
	log <- cl.Warn{"running benchmark with ", k, " reps"}
	for i := range initial {
		tn := time.Now()
		for j := 0; j < k; j++ {
			rand.Seed(time.Now().UnixNano())
			rand.Read(block)
			_ = fork.Hash(block, initial[i].Name, fork.List[1].ActivationHeight+2)
		}
		tnn := time.Now()
		initial[i].Ops = int(tnn.Sub(tn)) / k
		an := initial[i].Name
		pad := 14 - len(an)
		if pad > 0 {
			an += strings.Repeat(" ", pad)
		}
		sort.Sort(initial)
		log <- cl.Info{an, " ", initial[i].Ops, " ns/op ", int(time.Second) / initial[i].Ops, " ops/s"}
	}
	if yp, e := json.MarshalIndent(&initial, "", "  "); e == nil {
		log <- cl.Trace{"\n", string(yp)}
		ensureDir(conf)
		if e := ioutil.WriteFile(conf, yp, 0600); e != nil {
			log <- cl.Error{"error writing ", e, cl.Ine()}
		}
	} else {
		log <- cl.Error{"error marshalling ", e, cl.Ine()}
	}
	log <- cl.Trace{"results ", initial}
	time.Sleep(time.Second)
}

// EnsureDir checks a file could be written to a path,
// creates the directories as needed
func ensureDir(fileName string) {
	dirName := filepath.Dir(fileName)
	if _, sErr := os.Stat(dirName); sErr != nil {
		mErr := os.MkdirAll(dirName, os.ModePerm)
		if mErr != nil {
			panic(mErr)
		}
	}
}

// RightJustify takes a string and right justifies it by a width or crops it
func rightJustify(s string, w int) string {
	sw := len(s)
	diff := w - sw
	if diff > 0 {
		s = strings.Repeat(" ", diff) + s
	} else if diff < 0 {
		s = s[:w]
	}
	return s
}
