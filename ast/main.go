package ast

import (
	"github.com/systemshift/memex/pkg/sdk/types"
)

// NewModule is the exported plugin function that creates a new AST module
func NewModule(repo types.Repository) types.Module {
	return &Module{
		repo:     repo,
		parser:   NewParser(repo),
		analyzer: NewAnalyzer(repo),
		builder:  NewGraphBuilder(repo),
	}
}
