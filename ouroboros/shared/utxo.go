package shared

import (
	"encoding/json"
)

type UtxoTxID struct {
	ID string `json:"id"`
}

type Utxo struct {
	// TxOut "ref" fields.
	Transaction UtxoTxID `json:"transaction"`
	Index       uint32   `json:"index"`

	// TxOut fields.
	Address   string          `json:"address"`
	Value     Value           `json:"value"`
	DatumHash string          `json:"datumHash,omitempty"`
	Datum     string          `json:"datum,omitempty"`
	Script    json.RawMessage `json:"script,omitempty"`
}
