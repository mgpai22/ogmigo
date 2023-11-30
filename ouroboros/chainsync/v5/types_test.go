// Copyright 2021 Matt Ho
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v5

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/stretchr/testify/assert"
)

func TestV5(t *testing.T) {
	t.Run("TxFromV6", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/Tx_v6.json")
		assert.Nil(t, err)

		var expected chainsync.Tx
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)
		actual := TxFromV6(expected)

		// Signature: feb447eca6b819f88b4b6aac81a97200121dd4451bfccb65f73d31de296b5402eead42848808633ad1d4afbfad3a2f1967ce65516adca5cd373673f758a9c096
		bootstrap := "{\"key\":\"d88f6028cc3d6d335115de3737bc2fe80a9a57a21a2c7c228ebc33b222e0897b\",\"signature\":\"/rRH7Ka4GfiLS2qsgalyABId1EUb/Mtl9z0x3ilrVALurUKEiAhjOtHUr7+tOi8ZZ85lUWrcpc03NnP3WKnAlg==\",\"chainCode\":\"12340000\",\"addressAttributes\":\"Lw==\"}"
		assert.Equal(t, "feb447eca6b819f88b4b6aac81a97200121dd4451bfccb65f73d31de296b5402eead42848808633ad1d4afbfad3a2f1967ce65516adca5cd373673f758a9c096", expected.Signatories[0].Signature)
		assert.Equal(t, bootstrap, string(actual.Witness.Bootstrap[0]))
		assert.Equal(t, "IFb1lTq+ivhYQz6fAoPZQXuGgebeh5fIsM8rocK03mbss8yaUQpf871Qso2aAYaxjDadDHzMfUPRCJDpTyVxQg==", actual.Witness.Signatures["400019217786c3630fb121c455065b879055aa0ced5076a24abe8d6c837e0318"])
	})
}
