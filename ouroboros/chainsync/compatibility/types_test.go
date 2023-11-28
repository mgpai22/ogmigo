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

package compatibility

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/num"
	v5 "github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/chainsync/v5"
	"github.com/SundaeSwap-finance/ogmigo/v6/ouroboros/shared"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/tj/assert"
)

// TODO:
//   - break this up into smaller tests;
//   - put strings in a test_data folder;
//   - check the data we get to make sure it round tripped correctly
//   - test dynamodb marshalling
func TestCompatibleResult(t *testing.T) {
	t.Run("Roll Forward V5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/RollForward_v5.json")
		assert.Nil(t, err)

		var expected v5.ResultNextBlockV5
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultNextBlock
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected.ConvertToV6(), chainsync.ResultNextBlockPraos(compatible))

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultNextBlockV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		// Ideally we'd like to check the full type round tripped, but for various reasons
		// this is super awkward, so we just check a few important properties
		assert.NotNil(t, got.RollForward)
		assert.NotNil(t, got.RollForward.Block.Babbage)
		assert.Equal(t, got.RollForward.Block.Babbage.HeaderHash, got.RollForward.Block.Babbage.HeaderHash)
	})

	t.Run("Roll Backward V5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/RollBackward_v5.json")
		assert.Nil(t, err)

		var expected v5.ResultNextBlockV5
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultNextBlock
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected.ConvertToV6(), compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultNextBlockV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.NotNil(t, got.RollBackward)
		assert.EqualValues(t, got.RollBackward.Point.String(), expected.RollBackward.Point.String())
	})

	t.Run("Roll Forward V6", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/RollForward_v6.json")
		assert.Nil(t, err)

		var expected chainsync.ResultNextBlockPraos
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultNextBlock
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultNextBlockV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.NotNil(t, got.RollForward)
		assert.NotNil(t, got.RollForward.Block.Allegra)
		assert.Equal(t, got.RollForward.Block.Allegra.HeaderHash, got.RollForward.Block.Allegra.HeaderHash)
	})

	t.Run("Roll Backward V6", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/RollBackward_v6.json")
		assert.Nil(t, err)

		var expected chainsync.ResultNextBlockPraos
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultNextBlock
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultNextBlockV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.NotNil(t, got.RollBackward)
		assert.EqualValues(t, got.RollBackward.Point.ConvertToV6(), *expected.Point)
	})

	t.Run("Intersection Found V5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/IntersectionFound_v5.json")
		assert.Nil(t, err)

		var expected v5.ResultFindIntersectionV5
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultFindIntersection
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected.ConvertToV6(), compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultFindIntersectionV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.NotNil(t, got.IntersectionFound)
		assert.NotNil(t, got.IntersectionFound.Tip)
		assert.Equal(t, got.IntersectionFound.Tip.Hash, expected.IntersectionFound.Tip.Hash)
	})

	t.Run("Intersection Not Found V5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/IntersectionNotFound_v5.json")
		assert.Nil(t, err)

		var expected v5.ResultFindIntersectionV5
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultFindIntersection
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected.ConvertToV6(), compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultFindIntersectionV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.NotNil(t, got.IntersectionNotFound)
		assert.NotNil(t, got.IntersectionNotFound.Tip)
		assert.Equal(t, got.IntersectionNotFound.Tip.Hash, expected.IntersectionNotFound.Tip.Hash)
	})

	t.Run("Intersection Found V6", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/IntersectionFound_v6.json")
		assert.Nil(t, err)

		var expected chainsync.ResultFindIntersectionPraos
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultFindIntersection
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultFindIntersectionV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.NotNil(t, got.IntersectionFound)
		assert.NotNil(t, got.IntersectionFound.Tip)
		assert.Equal(t, got.IntersectionFound.Tip.Hash, expected.Tip.ID)
	})

	t.Run("Intersection Not Found V6", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/IntersectionNotFound_v6.json")
		assert.Nil(t, err)

		var expected chainsync.ResultFindIntersectionPraos
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResultFindIntersection
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResultFindIntersectionV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.NotNil(t, got.IntersectionNotFound)
	})

	t.Run("Real World Intersection Found", func(t *testing.T) {
		dataV5Result, err := os.ReadFile("test_data/RealWorld_IntersectionFound_v5.json")
		assert.Nil(t, err)
		var method9 CompatibleResult
		err = json.Unmarshal([]byte(dataV5Result), &method9)
		assert.Nil(t, err)
	})

	t.Run("Real World RollForward", func(t *testing.T) {
		dataV5Result, err := os.ReadFile("test_data/RealWorld_RollForward_v5.json")
		assert.Nil(t, err)
		var method10 CompatibleResult
		err = json.Unmarshal([]byte(dataV5Result), &method10)
		assert.Nil(t, err)
	})

	t.Run("Real World Taste Test", func(t *testing.T) {
		example, err := os.ReadFile("test_data/RealWorld_TasteTest_RollForward_v5.json")
		assert.Nil(t, err)
		var result CompatibleResult
		err = json.Unmarshal([]byte(example), &result)
		assert.Nil(t, err)
	})
}

func TestCompatibleResponse(t *testing.T) {
	t.Run("Full Response Roll Forward V5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/Response_NextBlock_v5.json")
		assert.Nil(t, err)

		var expected v5.ResponseV5
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResponsePraos
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected.ConvertToV6(), compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResponseV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)
	})

	t.Run("Full Response Roll Forward V5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/Response_NextBlock_v6.json")
		assert.Nil(t, err)

		var expected chainsync.ResponsePraos
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleResponsePraos
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ResponseV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)
	})
}

func TestDynamoDBMarshal(t *testing.T) {
	t.Run("Value v5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/Value_v5.json")
		assert.Nil(t, err)

		var compatible CompatibleValue
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		av, err := dynamodbattribute.Marshal(&compatible)
		assert.Nil(t, err)
		assert.NotNil(t, av.M)

		var got v5.ValueV5
		err = dynamodbattribute.Unmarshal(av, &got)
		assert.Nil(t, err)
		assert.EqualValues(t, compatible, got.ConvertToV6())
	})
	t.Run("Value v6", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/Value_v6.json")
		assert.Nil(t, err)

		var compatible CompatibleValue
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		av, err := dynamodbattribute.Marshal(&compatible)
		assert.Nil(t, err)
		assert.NotNil(t, av.M)

		var got v5.ValueV5
		err = dynamodbattribute.Unmarshal(av, &got)
		assert.Nil(t, err)
		assert.EqualValues(t, compatible, got.ConvertToV6())
	})
	// If we need a compatible Tx type, uncomment this test
	/*
	   t.Run("Tx v5", func(t *testing.T) {
	      rawData, err := os.ReadFile("test_data/Tx_v5.json")
	      assert.Nil(t, err)

	      var compatible CompatibleTx
	      err = json.Unmarshal(rawData, &compatible)
	      assert.Nil(t, err)

	      av, err := dynamodbattribute.Marshal(&compatible)
	      assert.Nil(t, err)
	      assert.NotNil(t, av.M)

	      var got v5.TxV5
	      err = dynamodbattribute.Unmarshal(av, &got)
	      assert.Nil(t, err)
	      assert.EqualValues(t, compatible, got.ConvertToV6())
	   })
	   t.Run("Tx v6", func(t *testing.T) {
	      rawData, err := os.ReadFile("test_data/Tx_v6.json")
	      assert.Nil(t, err)

	      var compatible CompatibleTx
	      err = json.Unmarshal(rawData, &compatible)
	      assert.Nil(t, err)

	      av, err := dynamodbattribute.Marshal(&compatible)
	      assert.Nil(t, err)
	      assert.NotNil(t, av.M)

	      var got v5.TxV5
	      err = dynamodbattribute.Unmarshal(av, &got)
	      assert.Nil(t, err)
	      assert.EqualValues(t, compatible, got.ConvertToV6())
	   })
	*/
}

func TestSerializeCompatibleValue(t *testing.T) {
	t.Run("Value v5", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/Value_v5.json")
		assert.Nil(t, err)

		var expected v5.ValueV5
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleValue
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected.ConvertToV6(), compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ValueV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, got)
	})
	t.Run("Value v6", func(t *testing.T) {
		rawData, err := os.ReadFile("test_data/Value_v6.json")
		assert.Nil(t, err)

		var expected shared.Value
		err = json.Unmarshal(rawData, &expected)
		assert.Nil(t, err)

		var compatible CompatibleValue
		err = json.Unmarshal(rawData, &compatible)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, compatible)

		bytes, err := json.Marshal(&compatible)
		assert.Nil(t, err)

		var got v5.ValueV5
		err = json.Unmarshal(bytes, &got)
		assert.Nil(t, err)

		assert.EqualValues(t, expected, got.ConvertToV6())
	})
}

func Test_ValueChecks(t *testing.T) {
	t.Run("string", func(t *testing.T) {
		equal1 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		equal2 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		if !shared.Equal(equal1, equal2) {
			t.Fatalf("%v and %v are not equal", equal1, equal2)
		}

		val1 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000001)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		val2 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		val3 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(12345678900)},
		)
		val4 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(1000000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(1234567890)},
		)
		if !shared.GreaterThanOrEqual(val1, val2) {
			t.Fatalf("%v is not greater than %v", val1, val2)
		}
		if !shared.LessThanOrEqual(val2, val1) {
			t.Fatalf("%v is not less than %v", val1, val2)
		}
		if !shared.GreaterThanOrEqual(val3, val4) {
			t.Fatalf("%v is not greater than %v", val3, val4)
		}
		if !shared.LessThanOrEqual(val4, val3) {
			t.Fatalf("%v is not less than %v", val4, val3)
		}
		if ok, err := shared.Enough(val3, val4); !ok {
			t.Fatalf("%v does not have enough assets for %v: %v", val4, val3, err)
		}
		if shared.Equal(val1, val2) {
			t.Fatalf("%v and %v are equal", val1, val2)
		}
		if shared.Equal(val3, val4) {
			t.Fatalf("%v and %v are equal", val3, val4)
		}

		val5 := shared.Add(val1, val2)
		val6 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(2000001)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(2469135780)},
		)
		if !shared.Equal(val5, val6) {
			t.Fatalf("%v is not the expected value (%v)", val5, val6)
		}

		val7 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(600000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(2345678900)},
		)
		val8 := shared.Subtract(val3, val7)
		val9 := shared.ValueFromCoins(
			shared.Coin{AssetId: shared.AdaAssetID, Amount: num.Int64(400000)},
			shared.Coin{AssetId: shared.FromSeparate("abra", "cadabra"), Amount: num.Int64(10000000000)},
		)
		if !shared.Equal(val8, val9) {
			t.Fatalf("%v is not the expected value (%v)", val8, val9)
		}
	})
}
