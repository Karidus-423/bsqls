package lsp

const /*Change*/ (
	NONE = iota
	FULL
	INCREMENTAL
)

type InitializeRequest struct {
	Request
	Params InitializeRequestParams `json:"params"`
}

type InitializeRequestParams struct {
	ClientInfo *ClientInfo `json:"clientInfo"`
	//TODO: Add more params to initialize
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResponse struct {
	Response
	Result InitializeResult `json:"result"`
}

type InitializeResult struct {
	Capabilities ServerCapabilities `json:"capabilities"`
	ServerInfo   ServerInfo         `json:"serverInfo"`
}

type TextDocumentSyncOptions struct {
	OpenClose bool `json:"openClose"`
	Change    int  `json:"change"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type ServerCapabilities struct {
	TextDocumentSync   TextDocumentSyncOptions `json:"textDocumentSync"`
	HoverProvider      bool                    `json:"hoverProvider"`
	DefinitionProvider bool                    `json:"definitionProvider"`
	PositionEncoding   string                  `json:"positionEncoding"`
}

func NewInitializeResponse(id int) InitializeResponse {
	return InitializeResponse{
		Response: Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: InitializeResult{
			Capabilities: ServerCapabilities{
				PositionEncoding: "utf-16",
				TextDocumentSync: TextDocumentSyncOptions{
					OpenClose: true,
					//TODO: Handle INCREMENTAL
					Change: FULL,
				},
				HoverProvider:      true,
				DefinitionProvider: true,
			},
			ServerInfo: ServerInfo{
				Name:    "bsqls",
				Version: "0.1.0",
			},
		},
	}
}
