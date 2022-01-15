package tx

import (
	"sanus/sanus-sdk/cc/transfer"
	"sanus/sanus-sdk/cc/utils"
)

type CCData struct {
	Type              string
	Amount            int                     `json:"amount"`
	Protocol          int64                   `json:"protocol"`
	Version           int64                   `json:"version"`
	Divisibility      int                     `json:"divisibility"`
	LockStatus        bool                    `json:"lock_status"`
	AggregationPolicy string                  `json:"aggregation_policy"`
	MultiSig          []transfer.MultiSigData `json:"multi_sig"`
	NoRules           bool                    `json:"no_rules"`
	Sha2              []byte                  `json:"sha_2"`
	TorrentHash       []byte                  `json:"torrent_hash"`
	Payments          []*utils.PaymentData    `json:"payments"`
}
