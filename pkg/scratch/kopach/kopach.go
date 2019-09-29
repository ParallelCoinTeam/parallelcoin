package kopach

import (
	"context"
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/p9c/rpcx/client"

	"github.com/p9c/pod/pkg/chain/config/netparams"
	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/scratch/kcpx"
	"github.com/p9c/pod/pkg/util"
)

// MineFunc is a function that stops on the semaphore closing and ostensibly
// mines the block provided
type MineFunc func(submit func(*util.Block) string, semaphore chan struct{},
	b []mining.BlockTemplate)

// Kopach is a worker
type Kopach struct {
	sync.Mutex
	params      *netparams.Params
	service     string
	group       string
	X           client.XClient
	password    string
	Controllers *[]string
	current     int
	Mine        MineFunc
	Semaphore   chan struct{}
	quit        chan struct{}
}

// NewKopach returns a new worker loaded with a mining function.
// - shutdown() stops the miner
// - done unblocks and returns nil when shutdown is complete
// - MineFunc hashes blocks and calls submit when it finds a solution,
// stops when the semaphore channel is closed
func NewKopach(service, group, address, password string, controllers []string,
	m MineFunc, nodiscovery bool, activeNet *netparams.Params) (k *Kopach,
	shutdown func(),
	done <-chan struct{}) {
	k = &Kopach{
		params:      activeNet,
		service:     service,
		group:       group,
		Mine:        m,
		password:    password,
		Controllers: &controllers,
		Semaphore:   make(chan struct{}),
		quit:        make(chan struct{}),
	}
	d := make(chan struct{})
	ticker := time.NewTicker(time.Second)
	go func() {
		_, stopServer := kcpx.Serve(address, "Kopach", password, k)
		select {
		case <-ticker.C:
			k.Lock()
			l := len(*k.Controllers)
			k.Unlock()
			if l > 0 {
				k.Lock()
				if k.X == nil {
					k.X = kcpx.NewXClient((*k.Controllers)[rand.Intn(len(*k.
						Controllers))], "Controller", k.password)
				}
				k.Unlock()
				deadline := time.Now()
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				k.Lock()
				rn := rand.Intn(len(*k.Controllers)) - 1
				if nodiscovery {
					rn = k.current
				}
				k.Unlock()
				kc := (*k.Controllers)[rn]
				err := k.X.Call(ctx, "Subscribe", kc, &deadline)
				if err != nil {
					log.ERROR("error sending block ", err)
					if nodiscovery {
						// in nodiscovery mode we roll to the next on failure
						k.Lock()
						k.current++
						if k.current > len(*k.Controllers) {
							k.current = 0
						}
						k.Unlock()
					} else {
						// in discovery mode failed controller is removed from
						// list, as when it returns it will be readded
						var out []string
						k.Lock()
						for i := range *k.Controllers {
							if i != rn {
								out = append(out, (*k.Controllers)[i])
							}
						}
						*k.Controllers = out
						k.Lock()
					}
				}
				go func() {
					<-ctx.Done()
					log.WARN("subscription accepted, expires ", deadline, )
					cancel()
				}()
			}
		case <-k.quit:
		}
		// stop the current job
		close(k.Semaphore)
		<-stopServer()
		close(d)
	}()
	shutdown = func() { close(k.quit) }
	done = d
	return
}

// Block delivers a new block to a Kopach
func (k *Kopach) Block(ctx context.Context, args *[]mining.BlockTemplate,
	reply *time.Time) (err error) {
	if *args == nil {
		log.WARN("empty block means don't work")
		return errors.New("empty block, not working :)")
	}
	defer func() {
		// most likely panic because close of closed channel, ignore
		if r := recover(); r != nil {
			log.DEBUG("Recovered in f", r)
			err = errors.New("worker is busy in submit")
		}
		// receiving a block while submitting usually will lead here,
		// so it is good that after this defer nothing further happens.
	}()
	*reply = time.Now()
	// kill the previous worker
	close(k.Semaphore)
	// replace the semaphore
	k.Lock()
	k.Semaphore = make(chan struct{})
	k.Unlock()
	// start the new job
	k.Mine(k.Submit, k.Semaphore, *args)
	return
}

// Submit sends out a solved block to the controller
func (k *Kopach) Submit(b *util.Block) (reply string) {
	defer func() {
		// blocks will now be processed again
		k.Lock()
		k.Semaphore = make(chan struct{})
		k.Unlock()
	}()
	// stop work during submission - any block received until we get a reply
	// from the controller or this thread dies will be ignored
	close(k.Semaphore)
	if k.X == nil {
		k.Lock()
		k.X = kcpx.NewXClient((*k.Controllers)[rand.Intn(len(*k.Controllers))],
			"Controller", k.password)
		k.Unlock()
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second/2)
	for try := 0; try < 3; try++ {
		err := k.X.Call(ctx, "Submit", &b, &reply)
		if err != nil {
			log.ERROR("error sending block ", err)
			return err.Error()
		}
	}
	<-ctx.Done()
	log.WARN("controller replied ", reply, " to block submit", )
	cancel()
	return reply
}
