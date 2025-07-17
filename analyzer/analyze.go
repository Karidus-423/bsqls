package analyzer

import (
	"bsqls/lsp"
	"fmt"
)

type State struct {
	Documents map[string]string //Map of doc -> text
}

func NewState() State {
	return State{Documents: map[string]string{}}
}

func (s *State) OpenDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) UpdateDocument(uri, text string) {
	s.Documents[uri] = text
}

func (s *State) HandleHover(id int, uri string, position lsp.Position) lsp.HoverResponse {
	//TODO: Look up definition locatin
	//TODO: Do this move cursor to actual location
	doc := s.Documents[uri]

	return lsp.HoverResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.HoverResult{
			Contents: fmt.Sprintf("File: %s, Chars: %d", uri, len(doc)),
		},
	}
}

func (s *State) HandleDefinition(id int, uri string, position lsp.Position) lsp.DefinitionResponse {
	//TODO: Look up type
	//TODO: Do this dynamically

	return lsp.DefinitionResponse{
		Response: lsp.Response{
			RPC: "2.0",
			ID:  &id,
		},
		Result: lsp.Location{
			Uri: uri,
			Range: lsp.Range{
				Start: lsp.Position{
					Line:      position.Line - 1,
					Character: 0,
				},
				End: lsp.Position{
					Line:      position.Line - 1,
					Character: 5,
				},
			},
		},
	}
}
