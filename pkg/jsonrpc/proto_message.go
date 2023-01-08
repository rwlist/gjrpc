package jsonrpc

import "encoding/json"

type MessageType string

const (
	MessageTypeRequest  MessageType = "request"
	MessageTypeResponse MessageType = "response"
	MessageTypeUnknown  MessageType = "unknown"
)

// DetectMessageType detects the type of the message from the raw bytes of JSON object.
func DetectMessageType(b []byte) MessageType {
	var tmpStruct struct {
		Method *string `json:"method"`
	}

	if err := json.Unmarshal(b, &tmpStruct); err != nil {
		return MessageTypeUnknown
	}

	if tmpStruct.Method != nil {
		return MessageTypeRequest
	}

	return MessageTypeResponse
}
