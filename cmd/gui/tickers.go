package gui

import (
	"github.com/p9c/pod/pkg/log"
	"time"
)

func (r *rcvar) DuOSgatherer() {
	go func() {
		for _ = range time.NewTicker(time.Second * 1).C {
			var err error
			// Status
			//r.GetDuOStatus()
			//err = r.Render("osStatus", r.status)
			//if err != nil {
			//	log.ERROR("error binding to webview:", err)
			//}
			// Hashes
			r.GetDuOShashesPerSec()
			err = r.Render("osHashes", r.hashes)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// NetHash
			r.GetDuOSnetworkHashesPerSec()
			err = r.Render("osNetHash", r.nethash)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// Height
			r.GetDuOSheight()
			err = r.Render("osHeight", r.height)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// BestBlock
			r.GetDuOSbestBlockHash()
			err = r.Render("osBestBlock", r.bestblock)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// Difficulty

			err = r.Render("osDifficulty", r.difficulty)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// BlockCount
			r.GetDuOSblockCount()
			err = r.Render("osBlockCount", r.blockcount)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// NetlastBlock
			r.GetDuOSnetworkLastBlock()
			err = r.Render("osNetlastBlock", r.netlastblock)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// Connections
			r.GetDuOSconnectionCount()
			err = r.Render("osConnections", r.connections)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// Balance
			r.GetDuOSbalance()
			err = r.Render("osBalance", r.balance)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// Unconfirmed
			r.GetDuOSunconfirmedBalance()
			err = r.Render("osUnconfirmed", r.unconfirmed)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// TxsNumber
			err = r.Render("osTxsNumber", r.txsnumber)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// Transactions
			//r.GetDuOStransactions(0, 10, "all")
			//err = r.Render("osTransactions", r.transactions)
			//if err != nil {
			//	log.ERROR("error binding to webview:", err)
			//}
			// Txs

			err = r.Render("osTxs", r.txs)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
			// LastTxs
			r.GetDuOSlastTxs()
			err = r.Render("osLastTxs", r.lasttxs)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}


			// LocalHost
			r.GetDuOSlocalLost()
			err = r.Render("osLocalHost", r.localhost)
			if err != nil {
				log.ERROR("error binding to webview:", err)
			}
//log.INFO("Compler rcvar -->>>>>", r)

		}
	}()

}
