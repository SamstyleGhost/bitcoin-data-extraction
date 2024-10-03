package custom

import "time"

type Transactions struct {
	Address    string  `json:"address,omitempty"`
	Amount     float64 `json:"amount,omitempty"`
	IsStandard bool    `json:"is_standard,omitempty"`
	NextTX     string  `json:"next_tx,omitempty"`
}

type TransactionRow struct {
	Date           string  `json:"date"`
	ReceivedFrom   string  `json:"receivedFrom"`
	ReceivedAmount float64 `json:"receivedAmount"`
	SentAmount     float64 `json:"sentAmount"`
	SentTo         string  `json:"sentTo"`
	Balance        float64 `json:"balance"`
	Transaction    string  `json:"transaction"`
}

type CashFlowTransaction struct {
	Found       bool           `json:"found,omitempty"`
	Label       string         `json:"label,omitempty"` // Will omit the Label value if it is not found
	TxID        string         `json:"txid,omitempty"`
	IsCoinbase  bool           `json:"is_coinbase,omitempty"`
	WalletID    string         `json:"wallet_id,omitempty"`
	BlockHeight int32          `json:"block_height,omitempty"`
	BlockPos    int32          `json:"block_pos,omitempty"`
	Time        time.Time      `json:"time,omitempty"`
	Size        int16          `json:"size,omitempty"`
	Ins         []Transactions `json:"ins,omitempty"`
	Outs        []Transactions `json:"outs,omitempty"`
}
