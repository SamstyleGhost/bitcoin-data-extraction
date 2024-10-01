package custom

import "time"

type Transactions struct {
	Address    string  `json:"address"`
	Amount     float64 `json:"amount"`
	IsStandard bool    `json:"is_standard"`
	NextTX     string  `json:"next_tx"`
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
	Found       bool           `json:"found"`
	Label       string         `json:"label,omitempty"` // Will omit the Label value if it is not found
	TxID        string         `json:"txid"`
	IsCoinbase  bool           `json:"is_coinbase"`
	WalletID    string         `json:"wallet_id"`
	BlockHeight int32          `json:"block_height"`
	BlockPos    int32          `json:"block_pos"`
	Time        time.Time      `json:"time"`
	Size        int16          `json:"size"`
	Ins         []Transactions `json:"ins"`
	Outs        []Transactions `json:"outs"`
}
