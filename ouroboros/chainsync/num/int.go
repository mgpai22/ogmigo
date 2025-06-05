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

package num

import (
	"errors"
	"fmt"
	"math/big"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// Int is a wrapper of sorts around big.Int. One of the intentions is to prevent users
// from being directly exposed to the big.Int API and usage of BigInt(). This is due
// to big.Int introducing bugs if not careful with its usage and its API.
type Int big.Int

func Int64(v int64) Int {
	bi := big.NewInt(v)
	return Int(*bi)
}

func Uint64(v uint64) Int {
	bi := big.NewInt(0).SetUint64(v)
	return Int(*bi)
}

func New(s string) (Int, bool) {
	bi, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return Int{}, false
	}

	return Int(*bi), true
}

func (i Int) Add(that Int) Int {
	sum := big.NewInt(0).Add(i.BigInt(), that.BigInt())
	return Int(*sum)
}

func (i Int) BigInt() *big.Int {
	bi := big.Int(i)
	return &bi
}

func (i Int) BigFloat() *big.Float {
	return big.NewFloat(0).SetInt(i.BigInt())
}

func (i Int) Int() int {
	return int(i.BigInt().Int64())
}

func (i Int) Int64() int64 {
	return i.BigInt().Int64()
}

func (i Int) Uint64() uint64 {
	return i.BigInt().Uint64()
}

func (i Int) MarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	item.N = aws.String(i.BigInt().String())
	return nil
}

func (i Int) MarshalJSON() ([]byte, error) {
	s := i.BigInt().String()
	return []byte(s), nil
}

func (i Int) String() string {
	return i.BigInt().String()
}

func (i Int) Sub(that Int) Int {
	sum := big.NewInt(0).Sub(i.BigInt(), that.BigInt())
	return Int(*sum)
}

func (i Int) Mul(that Int) Int {
	product := big.NewInt(0).Mul(i.BigInt(), that.BigInt())
	return Int(*product)
}

func (i Int) Div(that Int) Int {
	quotient := big.NewInt(0).Div(i.BigInt(), that.BigInt())
	return Int(*quotient)
}

func (i Int) LessThan(that Int) bool {
	return i.BigInt().Cmp(that.BigInt()) == -1
}

func (i Int) GreaterThan(that Int) bool {
	return i.BigInt().Cmp(that.BigInt()) == 1
}

func (i Int) Equal(that Int) bool {
	return i.BigInt().Cmp(that.BigInt()) == 0
}

func (i *Int) UnmarshalDynamoDBAttributeValue(item *dynamodb.AttributeValue) error {
	if aws.BoolValue(item.NULL) {
		return nil
	}
	if item.N == nil {
		return errors.New("unable to unmarshal invalid Int: N not set")
	}

	s := aws.StringValue(item.N)
	v, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return fmt.Errorf("failed to parse number, %v", s)
	}

	number := Int(*v)
	*i = number

	return nil
}

func (i *Int) UnmarshalJSON(data []byte) error {
	if data == nil {
		return nil
	}

	s := string(data)
	v, ok := big.NewInt(0).SetString(s, 10)
	if !ok {
		return fmt.Errorf("failed to parse number, %v", s)
	}

	number := Int(*v)
	*i = number

	return nil
}
