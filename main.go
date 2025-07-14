package main

import (
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

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			logger.Printf("ERR: %s", err)
		}
		HandleMessage(logger, method, contents)
	}
}

func HandleMessage(logger *log.Logger, method string, contents []byte) {
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
		logger.Printf("Sending: %s", reply)

		writer := os.Stdout
		writer.Write([]byte(reply))
	case "textDocument/didOpen":
		logger.Println("Opended Document")
	}
}

func GetLogger(filename string) *log.Logger {
	logfile, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		panic("Bad file given to logger.")
	}

	return log.New(logfile, "[bsqls]", log.Ldate|log.Ltime|log.Lshortfile)
}
