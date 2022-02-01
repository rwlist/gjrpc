package jsonrpc

import (
	"encoding/json"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/assert"
)

func TestRequestUnmarshal(t *testing.T) {
	samples := []string{
		`{"jsonrpc": "2.0", "method": "auth.oauth", "params": null, "id": 1}`,
		`{"jsonrpc": "2.0", "method": "auth.oauth", "params": null, "id": "1"}`,
		`{"jsonrpc": "2.0", "method": "auth.oauth", "params": "abc", "id": "1"}`,
		`{"jsonrpc": "2.0", "method": "auth.oauth", "params": "abc"}`,
		`{"jsonrpc": "2.0", "method": "auth.oauth", "params": []}`,
		`{"jsonrpc": "2.0", "method": "auth.oauth", "params": {}}`,
	}

	for _, data := range samples {
		var req Request
		assert.NoError(t, json.Unmarshal([]byte(data), &req))
		spew.Dump(req)
	}
}

func TestResponseMarshal(t *testing.T) {
	samples := []Response{
		{
			Version: "2.0",
			Result:  nil,
			Error:   &Error{Code: -32700, Message: "Parse error"},
			ID:      nil,
		}, {
			Version: "2.0",
			Result: struct {
				Token string `json:"token"`
			}{Token: "abc"},
			Error: nil,
			ID:    nil,
		},
	}

	for _, data := range samples {
		res, err := json.Marshal(data)
		assert.NoError(t, err)
		spew.Dump(res)
	}
}
