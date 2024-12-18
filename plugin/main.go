package main

import (
    "github.com/systemshift/memex/pkg/sdk/types"
    "github.com/systemshift/memex-ast/ast"
)

func NewModule(repo types.Repository) types.Module {
    return ast.NewModule(repo)
}
