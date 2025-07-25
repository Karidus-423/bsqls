package lsp

type DefinitionRequest struct {
	Request
	Params DefinitionParams `json:"params"`
}

type DefinitionParams struct {
	TextDocumentPositionParams
}

type DefinitionResponse struct {
	Response
	Result Location `json:"result"`
}

type Location struct {
	Uri   string `json:"uri"`
	Range Range  `json:"range"`
}
