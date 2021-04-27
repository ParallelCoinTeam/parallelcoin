package txscript_test

import (
	"encoding/hex"
	"fmt"
	
	"github.com/p9c/pod/pkg/chain/config/netparams"
	chainhash "github.com/p9c/matrjoska/pkg/chainhash"
	txscript "github.com/p9c/matrjoska/pkg/txscript"
	"github.com/p9c/matrjoska/pkg/wire"
	ecc "github.com/p9c/matrjoska/pkg/ecc"
	"github.com/p9c/pod/pkg/util"
)

// This example demonstrates creating a script which pays to a bitcoin address. It also prints the created script hex
// and uses the DisasmString function to display the disassembled script.
func ExamplePayToAddrScript() {
	// Parse the address to send the coins to into a util.Address which is useful to ensure the accuracy of the address and determine the address type.  It is also required for the upcoming call to
	// PayToAddrScript.
	addressStr := "12gpXQVcCL2qhTNQgyLVdCFG2Qs2px98nV"
	address, e := util.DecodeAddress(addressStr, &netparams.MainNetParams)
	if e != nil {
		return
	}
	// Create a public key script that pays to the address.
	script, e := txscript.PayToAddrScript(address)
	if e != nil {
		return
	}
	fmt.Printf("Script Hex: %x\n", script)
	disasm, e := txscript.DisasmString(script)
	if e != nil {
		return
	}
	fmt.Println("Script Disassembly:", disasm)
	// Output:
	// Script Hex: 76a914128004ff2fcaf13b2b91eb654b1dc2b674f7ec6188ac
	// Script Disassembly: OP_DUP OP_HASH160 128004ff2fcaf13b2b91eb654b1dc2b674f7ec61 OP_EQUALVERIFY OP_CHECKSIG
}

// // This example demonstrates extracting information from a standard public key script.
// func ExampleExtractPkScriptAddrs() {
// 	// Start with a standard pay-to-pubkey-hash script.
// 	scriptHex := "76a914128004ff2fcaf13b2b91eb654b1dc2b674f7ec6188ac"
// 	script, e := hex.DecodeString(scriptHex)
// 	if e != nil  {
// 		L.Script// 		return
// 	}
// 	// Extract and print details from the script.
// 	scriptClass, addresses, reqSigs, e := txscript.ExtractPkScriptAddrs(
// 		script, &netparams.MainNetParams)
// 	if e != nil  {
// 		L.Script// 		return
// 	}
// 	fmt.Println("Script Class:", scriptClass)
// 	fmt.Println("Addresses:", addresses)
// 	fmt.Println("Required Signatures:", reqSigs)
// 	// Output:
// 	// Script Class: pubkeyhash
// 	// Addresses: [12gpXQVcCL2qhTNQgyLVdCFG2Qs2px98nV]
// 	// Required Signatures: 1
// }

// This example demonstrates manually creating and signing a redeem transaction.
func ExampleSignTxOutput() {
	// Ordinarily the private key would come from whatever storage mechanism is being used, but for this example just hard code it.
	privKeyBytes, e := hex.DecodeString("22a47fa09a223f2aa079edf85a7c2" +
		"d4f8720ee63e502ee2869afab7de234b80c",
	)
	if e != nil {
		return
	}
	privKey, pubKey := ecc.PrivKeyFromBytes(ecc.S256(), privKeyBytes)
	pubKeyHash := util.Hash160(pubKey.SerializeCompressed())
	addr, e := util.NewAddressPubKeyHash(pubKeyHash,
		&netparams.MainNetParams,
	)
	if e != nil {
		return
	}
	// For this example, create a fake transaction that represents what would ordinarily be the real transaction that is being spent.  It contains a single output that pays to address in the amount of 1 DUO.
	originTx := wire.NewMsgTx(wire.TxVersion)
	prevOut := wire.NewOutPoint(&chainhash.Hash{}, ^uint32(0))
	txIn := wire.NewTxIn(prevOut, []byte{txscript.OP_0, txscript.OP_0}, nil)
	originTx.AddTxIn(txIn)
	pkScript, e := txscript.PayToAddrScript(addr)
	if e != nil {
		return
	}
	txOut := wire.NewTxOut(100000000, pkScript)
	originTx.AddTxOut(txOut)
	originTxHash := originTx.TxHash()
	// Create the transaction to redeem the fake transaction.
	redeemTx := wire.NewMsgTx(wire.TxVersion)
	// Add the input(s) the redeeming transaction will spend.  There is no signature script at this point since it hasn't been created or signed yet, hence nil is provided for it.
	prevOut = wire.NewOutPoint(&originTxHash, 0)
	txIn = wire.NewTxIn(prevOut, nil, nil)
	redeemTx.AddTxIn(txIn)
	// Ordinarily this would contain that actual destination of the funds, but for this example don't bother.
	txOut = wire.NewTxOut(0, nil)
	redeemTx.AddTxOut(txOut)
	// Sign the redeeming transaction.
	lookupKey := func(a util.Address) (*ecc.PrivateKey, bool, error) {
		// Ordinarily this function would involve looking up the private key for the provided address, but since the only thing being signed in this example uses the address associated with the private key from above, simply return it with the compressed flag set since the address is using the associated compressed public key.
		// NOTE: If you want to prove the code is actually signing the transaction properly, uncomment the following line which intentionally returns an invalid key to sign with, which in turn will result in a failure during the script execution when verifying the signature.
		// privKey.D.SetInt64(12345)
		//
		return privKey, true, nil
	}
	// Notice that the script database parameter is nil here since it isn't used.  It must be specified when pay-to-script-hash transactions are being signed.
	sigScript, e := txscript.SignTxOutput(&netparams.MainNetParams,
		redeemTx, 0, originTx.TxOut[0].PkScript, txscript.SigHashAll,
		txscript.KeyClosure(lookupKey), nil, nil,
	)
	if e != nil {
		return
	}
	redeemTx.TxIn[0].SignatureScript = sigScript
	// Prove that the transaction has been validly signed by executing the script pair.
	flags := txscript.ScriptBip16 | txscript.ScriptVerifyDERSignatures |
		txscript.ScriptStrictMultiSig |
		txscript.ScriptDiscourageUpgradableNops
	vm, e := txscript.NewEngine(originTx.TxOut[0].PkScript, redeemTx, 0,
		flags, nil, nil, -1,
	)
	if e != nil {
		return
	}
	if e := vm.Execute(); dbg.Chk(e) {
		return
	}
	fmt.Println("Transaction successfully signed")
	// Output:
	// Transaction successfully signed
}
