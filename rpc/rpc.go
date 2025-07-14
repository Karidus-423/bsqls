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

var ErrDelimiterNotFound = errors.New("Did not find \"\r\n\r\n\" in message.")

func GetHeaderInfo(msg []byte) ([]byte, []byte, int, error) {
	header, content, found := bytes.Cut(msg, []byte{'\r', '\n', '\r', '\n'})
	if found == false {
		return nil, nil, 0, ErrDelimiterNotFound
	}

	header_val_str := header[len("Content-Length: "):]
	header_val, err := strconv.Atoi(string(header_val_str))
	if err != nil {
		return nil, nil, 0, err
	}

	return header, content, header_val, nil
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
	_, content, header_val, err := GetHeaderInfo(msg)
	if err != nil {
		return "n/a", nil, err
	}

	var base_msg BaseMessage
	if err := json.Unmarshal(content[:header_val], &base_msg); err != nil {
		return "n/a", nil, err
	}

	return base_msg.Method, content[:header_val], nil
}

func Split(data []byte, _ bool) (advance int, token []byte, err error) {
	header, content, header_val, err := GetHeaderInfo(data)
	if errors.Is(err, ErrDelimiterNotFound) {
		return 0, nil, nil
	}

	//Need more time to keep scanning
	if len(content) < header_val {
		return 0, nil, nil
	}

	total_len := len(header) + 4 + header_val
	return total_len, data[:total_len], nil
}
