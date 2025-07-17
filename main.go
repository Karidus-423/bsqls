package main

import (
	"bsqls/analyzer"
	"bsqls/lsp"
	"bsqls/rpc"
	"bufio"
	"encoding/json"
	"log"
	"os"
)

func main() {
	logger := GetLogger("/home/karidus/personal/bsqls/log.txt")
	logger.Println("Logger Started.")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analyzer.NewState()

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("ERR: %s", err)
		}
		HandleMessage(logger, state, method, contents)
	}
}

func HandleMessage(logger *log.Logger, state analyzer.State, method string, contents []byte) {
	logger.Printf("Method Recieved: %s", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Error Parsing: %s", err)
		}

		logger.Printf("Connection Established: %s ver.%s",
			request.Params.ClientInfo.Name,
			request.Params.ClientInfo.Version,
		)

		msg := lsp.NewInitializeResponse(request.ID)
		reply := rpc.EncodeMessage(msg)

		writer := os.Stdout
		writer.Write([]byte(reply))

	case "textDocument/didOpen":
		var request lsp.DidOpenTextDocumentNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed textDocument/didOpen: %s", err)
		}
		logger.Printf("Opened File: %s", request.Params.TextDocument.Uri)

		state.OpenDocument(request.Params.TextDocument.Uri, request.Params.TextDocument.Text)

	case "textDocument/didChange":
		var request lsp.TextDocumentDidChangeNotification
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed textDocument/didChange: %s", err)
		}
		logger.Printf("Changed: %s", request.Params.TextDocument.Uri)
		for _, change := range request.Params.ContentChanges {
			state.UpdateDocument(
				request.Params.TextDocument.Uri,
				change.Text)
		}
	}
}

func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Bad file given to logger.")
	}

	return log.New(logfile, "[bsqls]", log.Ldate|log.Ltime|log.Lshortfile)
}
