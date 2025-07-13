package rpc_test

import (
	"bsqls/rpc"
	"strconv"
	"testing"
)

type EncodingExample struct {
	Testing bool
}

func TestEncode(t *testing.T) {
	expected := "Content-Length: 16\r\n\r\n{\"Testing\":true}"
	actual := rpc.EncodeMessage(EncodingExample{Testing: true})
	if expected != actual {
		t.Fatalf("Expected: %s, Actual: %s", expected, actual)
	}
}

func TestDecode(t *testing.T) {
	test_method := "{\"method\": \"textDocument/completion\"}"
	msg_in := "Content-Length: " + strconv.Itoa(len(test_method)) +
		"\r\n\r\n" + test_method
	method, content, err := rpc.DecodeMessage([]byte(msg_in))
	content_len := len(content)

	if err != nil {
		t.Fatal(err)
	}

	if content_len != len(test_method) {
		t.Fatalf("Expected: %d, Actual: %d", len(test_method), content_len)
	}

	if method != "textDocument/completion" {
		t.Fatalf("Expected: textDocument/completion, Actual: %s", method)
	}

}
