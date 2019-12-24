package models

import (
	"time"
)

type
	DuoUIbalance struct {
		Balance string `json:"balance"`
	}
type
	DuoUIunconfirmed struct {
		Unconfirmed string `json:"unconfirmed"`
	}
type
	DuoUItransactions struct {
		Txs       []DuoUItx `json:"txs"`
		TxsNumber int       `json:"txsnumber"`
	}
type
	DuoUItx struct {
		TxID          string
		Amount        float64
		Category      string
		Confirmations int64
		Time          time.Time
		Added         time.Time
	}

type
	DuoUItransactionsNumber struct {
		TxsNumber int `json:"txsnumber"`
	}
type
	DuoUItransactionsExcerpts struct {
		Txs           []DuoUItransactionExcerpt `json:"txs"`
		TxsNumber     int                       `json:"txsnumber"`
		Balance       float64                   `json:"balance"`
		BalanceHeight float64                   `json:"balanceheight"`
	}
type
	DuoUItransactionExcerpt struct {
		Balance       float64 `json:"balance"`
		Amount        float64 `json:"amount"`
		Category      string  `json:"category"`
		Confirmations int64   `json:"confirmations"`
		Time          string  `json:"time"`
		TxID          string  `json:"txid"`
		Comment       string  `json:"comment,omitempty"`
	}
