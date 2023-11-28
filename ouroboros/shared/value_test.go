package shared

import (
	"fmt"
	"testing"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
	"github.com/tj/assert"
)

func Test_AssetsExceptADA(t *testing.T) {
	v := Value{
		"ada": {
			"lovelace": num.Uint64(1),
		},
		"policy1": {
			"asset1": num.Uint64(1),
			"asset2": num.Uint64(2),
		},
		"policy2": {
			"asset3": num.Uint64(3),
			"asset4": num.Uint64(4),
		},
	}

	a, _ := v.AssetsExceptAda()
	fmt.Printf("%v\n", a)
	assert.EqualValues(t, a, Value{
		"policy1": {
			"asset1": num.Uint64(1),
			"asset2": num.Uint64(2),
		},
		"policy2": {
			"asset3": num.Uint64(3),
			"asset4": num.Uint64(4),
		},
	})
}

func Test_Enough(t *testing.T) {
	have := Value{
		"ada": {
			"lovelace": num.Uint64(437041203),
		},
		"da8c30857834c6ae7203935b89278c532b3995245295456f993e1d24": {
			"4c51": num.Uint64(14310359231),
		},
	}
	want := Value{
		"da8c30857834c6ae7203935b89278c532b3995245295456f993e1d24": {
			"4c51": num.Uint64(1023291),
		},
		"25c5de5f5b286073c593edfd77b48abc7a48e5a4f3d4cd9d428ff935": {
			"55534454": num.Uint64(3449),
		},
	}
	ok, _ := Enough(have, want)
	assert.False(t, ok)
}
