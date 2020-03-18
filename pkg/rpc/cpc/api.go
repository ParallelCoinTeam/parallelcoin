package api

import (
	"github.com/p9c/pod/cmd/node/rpc"
	"github.com/p9c/pod/pkg/rpc/btcjson"
)

// StartAPI starts up the api handler server that receives rpc.API messages and runs the handler and returns the result
// Note that the parameters are type asserted to prevent the consumer of the API from sending wrong message types not
// because it's necessary since they are interfaces end to end
func StartAPI(server *rpc.Server, quit chan struct{}) {
	nrh := rpc.RPCHandlers
	go func() {
		var err error
		var res interface{}
		for {
			select {
			case msg := <-nrh["addnode"].Call:
				if _, err = nrh["addnode"].
					Fn(server, msg.Params.(btcjson.AddNodeCmd),
						nil); L.Check(err) {
				}
				msg.Ch.(chan rpc.AddNodeRes) <- rpc.AddNodeRes{
					Res: nil, Err: err}
			case msg := <-nrh["createrawtransaction"].Call:
				if res, err = nrh["createrawtransaction"].
					Fn(server, msg.Params.(btcjson.CreateRawTransactionCmd),
						nil); L.Check(err) {
				}
				msg.Ch.(chan rpc.CreateRawTransactionRes) <- rpc.CreateRawTransactionRes{
					Res: res.(string), Err: err}
			case msg := <-nrh["decoderawtransaction"].Call:
				var ret btcjson.TxRawDecodeResult
				if res, err = nrh["decoderawtransaction"].
					Fn(server, msg.Params.(btcjson.DecodeRawTransactionCmd),
						nil); L.Check(err) {
				} else {
					ret = res.(btcjson.TxRawDecodeResult)
				}
				msg.Ch.(chan rpc.DecodeRawTransactionRes) <- rpc.DecodeRawTransactionRes{
					Res: ret, Err: err}
			case msg := <-nrh["decodescript"].Call:
				if res, err = nrh["decodescript"].
					Fn(server, msg.Params.(btcjson.DecodeScriptCmd),
						nil); L.Check(err) {
				}
				msg.Ch.(chan rpc.DecodeScriptRes) <- rpc.DecodeScriptRes{
					Res: res.(btcjson.DecodeScriptResult), Err: err}
			case msg := <-nrh["estimatefee"].Call:
				if res, err = nrh["estimatefee"].
					Fn(server, msg.Params.(btcjson.EstimateFeeCmd),
						nil); L.Check(err) {
				}
				msg.Ch.(chan rpc.EstimateFeeRes) <- rpc.EstimateFeeRes{
					Res: res.(float64), Err: err}
			case <-quit:
				return
			}
		}
	}()
}
