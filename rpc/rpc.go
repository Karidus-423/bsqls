package rpc

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type BaseMessage struct {
	Method string `json:"method"`
}

func EncodeMessage(msg any) string {
	content, err := json.Marshal(msg)
	if err != nil {
		//TODO: Handle this better.
		panic(err)
	}

	return fmt.Sprintf("Content-Length: %d\r\n\r\n%s", len(content), content)
}

func DecodeMessage(msg []byte) (string, []byte, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})

	if found == false {
		return "n/a", nil, errors.New("Did not find \"\r\n\r\n\" in message.")
	}

	header_val_str := header[len("Content-Length: "):]
	header_val, err := strconv.Atoi(string(header_val_str))
	if err != nil {
		return "n/a", nil, err
	}

	var base_msg BaseMessage
	if err := json.Unmarshal(content[:header_val], &base_msg); err != nil {
		return "n/a", nil, err
	}

	return base_msg.Method, content[:header_val], nil
}
