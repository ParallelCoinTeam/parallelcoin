package kopach

import (
	"context"
	"sync"
	"time"

	"github.com/p9c/rpcx/client"

	"github.com/p9c/pod/pkg/chain/mining"
	"github.com/p9c/pod/pkg/log"
	"github.com/p9c/pod/pkg/scratch/kcpx"
	"github.com/p9c/pod/pkg/util"
)

// Subscriber is the information required to track a subscriber
type Subscriber struct {
	X       client.XClient
	Timeout time.Time
}

// Subscribers is a map of addresses and deadlines.
// A worker thread will check them on a ticker and remove the outdated items,
// the rest get sent blocks
type Subscribers map[string]Subscriber

// SubmitFunc is a function that accepts a block and returns a string result
// of the processing of the block
type SubmitFunc func(b *util.Block) string

// Controller keeps a subscriber list and has a block submit function for
// when workers find a solution
type Controller struct {
	sync.Mutex
	password    string
	Subscribers *Subscribers
	timeout     time.Duration
	submit      SubmitFunc
	lastBlock   *[]mining.BlockTemplate
	quit        chan struct{}
}

// NewController creates a new controller
// - blockChan needs to be connected to a function that relays a block template
// for each algorithm in a slice when a new block connects.
// It should send an empty slice to indicate stop work when chain is out of
// date or node is disconnected
// - submitFunc needs to be set to a function that submits the solved blocks
// sent back by workers to the node for processing
func NewController(password string, submit SubmitFunc,
	blockChan chan []mining.BlockTemplate, t time.Duration) (c *Controller,
	shutdown func(), done <-chan struct{}) {
	subs := make(Subscribers)
	c = &Controller{password: password, timeout: t,
		Subscribers: &subs, submit: submit}
	ticker := time.NewTicker(time.Second)
	quit := make(chan struct{})
	go func() {
		select {
		case <-ticker.C:
			c.Purge()
			c.Lock()
			// TODO: lastBlock could be updated here saving miners time adjusting
			//  timestamp (ahaha)
			c.SendBlock(*c.lastBlock)
			c.Unlock()
		case b := <-blockChan:
			c.Lock()
			c.lastBlock = &b
			c.Unlock()
			c.SendBlock(b)
		case <-quit:
		}
	}()
	shutdown = func() { close(quit) }
	return
}

// CopySubscribers creates a new slice with all of the XClients that have not
// yet timed out. These should be used to send new Block messages, presumably
func (c *Controller) CopySubscribers() (subs map[string]client.XClient) {
	c.Lock()
	subs = make(map[string]client.XClient)
	for i := range *c.Subscribers {
		subs[i] = (*c.Subscribers)[i].X
	}
	c.Lock()
	return
}

// Purge deletes all entries past timeout
func (c *Controller) Purge() {
	c.Lock()
	var deleters []string
	tn := time.Now()
	for i := range *c.Subscribers {
		if tn.Sub((*c.Subscribers)[i].Timeout) < 0 {
			deleters = append(deleters, i)
		}
	}
	for i := range deleters {
		delete(*c.Subscribers, deleters[i])
	}
	c.Unlock()
}

// SendBlock sends a block out to all the Subscribers
func (c *Controller) SendBlock(b []mining.BlockTemplate) {
	subs := c.CopySubscribers()
	for i := range subs {
		if subs[i] == nil {
			subs[i] = kcpx.NewXClient(i, "Kopach", c.password)
			// update the controller for this case
			c.Lock()
			if cs, ok := (*c.Subscribers)[i]; ok {
				cs.X = subs[i]
				c.Unlock()
			} else {
				c.Unlock()
				// it disappeared, so don't run it (must have unsubscribed)
				continue
			}
		}
		receivedTime := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := subs[i].Call(ctx, "Block", &b, &receivedTime)
		if err != nil {
		log.ERROR(err)
log.ERROR("error sending block ", err)
			continue
		}
		go func() {
			<-ctx.Done()
			log.WARN("worker ", i, " reported receiving ", receivedTime)
			cancel()
		}()
	}
}

// Subscribe adds a worker to the subscriber list and returns the time at
// which it will be removed from the subscriber list if it does not
// beforehand send a subsequent subscription message
func (c *Controller) Subscribe(ctx context.Context, args *string,
	reply *time.Time) (err error) {
	tn := time.Now()
	*reply = tn
	c.Lock()
	(*c.Subscribers)[*args] = Subscriber{Timeout: tn.Add(c.timeout)}
	c.Unlock()
	return
}

// Unsubscribe removes a worker from the subscriber list,
// returns 1 if it *was* subscribed and thus removed, or 0
func (c *Controller) Unsubscribe(ctx context.Context, args *string,
	reply *int) (err error) {
	c.Lock()
	if _, ok := (*c.Subscribers)[*args]; ok {
		delete(*c.Subscribers, *args)
		*reply = 1
	} else {
		*reply = 0
	}
	c.Lock()
	return
}

// Submit calls the submit function on the provided block and returns a
// message indicating the result (connected, orphan, error)
func (c *Controller) Submit(ctx context.Context, args *util.Block,
	reply *string) (err error) {
	*reply = c.submit(args)
	return
}
