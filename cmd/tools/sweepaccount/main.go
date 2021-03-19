// Copyright (c) 2015-2016 The btcsuite developers

package main

import (
	"errors"
	"fmt"
	"github.com/p9c/pod/pkg/logg"
	"io/ioutil"
	"os"
	"path/filepath"
	
	"github.com/p9c/pod/pkg/util/qu"
	
	"github.com/jessevdk/go-flags"
	"golang.org/x/crypto/ssh/terminal"
	
	"github.com/p9c/pod/app/appdata"
	"github.com/p9c/pod/pkg/blockchain/chaincfg/netparams"
	"github.com/p9c/pod/pkg/blockchain/chainhash"
	"github.com/p9c/pod/pkg/blockchain/tx/txauthor"
	"github.com/p9c/pod/pkg/blockchain/tx/txrules"
	"github.com/p9c/pod/pkg/blockchain/tx/txscript"
	"github.com/p9c/pod/pkg/blockchain/wire"
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/rpcclient"
	"github.com/p9c/pod/pkg/util"
	"github.com/p9c/pod/pkg/util/cfgutil"
)

var (
	walletDataDirectory = appdata.Dir("mod", false)
	newlineBytes        = []byte{'\n'}
)

func fatalf(
	format string, args ...interface{},
) {
	var e error
	if _, e = fmt.Fprintf(os.Stderr, format, args...); E.Chk(e) {
	}
	if _, e = os.Stderr.Write(newlineBytes); E.Chk(e) {
	}
	os.Exit(1)
}
func errContext(
	ee error, context string,
) (e error) {
	return fmt.Errorf("%s: %v", context, ee)
}

// Flags.
var opts = struct {
	TestNet3              bool                `long:"testnet" description:"Use the test bitcoin network (version 3)"`
	SimNet                bool                `long:"simnet" description:"Use the simulation bitcoin network"`
	RPCConnect            string              `short:"c" long:"connect" description:"Hostname[:port] of wallet RPC server"`
	RPCUsername           string              `short:"u" long:"rpcuser" description:"Wallet RPC username"`
	RPCCertificateFile    string              `long:"cafile" description:"Wallet RPC TLS certificate"`
	FeeRate               *cfgutil.AmountFlag `long:"feerate" description:"Transaction fee per kilobyte"`
	SourceAccount         string              `long:"sourceacct" description:"Account to sweep outputs from"`
	DestinationAccount    string              `long:"destacct" description:"Account to send sweeped outputs to"`
	RequiredConfirmations int64               `long:"minconf" description:"Required confirmations to include an output"`
}{
	TestNet3:              false,
	SimNet:                false,
	RPCConnect:            "localhost",
	RPCUsername:           "",
	RPCCertificateFile:    filepath.Join(walletDataDirectory, "rpc.cert"),
	FeeRate:               cfgutil.NewAmountFlag(txrules.DefaultRelayFeePerKb),
	SourceAccount:         "imported",
	DestinationAccount:    "default",
	RequiredConfirmations: 1,
}

// Parse and validate flags.
func init() {
	
	// Unset localhost defaults if certificate file can not be found.
	certFileExists, e := cfgutil.FileExists(opts.RPCCertificateFile)
	if e != nil {
		fatalf("%v", err)
	}
	if !certFileExists {
		opts.RPCConnect = ""
		opts.RPCCertificateFile = ""
	}
	_, e = flags.Parse(&opts)
	if e != nil {
		os.Exit(1)
	}
	if opts.TestNet3 && opts.SimNet {
		fatalf("Multiple bitcoin networks may not be used simultaneously")
	}
	var activeNet = &netparams.MainNetParams
	if opts.TestNet3 {
		activeNet = &netparams.TestNet3Params
	} else if opts.SimNet {
		activeNet = &netparams.SimNetParams
	}
	if opts.RPCConnect == "" {
		fatalf("RPC hostname[:port] is required")
	}
	rpcConnect, e := cfgutil.NormalizeAddress(opts.RPCConnect, activeNet.WalletRPCServerPort)
	if e != nil {
		fatalf("Invalid RPC network address `%v`: %v", opts.RPCConnect, err)
	}
	opts.RPCConnect = rpcConnect
	if opts.RPCUsername == "" {
		fatalf("RPC username is required")
	}
	certFileExists, e = cfgutil.FileExists(opts.RPCCertificateFile)
	if e != nil {
		fatalf("%v", err)
	}
	if !certFileExists {
		fatalf("RPC certificate file `%s` not found", opts.RPCCertificateFile)
	}
	if opts.FeeRate.Amount > 1e6 {
		fatalf("Fee rate `%v/kB` is exceptionally high", opts.FeeRate.Amount)
	}
	if opts.FeeRate.Amount < 1e2 {
		fatalf("Fee rate `%v/kB` is exceptionally low", opts.FeeRate.Amount)
	}
	if opts.SourceAccount == opts.DestinationAccount {
		fatalf("Source and destination accounts should not be equal")
	}
	if opts.RequiredConfirmations < 0 {
		fatalf("Required confirmations must be non-negative")
	}
}

// noInputValue describes an error returned by the input source when no inputs were selected because each previous
// output value was zero. Callers of txauthor.NewUnsignedTransaction need not report these errors to the user.
type noInputValue struct {
}

func (noInputValue) Error() string { return "no input value" }

// makeInputSource creates an InputSource that creates inputs for every unspent output with non-zero output values. The
// target amount is ignored since every output is consumed. The InputSource does not return any previous output scripts
// as they are not needed for creating the unsinged transaction and are looked up again by the wallet during the call to
// signrawtransaction.
func makeInputSource(
	outputs []btcjson.ListUnspentResult,
) txauthor.InputSource {
	var (
		totalInputValue util.Amount
		inputs          = make([]*wire.TxIn, 0, len(outputs))
		inputValues     = make([]util.Amount, 0, len(outputs))
		sourceErr       error
	)
	for _, output := range outputs {
		outputAmount, e := util.NewAmount(output.Amount)
		if e != nil {
			sourceErr = fmt.Errorf(
				"invalid amount `%v` in listunspent result",
				output.Amount,
			)
			break
		}
		if outputAmount == 0 {
			continue
		}
		if !saneOutputValue(outputAmount) {
			sourceErr = fmt.Errorf(
				"impossible output amount `%v` in listunspent result",
				outputAmount,
			)
			break
		}
		totalInputValue += outputAmount
		previousOutPoint, e := parseOutPoint(&output)
		if e != nil {
			sourceErr = fmt.Errorf(
				"invalid data in listunspent result: %v",
				e,
			)
			break
		}
		inputs = append(inputs, wire.NewTxIn(&previousOutPoint, nil, nil))
		inputValues = append(inputValues, outputAmount)
	}
	if sourceErr == nil && totalInputValue == 0 {
		sourceErr = noInputValue{}
	}
	return func(util.Amount) (util.Amount, []*wire.TxIn, []util.Amount, [][]byte, error) {
		return totalInputValue, inputs, inputValues, nil, sourceErr
	}
}

// makeDestinationScriptSource creates a ChangeSource which is used to receive all correlated previous input value. A
// non-change address is created by this function.
func makeDestinationScriptSource(
	rpcClient *rpcclient.Client, accountName string,
) txauthor.ChangeSource {
	return func() ([]byte, error) {
		destinationAddress, e := rpcClient.GetNewAddress(accountName)
		if e != nil {
			return nil, e
		}
		return txscript.PayToAddrScript(destinationAddress)
	}
}
func main() {
	e := sweep()
	if e != nil {
		fatalf("%v", err)
	}
}
func sweep() (e error) {
	rpcPassword, e := promptSecret("Wallet RPC password")
	if e != nil {
		return errContext(e, "failed to read RPC password")
	}
	// Open RPC client.
	rpcCertificate, e := ioutil.ReadFile(opts.RPCCertificateFile)
	if e != nil {
		return errContext(e, "failed to read RPC certificate")
	}
	rpcClient, e := rpcclient.New(
		&rpcclient.ConnConfig{
			Host:         opts.RPCConnect,
			User:         opts.RPCUsername,
			Pass:         rpcPassword,
			Certificates: rpcCertificate,
			HTTPPostMode: true,
		}, nil, qu.T(),
	)
	if e != nil {
		return errContext(e, "failed to create RPC client")
	}
	defer rpcClient.Shutdown()
	// Fetch all unspent outputs, ignore those not from the source account, and group by their destination address. Each
	// grouping of outputs will be used as inputs for a single transaction sending to a new destination account address.
	unspentOutputs, e := rpcClient.ListUnspent()
	if e != nil {
		return errContext(e, "failed to fetch unspent outputs")
	}
	sourceOutputs := make(map[string][]btcjson.ListUnspentResult)
	for _, unspentOutput := range unspentOutputs {
		if !unspentOutput.Spendable {
			continue
		}
		if unspentOutput.Confirmations < opts.RequiredConfirmations {
			continue
		}
		if unspentOutput.Account != opts.SourceAccount {
			continue
		}
		sourceAddressOutputs := sourceOutputs[unspentOutput.Address]
		sourceOutputs[unspentOutput.Address] = append(sourceAddressOutputs, unspentOutput)
	}
	var privatePassphrase string
	if len(sourceOutputs) != 0 {
		privatePassphrase, e = promptSecret("Wallet private passphrase")
		if e != nil {
			return errContext(e, "failed to read private passphrase")
		}
	}
	var totalSwept util.Amount
	var numErrors int
	var reportError = func(format string, args ...interface{}) {
		if _, e = fmt.Fprintf(os.Stderr, format, args...); E.Chk(e) {
		}
		if _, e = os.Stderr.Write(newlineBytes); E.Chk(e) {
		}
		numErrors++
	}
	for _, previousOutputs := range sourceOutputs {
		inputSource := makeInputSource(previousOutputs)
		destinationSource := makeDestinationScriptSource(rpcClient, opts.DestinationAccount)
		tx, e := txauthor.NewUnsignedTransaction(
			nil, opts.FeeRate.Amount,
			inputSource, destinationSource,
		)
		if e != nil {
			if e != (noInputValue{}) {
				reportError("Failed to create unsigned transaction: %v", err)
			}
			continue
		}
		// Unlock the wallet, sign the transaction, and immediately lock.
		e = rpcClient.WalletPassphrase(privatePassphrase, 60)
		if e != nil {
			reportError("Failed to unlock wallet: %v", err)
			continue
		}
		signedTransaction, complete, e := rpcClient.SignRawTransaction(tx.Tx)
		_ = rpcClient.WalletLock()
		if e != nil {
			reportError("Failed to sign transaction: %v", err)
			continue
		}
		if !complete {
			reportError("Failed to sign every input")
			continue
		}
		// Publish the signed sweep transaction.
		txHash, e := rpcClient.SendRawTransaction(signedTransaction, false)
		if e != nil {
			reportError("Failed to publish transaction: %v", err)
			continue
		}
		outputAmount := util.Amount(tx.Tx.TxOut[0].Value)
		I.F(
			"Swept %v to destination account with transaction %v\n",
			outputAmount, txHash,
		)
		totalSwept += outputAmount
	}
	numPublished := len(sourceOutputs) - numErrors
	transactionNoun := pickNoun(numErrors, "transaction", "transactions")
	if numPublished != 0 {
		I.F(
			"Swept %v to destination account across %d %s\n",
			totalSwept, numPublished, transactionNoun,
		)
	}
	if numErrors > 0 {
		return fmt.Errorf("failed to publish %d %s", numErrors, transactionNoun)
	}
	return nil
}
func promptSecret(what string) (string, error) {
	fmt.Printf("%s: ", what)
	fd := int(os.Stdin.Fd())
	input, e := terminal.ReadPassword(fd)
	fmt.Println()
	if e != nil {
		return "", e
	}
	return string(input), nil
}

func saneOutputValue(
	amount util.Amount,
) bool {
	return amount >= 0 && amount <= util.MaxSatoshi
}

func parseOutPoint(
	input *btcjson.ListUnspentResult,
) (wire.OutPoint, error) {
	txHash, e := chainhash.NewHashFromStr(input.TxID)
	if e != nil {
		return wire.OutPoint{}, e
	}
	return wire.OutPoint{Hash: *txHash, Index: input.Vout}, nil
}

func pickNoun(
	n int, singularForm, pluralForm string,
) string {
	if n == 1 {
		return singularForm
	}
	return pluralForm
}

var subsystem = logg.AddLoggerSubsystem()
var ftl, err, wrn, inf, dbg, trc logg.LevelPrinter = logg.GetLogPrinterSet(subsystem)

func init() {
	// var _ = logg.AddFilteredSubsystem(subsystem)
	// var _ = logg.AddHighlightedSubsystem(subsystem)
	F.Ln("F.Ln")
	E.Ln("E.Ln")
	W.Ln("W.Ln")
	I.Ln("I.Ln")
	D.Ln("D.Ln")
	F.Ln("T.Ln")
	F.F("%s", "F.F")
	E.F("%s", "E.F")
	W.F("%s", "W.F")
	I.F("%s", "I.F")
	D.F("%s", "D.F")
	T.F("%s", "T.F")
	ftl.C(func() string { return "ftl.C" })
	err.C(func() string { return "err.C" })
	W.C(func() string { return "W.C" })
	I.C(func() string { return "inf.C" })
	D.C(func() string { return "D.C" })
	T.C(func() string { return "T.C" })
	ftl.C(func() string { return "ftl.C" })
	E.Chk(errors.New("E.Chk"))
	W.Chk(errors.New("W.Chk"))
	I.Chk(errors.New("inf.Chk"))
	D.Chk(errors.New("D.Chk"))
	T.Chk(errors.New("T.Chk"))
}
