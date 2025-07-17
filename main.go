package main

import (
	"bsqls/analyzer"
	"bsqls/lsp"
	"bsqls/rpc"
	"bufio"
	"encoding/json"
	"io"
	"log"
	"os"
)

func main() {
	logger := GetLogger("/home/karidus/personal/bsqls/log.txt")
	logger.Println("Logger Started.")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := analyzer.NewState()

	writer := os.Stdout

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("ERR: %s", err)
		}
		HandleMessage(logger, writer, state, method, contents)
	}
}

func HandleMessage(logger *log.Logger, writer io.Writer, state analyzer.State, method string, contents []byte) {
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

		initialize := lsp.NewInitializeResponse(request.ID)
		WriteResponse(initialize, writer)

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

	case "textDocument/hover":
		var request lsp.HoverRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed textDocument/hover: %s", err)
		}

		hover := state.HandleHover(request.ID,
			request.Params.TextDocument.Uri,
			request.Params.Position,
		)

		WriteResponse(hover, writer)

	case "textDocument/definition":
		var request lsp.DefinitionRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			logger.Printf("Failed textDocument/definition: %s", err)
		}

		definition := state.HandleDefinition(
			request.ID,
			request.Params.TextDocument.Uri,
			request.Params.Position,
		)

		WriteResponse(definition, writer)
	}

}

func WriteResponse(msg any, writer io.Writer) {
	reply := rpc.EncodeMessage(msg)
	writer.Write([]byte(reply))
}

func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Bad file given to logger.")
	}

	return log.New(logfile, "[bsqls]", log.Ldate|log.Ltime|log.Lshortfile)
}
