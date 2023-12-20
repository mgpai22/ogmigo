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

		var expectedV6 chainsync.Tx
		err = json.Unmarshal(rawData, &expectedV6)
		assert.Nil(t, err)
		v5Conversion := TxFromV6(expectedV6)
		v6Conversion := v5Conversion.ConvertToV6()

		var bootstrapV5Sigs []chainsync.Signature
		for _, x := range v5Conversion.Witness.Bootstrap {
			var sig chainsync.Signature
			json.Unmarshal(x, &sig)
			bootstrapV5Sigs = append(bootstrapV5Sigs, sig)
		}

		network := "mainnet"
		bootstrap := "{\"key\":\"d88f6028cc3d6d335115de3737bc2fe80a9a57a21a2c7c228ebc33b222e0897b\",\"signature\":\"/rRH7Ka4GfiLS2qsgalyABId1EUb/Mtl9z0x3ilrVALurUKEiAhjOtHUr7+tOi8ZZ85lUWrcpc03NnP3WKnAlg==\",\"chainCode\":\"12340000\",\"addressAttributes\":\"Lw==\"}"
		assert.Equal(t, "2f", expectedV6.Signatories[0].AddressAttributes)
		assert.Equal(t, "12340000", expectedV6.Signatories[0].ChainCode)
		assert.Equal(t, "d88f6028cc3d6d335115de3737bc2fe80a9a57a21a2c7c228ebc33b222e0897b", expectedV6.Signatories[0].Key)
		assert.Equal(t, "feb447eca6b819f88b4b6aac81a97200121dd4451bfccb65f73d31de296b5402eead42848808633ad1d4afbfad3a2f1967ce65516adca5cd373673f758a9c096", expectedV6.Signatories[0].Signature)
		assert.Equal(t, bootstrap, string(v5Conversion.Witness.Bootstrap[0]))
		assert.Equal(t, "IFb1lTq+ivhYQz6fAoPZQXuGgebeh5fIsM8rocK03mbss8yaUQpf871Qso2aAYaxjDadDHzMfUPRCJDpTyVxQg==", v5Conversion.Witness.Signatures["400019217786c3630fb121c455065b879055aa0ced5076a24abe8d6c837e0318"])
		assert.Equal(t, json.RawMessage(network), v5Conversion.Body.Network)
		assert.Equal(t, network, v6Conversion.Network)
		assert.Equal(t, expectedV6.Signatories[0].AddressAttributes, v6Conversion.Signatories[2].AddressAttributes)
		assert.Equal(t, expectedV6.Signatories[0].ChainCode, v6Conversion.Signatories[2].ChainCode)
		assert.Equal(t, expectedV6.Signatories[0].Key, v6Conversion.Signatories[2].Key)
		assert.Equal(t, expectedV6.Signatories[0].Signature, v6Conversion.Signatories[2].Signature)
	})
}
