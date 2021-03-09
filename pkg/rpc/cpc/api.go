package api

import (
	"github.com/p9c/pod/pkg/rpc/btcjson"
	"github.com/p9c/pod/pkg/rpc/chainrpc"
	qu "github.com/p9c/pod/pkg/util/qu"
)

// StartAPI starts up the api handler server that receives rpc.API messages and runs the handler and returns the result
// Note that the parameters are type asserted to prevent the consumer of the API from sending wrong message types not
// because it's necessary since they are interfaces end to end
func StartAPI(server *chainrpc.Server, quit qu.C) {
	nrh := chainrpc.RPCHandlers
	go func() {
		var e error
		var res interface{}
		for {
			select {
			case msg := <-nrh["addnode"].Call:
				if _, e = nrh["addnode"].
					Fn(
						server, msg.Params.(btcjson.AddNodeCmd),
						nil,
				); err.Chk(e) {
				}
				msg.Ch.(chan chainrpc.AddNodeRes) <- chainrpc.AddNodeRes{
					Res: nil, Err: err,
				}
			case msg := <-nrh["createrawtransaction"].Call:
				if res, e = nrh["createrawtransaction"].
					Fn(server, msg.Params.(btcjson.CreateRawTransactionCmd), nil); err.Chk(e) {
				}
				msg.Ch.(chan chainrpc.CreateRawTransactionRes) <- chainrpc.CreateRawTransactionRes{
					Res: res.(*string), Err: err,
				}
			case msg := <-nrh["decoderawtransaction"].Call:
				var ret btcjson.TxRawDecodeResult
				if res, e = nrh["decoderawtransaction"].Fn(
					server, msg.Params.(btcjson.DecodeRawTransactionCmd),
					nil,
				); err.Chk(e) {
				} else {
					ret = res.(btcjson.TxRawDecodeResult)
				}
				msg.Ch.(chan chainrpc.DecodeRawTransactionRes) <- chainrpc.DecodeRawTransactionRes{
					Res: &ret, Err: err,
				}
			case msg := <-nrh["decodescript"].Call:
				if res, e = nrh["decodescript"].Fn(server, msg.Params.(btcjson.DecodeScriptCmd), nil); err.Chk(e) {
				}
				msg.Ch.(chan chainrpc.DecodeScriptRes) <- chainrpc.DecodeScriptRes{
					Res: res.(*btcjson.DecodeScriptResult), Err: err,
				}
			case msg := <-nrh["estimatefee"].Call:
				if res, e = nrh["estimatefee"].
					Fn(
						server, msg.Params.(btcjson.EstimateFeeCmd),
						nil,
				); err.Chk(e) {
				}
				msg.Ch.(chan chainrpc.EstimateFeeRes) <- chainrpc.EstimateFeeRes{
					Res: res.(*float64), Err: err,
				}
			case <-quit.Wait():
				return
			}
		}
	}()
}
