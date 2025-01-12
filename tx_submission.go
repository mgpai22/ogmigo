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

package ogmigo

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/buger/jsonparser"
)

type Response struct {
	Transaction ResponseTx `json:"transaction,omitempty"  dynamodbav:"transaction,omitempty"`
}

type ResponseTx struct {
	ID string `json:"id,omitempty" dynamodbav:"id,omitempty"`
}

type SubmitTx struct {
	Cbor string `json:"cbor"`
}

// type Response struct {
// 	Type        string
// 	Version     string
// 	ServiceName string `json:"servicename"`
// 	MethodName  string `json:"methodname"`
// 	Reflection  interface{}
// 	Result      json.RawMessage
// }

// SubmitTx submits the transaction via ogmios
// https://ogmios.dev/mini-protocols/local-tx-submission/
func (c *Client) SubmitTx(ctx context.Context, data string) (s *SubmitTxResponse, err error) {
	tx := SubmitTx{
		Cbor: data,
	}
	var (
		payload = makePayload("submitTransaction", Map{"transaction": tx}, Map{})
		raw     json.RawMessage
	)
	if err := c.query(ctx, payload, &raw); err != nil {
		return nil, fmt.Errorf("failed to submit TX: %w", err)
	}

	return readSubmitTx(raw)
}

func readSubmitTx(data []byte) (r *SubmitTxResponse, err error) {
	e, err1 := readSubmitTxError(data)
	id, err2 := readSubmitTxResult(data)
	if err1 != nil && err2 != nil {
		return nil, fmt.Errorf("could not parse submit tx response; neither error (%w) nor result (%w)", err1, err2)
	}
	if err1 == nil {
		return &SubmitTxResponse{Error: e}, nil
	}
	if err2 == nil {
		return &SubmitTxResponse{ID: id}, nil
	}
	return nil, fmt.Errorf("could not parse submit tx response: %s", string(data))
}

type SubmitTxResponse struct {
	ID    string
	Error *SubmitTxError
}

type SubmitTxError struct {
	Code    int
	Message string
	Data    json.RawMessage
}

func readSubmitTxError(data []byte) (*SubmitTxError, error) {
	value, _, _, err := jsonparser.Get(data, "error")
	if err != nil {
		return nil, fmt.Errorf("failed to parse SubmitTx error: %w %s", err, data)
	}
	var e SubmitTxError
	if err := json.Unmarshal(value, &e); err != nil {
		return nil, fmt.Errorf("failed to parse SubmitTx error: %w %s", err, data)
	}
	return &e, nil
}

func readSubmitTxResult(data []byte) (string, error) {
	value, dataType, _, err := jsonparser.Get(data, "result")
	if err != nil {
		return "", fmt.Errorf("failed to parse SubmitTx response: %w %s", err, string(data))
	}

	switch dataType {
	case jsonparser.Object:
		var result struct {
			Transaction struct {
				ID string
			}
		}
		if err := json.Unmarshal(value, &result); err != nil {
			return "", fmt.Errorf("failed to parse SubmitTx response: %w", err)
		}
		return result.Transaction.ID, nil
	default:
		return "", fmt.Errorf("failed to parser SubmitTx response: %w", err)
	}
}
