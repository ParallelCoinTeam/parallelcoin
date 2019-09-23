// +build ignore

package kopach

import (
	"context"
	"crypto/sha1"
	"fmt"
	"math/big"
	"math/rand"
	"net"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/p9c/rpcx/client"
	"github.com/p9c/rpcx/server"
	"github.com/dchest/blake256"
	"golang.org/x/crypto/pbkdf2"

	"github.com/p9c/pod/cmd/node/state"
	blockchain "github.com/p9c/pod/pkg/chain"
	"github.com/p9c/pod/pkg/chain/fork"
	"github.com/p9c/pod/pkg/chain/mining"
	txscript "github.com/p9c/pod/pkg/chain/tx/script"
	"github.com/p9c/pod/pkg/chain/wire"
	"github.com/p9c/pod/pkg/pod"
	"github.com/p9c/pod/pkg/util"
)

const (
	// // maxNonce is the maximum value a nonce can be in a block header.
	// maxNonce = 2 ^ 32 - 1
	// maxExtraNonce is the maximum value an extra nonce used in a coinbase
	// transaction can be.
	maxExtraNonce = 2 ^ 64 - 1
	CoinbaseFlags = "/P2SH/pod/"
)

// Ping is the ping message
type Ping struct {
	Time   time.Time
	Sender string
}

// Pong is the pong message
type Pong struct {
	Time time.Time
}

// Block is a collection of block templates
type Block struct {
	Templates []mining.BlockTemplate
}

// Subscriber is stored in a sync.Map with the address as the key
type Subscriber struct {
	LastPing time.Time
	XClient  client.XClient
}

// Controller is a service for the node that delivers work over rpcx
//
type Controller struct {
	BC                     *blockchain.BlockChain
	Cfg                    *pod.Config
	StateCfg               *state.Config
	Subscribers            sync.Map
	Blocks                 atomic.Value
	BlockTemplateGenerator *mining.BlkTmplGenerator
	ProcessBlock           func(block *util.Block, flags blockchain.BehaviorFlags) (bool, error)
	ConnectedCount         func() int32
	IsCurrent              func() bool
	Quit                   chan struct{}
}

type selections []selection

var _ sort.Interface = selections{}

func (ss selections) Len() int           { return len(ss) }
func (ss selections) Less(i, j int) bool { return ss[i].ops.Cmp(ss[j].ops) < 0 }
func (ss selections) Swap(i, j int)      { ss[i], ss[j] = ss[j], ss[i] }

type selection struct {
	algo  string
	ops   *big.Int
	coeff uint64
	block *mining.BlockTemplate
}

// StartController starts up a kopach miner controller
func StartController(c *Controller) {
	log.WARN("starting kopach miner controller")
	address := *c.Cfg.Controller
	password := *c.Cfg.MinerPass
	if address == "" {
		return
	}
	s := NewKCPService(address, password)
	s.DisableHTTPGateway = true
	initialTemplates := Block{}
	log.WARN("creating initial blocks")
	for i := range fork.List[fork.GetCurrent(c.BC.BestSnapshot().Height+1)].Algos {
		// Choose a payment address at random.
		//
		rand.Seed(time.Now().UnixNano())
		payToAddr := c.StateCfg.ActiveMiningAddrs[rand.Intn(len(*c.Cfg.MiningAddrs))]
		if bt, err := c.BlockTemplateGenerator.NewBlockTemplate(payToAddr, i); err == nil {
			initialTemplates.Templates = append(initialTemplates.Templates, *bt)
		}
	}
	c.Blocks.Store(initialTemplates)
	if err := s.RegisterName("Controller", c, ""); err != nil {
		log.ERROR("failed to register controller ", err)
	}
	c.BC.Subscribe(func(ntfn *blockchain.Notification) {
		INFO("kopach controller new block")
		if ntfn.Type == blockchain.NTBlockConnected {
			if block, ok := ntfn.Data.(*util.Block); !ok {
				log.WARN("chain connected notification is not a block")
			} else {
				log.WARN("new block height ", block.Height())
				tmpl := Block{}
				for i := range fork.List[fork.GetCurrent(block.Height()+1)].Algos {
					// Choose a payment address at random.
					rand.Seed(time.Now().UnixNano())
					payToAddr := c.StateCfg.ActiveMiningAddrs[rand.Intn(len(*c.Cfg.MiningAddrs))]
					if bt, err := c.BlockTemplateGenerator.NewBlockTemplate(payToAddr, i); err == nil {
						log.TRACE("created template for ", i)
						tmpl.Templates = append(tmpl.Templates, *bt)
					}
				}
				// atomically store the block template set for first initial ping response
				log.TRACE("storing templates in atomic")
				c.Blocks.Store(tmpl)
				// prune current subscribers list
				var delKeys []string
				c.Subscribers.Range(func(key, value interface{}) bool {
					if time.Now().Sub(value.(Subscriber).LastPing) > time.Second*5 {
						// note which keys are stale
						delKeys = append(delKeys, key.(string))
					}
					return true
				})
				// delete the stale keys
				log.TRACE("deleting timed out subscribers ", len(delKeys))
				for _, x := range delKeys {
					c.Subscribers.Delete(x)
				}
				if c.ConnectedCount() > 1 && c.IsCurrent() {
					sendXC := make(map[string]Subscriber)
					c.Subscribers.Range(func(key, value interface{}) bool {
						sub := value.(Subscriber)
						sendXC[key.(string)] = sub
						return true
					})
					for i, x := range sendXC {
						go func() {
							log.WARN("sending to subscriber ", i)
							done := make(chan *client.Call, 2)
							pong := &Pong{}
							_, err := x.XClient.Go(context.Background(), "Block", &tmpl, &pong, done)
							if err != nil {
								log.ERROR("error calling Worker Block ", err)
								select {
								case <-c.Quit:
								case <-time.After(time.Millisecond * 100):
								case <-done:
								}
							}
						}()
					}
				}
			}
		}
	})
	// server goroutine
	go func() {
		log.WARN("serving Controller")
		if err := s.Serve("kcp", address); err != nil {
			log.DEBUG("error serving Controller ", err)
		}
	}()
}

// Ping is an 'are you alive' query, this does double duty by prompting the sending of a new block template
func (c *Controller) Ping(ctx context.Context, args *Ping, reply *Pong) error {
	c.Subscribers.Range(func(key, value interface{}) bool {
		if key.(string) == args.Sender {
			// update timestamp if found
			//
			value = time.Now()
			return false
		}
		return true
	})
	xc := NewKCPConnection("Worker", args.Sender, *c.Cfg.MinerPass)
	c.Subscribers.Store(args.Sender, Subscriber{LastPing: time.Now(), XClient: xc})
	done := make(chan *client.Call, 2)
	pong := &Pong{}
	block := c.Blocks.Load()
	_, err := xc.Go(context.Background(), "Block", &block, &pong, done)
	if err != nil {
		log.ERROR("error calling Worker Block ", err)
		c.Subscribers.Delete(args.Sender)
	}
	select {
	// case <-time.After(time.Millisecond * 100):
	// 	log.WARN("timeout"}
	case <-done:
		log.WARN("successfully sent block to worker")
	}
	log.TRACE("received ping ", args.Sender)
	*reply = Pong{Time: time.Now()}
	return nil
}

func (c *Controller) Submit(ctx context.Context, args *wire.MsgBlock, reply *string) error {
	fmt.Println(*args)
	log.WARN("submitting block")
	isOrphan, err := c.ProcessBlock(util.NewBlock(args), blockchain.BFNone)
	if err != nil {
		// Anything other than a rule violation is an unexpected error, so log
		// that error as an internal error.
		if _, ok := err.(blockchain.RuleError); !ok {
			*reply = "Unexpected error while processing block submitted via CPU miner: " + err.Error()
		} else {
			*reply = "block submitted via CPU miner rejected: " + err.Error()
		}
	}
	if isOrphan {
		*reply = "block is an orphan, sorry, try again next time"
	} else {
		*reply = "the block was accepted, you win!"
	}
	return nil
}

// Worker is the service that gets work for a miner
//
type Worker struct {
	XClient client.XClient
	newWork chan struct{}
	benches benches
	bias    int
	threads int
	quit    chan struct{}
}

// NewWorker creates a new worker
func NewWorker(controller, password, listener, dataDir string, bias, threads int, quit chan struct{}) {
	if bias > 8 {
		bias = 8
	}
	if bias < -8 {
		bias = -8
	}
	c := NewKCPConnection("Controller", controller, password)
	w := &Worker{XClient: c, bias: bias, threads: threads, quit: quit}
	conf := filepath.Join(dataDir, "bench.json")
	log.WARN("conf ", conf)
	if !FileExists(conf) {
		log.WARN("run benchmark")
		Benchmark(conf)
	}
	// load benchmark
	//
	w.benches = benches{}
	if b, e := ioutil.ReadFile(conf); e != nil {
		log.ERROR("error reading benches ", e)
	} else {
		e = json.Unmarshal(b, &w.benches)
		if e != nil {
			log.ERROR("error unmarshaling ", e)
		}
	}
	var counter atomic.Value
	counter.Store(0)
	// ping loop keeps server's stamp fresh
	go func() {
	out:
		for {
			if counter.Load().(int) > 5 {
				if err := c.Close(); err != nil {
					log.ERROR("error closing XClient ", err)
				}
				c = NewKCPConnection("Controller", controller, password)
			}
			tn := Ping{Time: time.Now(), Sender: listener}
			tnp := Pong{Time: tn.Time}
			done := make(chan *client.Call, 2)
			_, err := c.Go(context.Background(), "Ping", &tn, &tnp, done)
			if err != nil {
				log.ERROR("error calling Controller Ping ", err)
			}
			select {
			case <-quit:
				break out
			case <-time.After(time.Second):
				counter.Store(counter.Load().(int) + 1)
			case p := <-done:
				counter.Store(0)
				log.TRACE("pong ", time.Now().Sub(p.Reply.(*Pong).Time))
				time.Sleep(time.Second)
			}
		}
		if err := c.Close(); err != nil {
			log.ERROR("error closing connection ", err)
		}
	}()
	// start up the worker rpc server
	//
	s := NewKCPService(listener, password)
	if err := s.RegisterName("Worker", w, ""); err != nil {
		log.ERROR("failed to register controller ", err)
	}
	go func() {
		log.WARN("serving Controller")
		if err := s.Serve("kcp", listener); err != nil {
			log.DEBUG("error serving Controller ", err)
		}
	}()
	<-quit
}

// Block is the worker RPC method to give it a new set of block templates
func (w *Worker) Block(ctx context.Context, args *Block, reply *Pong) error {
	if w.newWork == nil {
		w.newWork = make(chan struct{})
	} else {
		w.newWork <- struct{}{}
	}
	// signal previous workers to stop
	*reply = Pong{Time: time.Now()}
	block := (*args).Templates
	height := block[0].Height
	var se selections
	for i := range block {
		se = append(se, selection{
			algo:  fork.GetAlgoName(block[i].Block.Header.Version, height),
			block: &block[i],
		})
	}
	for i := range se {
		ops := w.benches.getOps(se[i].algo)
		bOps := big.NewInt(int64(ops))
		se[i].ops = big.NewInt(1).Mul(bOps, fork.CompactToBig(block[i].Block.Header.Bits))
	}
	sort.Sort(se)
	for i := range se {
		coeff := big.NewInt(1).Div(se[i].ops, se[0].ops)
		se[i].coeff = coeff.Uint64()
		log.DEBUG(se[i].algo, " ", se[i].coeff)
	}

	var out selections
	switch {
	case w.bias < 0:
		out = se[-w.bias:]
	case w.bias > 0:
		out = se[:9-w.bias]
	default:
		out = se
	}
	choice := out[rand.Intn(len(out))]
	log.WARN("selection ", choice.algo, " ", choice.coeff, " threads ",
		w.threads)
	// solve the block!
	for i := 0; i < w.threads; i++ {
		go func() {
		outest:
			for {
				log.WARN("solving block thread ", i)
			sb := choice.block.Block
			targetDifficulty := fork.CompactToBig(sb.Header.Bits)
			en, _ := wire.RandomUint64()
			if err := UpdateExtraNonce(sb, height, en); err != nil {
				log.ERROR("error updating extra nonce ", err)
		}
		rn, _ := wire.RandomUint64()
		rNonce := uint32(rn)
		mn := uint32(1 << 16) // typically the timestamp will update before one core exhausts this
		for i := rNonce; i <= rNonce+mn; i++ {
			sb.Header.Nonce = i
			hash := sb.Header.BlockHashWithAlgos(height)
			bigHash := blockchain.HashToBig(&hash)
			if bigHash.Cmp(targetDifficulty) <= 0 {
				// submit to controller
				//
				var result string
				done := make(chan *client.Call, 2)
				// solved := util.NewBlock(sb)
				log.ERROR("yay a block! ", sb.Header.BlockHashWithAlgos(
					height))
			submission := *sb
			_, err := w.XClient.Go(context.Background(), "Submit", &submission, &result, done)
			if err != nil {
				log.ERROR("error calling Submit ", err)
		}
		select {
		case <-w.quit:
			break outest
		case p := <-done:
			log.WARN("result of submit: ", *p.Reply.(*string))
	}
}
select {
case <-w.newWork:
break outest
default:
}
}
}
}()
}
return nil
}

// NewKCPService creates a new KCP service with an encryption salt/password
func NewKCPService(address, password string) *server.Server {
	return server.NewServer(server.WithBlockCrypt(getBC(password)))
}

// NewKCPConnection creates a new encrypted KCP connection
func NewKCPConnection(service, address, password string) client.XClient {
	option := client.DefaultOption
	option.Block = getBC(password)
	d := client.NewPeer2PeerDiscovery("kcp@"+address, "")
	xClient := client.NewXClient(service, client.Failtry, client.RoundRobin, d, option)
	cs := &configUDPSession{}
	pc := client.NewPluginContainer()
	pc.Add(cs)
	xClient.SetPlugins(pc)
	return xClient
}

func getBC(password string) kcp.BlockCrypt {
	h := blake256.Sum256([]byte(password))
	pass := pbkdf2.Key([]byte(password), h[:], 4096, 32, sha1.New)
	bc, _ := kcp.NewAESBlockCrypt(pass)
	return bc
}

type configUDPSession struct{}

func (p *configUDPSession) ConnCreated(conn net.Conn) (net.Conn, error) {
	session, ok := conn.(*kcp.UDPSession)
	if !ok {
		return conn, nil
	}

	session.SetACKNoDelay(true)
	session.SetStreamMode(true)
	return conn, nil
}

// FileExists reports whether the named file or directory exists.
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return err == nil
}

// RightJustify takes a string and right justifies it by a width or crops it
func RightJustify(s string, w int) string {
	sw := len(s)
	diff := w - sw
	if diff > 0 {
		s = strings.Repeat(" ", diff) + s
	} else if diff < 0 {
		s = s[:w]
	}
	return s
}

// standardCoinbaseScript returns a standard script suitable for use as the
// signature script of the coinbase transaction of a new block.  In particular,
// it starts with the block height that is required by version 2 blocks and
// adds the extra nonce as well as additional coinbase flags.
func standardCoinbaseScript(nextBlockHeight int32, extraNonce uint64) ([]byte, error) {
	return txscript.NewScriptBuilder().AddInt64(int64(nextBlockHeight)).
		AddInt64(int64(extraNonce)).AddData([]byte(CoinbaseFlags)).
		Script()
}

func UpdateExtraNonce(msgBlock *wire.MsgBlock, blockHeight int32, extraNonce uint64) error {
	coinbaseScript, err := standardCoinbaseScript(blockHeight, extraNonce)
	if err != nil {
		return err
	}
	if len(coinbaseScript) > blockchain.MaxCoinbaseScriptLen {
		return fmt.Errorf(
			"coinbase transaction script length of %d is out of range (min: %d, max: %d)",
			len(coinbaseScript), blockchain.MinCoinbaseScriptLen,
			blockchain.MaxCoinbaseScriptLen)
	}
	msgBlock.Transactions[0].TxIn[0].SignatureScript = coinbaseScript
	// TODO(davec): A util.Block should use saved in the state to avoid
	// recalculating all of the other transaction hashes.
	// block.Transactions[0].InvalidateCache() Recalculate the merkle root with
	// the updated extra nonce.
	block := util.NewBlock(msgBlock)
	merkles := blockchain.BuildMerkleTreeStore(block.Transactions(), false)
	msgBlock.Header.MerkleRoot = *merkles[len(merkles)-1]
	return nil
}
