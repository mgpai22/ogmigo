package statequery

import (
	"math/big"
)

type EraStart struct {
	Time  EraSeconds `json:"time,omitempty"`
	Slot  big.Int    `json:"slot,omitempty"`
	Epoch big.Int    `json:"epoch,omitempty"`
}

type EraSeconds struct {
	Seconds big.Int `json:"seconds"`
}

type EraMilliseconds struct {
	Milliseconds big.Int `json:"milliseconds"`
}
