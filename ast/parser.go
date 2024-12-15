package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"memex/internal/memex/core"
)

// Parser handles Go source parsing
type Parser struct {
	fset  *token.FileSet
	repo  core.Repository
	files map[string]*ast.File
}

// NewParser creates a new parser
func NewParser(repo core.Repository) *Parser {
	return &Parser{
		fset:  token.NewFileSet(),
		repo:  repo,
		files: make(map[string]*ast.File),
	}
}

// ParsePath parses Go files in a path
func (p *Parser) ParsePath(path string) error {
	// Check if path is a directory
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("checking path: %w", err)
	}

	if info.IsDir() {
		// Parse directory
		pkgs, err := parser.ParseDir(p.fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parsing directory: %w", err)
		}

		// Store all files
		for _, pkg := range pkgs {
			for filename, file := range pkg.Files {
				p.files[filename] = file
			}
		}
	} else {
		// Parse single file
		file, err := parser.ParseFile(p.fset, path, nil, parser.ParseComments)
		if err != nil {
			return fmt.Errorf("parsing file: %w", err)
		}
		p.files[path] = file
	}

	return nil
}

// GetPackages returns unique package names
func (p *Parser) GetPackages() []string {
	packages := make(map[string]bool)
	for _, file := range p.files {
		if file.Name != nil {
			packages[file.Name.Name] = true
		}
	}

	result := make([]string, 0, len(packages))
	for pkg := range packages {
		result = append(result, pkg)
	}
	return result
}

// GetImports returns all imports
func (p *Parser) GetImports() map[string][]string {
	// Map of package to its imports
	imports := make(map[string][]string)

	for _, file := range p.files {
		if file.Name == nil {
			continue
		}
		pkg := file.Name.Name

		// Get imports for this file
		for _, imp := range file.Imports {
			if imp.Path != nil {
				path := imp.Path.Value
				// Remove quotes from import path
				path = path[1 : len(path)-1]
				imports[pkg] = append(imports[pkg], path)
			}
		}
	}

	return imports
}

// GetTypes returns type declarations
func (p *Parser) GetTypes() []*ast.TypeSpec {
	var types []*ast.TypeSpec

	for _, file := range p.files {
		ast.Inspect(file, func(n ast.Node) bool {
			if typeSpec, ok := n.(*ast.TypeSpec); ok {
				types = append(types, typeSpec)
			}
			return true
		})
	}

	return types
}

// GetFunctions returns function declarations
func (p *Parser) GetFunctions() []*ast.FuncDecl {
	var funcs []*ast.FuncDecl

	for _, file := range p.files {
		ast.Inspect(file, func(n ast.Node) bool {
			if funcDecl, ok := n.(*ast.FuncDecl); ok {
				funcs = append(funcs, funcDecl)
			}
			return true
		})
	}

	return funcs
}

// Position returns position information for a node
func (p *Parser) Position(node ast.Node) token.Position {
	return p.fset.Position(node.Pos())
}

// Files returns the parsed files
func (p *Parser) Files() map[string]*ast.File {
	return p.files
}

// FileSet returns the token file set
func (p *Parser) FileSet() *token.FileSet {
	return p.fset
}
