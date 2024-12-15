package ast

import (
	"fmt"
	"go/ast"

	"memex/internal/memex/core"
)

// TypeInfo holds type analysis information
type TypeInfo struct {
	Name        string
	IsStruct    bool
	IsInterface bool
	Methods     []string
	Embedded    []string
}

// Analyzer analyzes code relationships
type Analyzer struct {
	repo   core.Repository
	parser *Parser
	types  map[string]*TypeInfo
	calls  map[string][]string
	uses   map[string][]string
}

// NewAnalyzer creates a new analyzer
func NewAnalyzer(repo core.Repository) *Analyzer {
	return &Analyzer{
		repo:  repo,
		types: make(map[string]*TypeInfo),
		calls: make(map[string][]string),
		uses:  make(map[string][]string),
	}
}

// SetParser sets the parser
func (a *Analyzer) SetParser(p *Parser) {
	a.parser = p
}

// Analyze analyzes parsed files
func (a *Analyzer) Analyze() error {
	if a.parser == nil {
		return fmt.Errorf("parser not set")
	}

	// Analyze types
	if err := a.analyzeTypes(); err != nil {
		return fmt.Errorf("analyzing types: %w", err)
	}

	// Analyze function calls
	if err := a.analyzeCalls(); err != nil {
		return fmt.Errorf("analyzing calls: %w", err)
	}

	// Analyze type usage
	if err := a.analyzeUses(); err != nil {
		return fmt.Errorf("analyzing uses: %w", err)
	}

	return nil
}

// analyzeTypes analyzes type declarations
func (a *Analyzer) analyzeTypes() error {
	for _, typeSpec := range a.parser.GetTypes() {
		info := &TypeInfo{
			Name: typeSpec.Name.Name,
		}

		switch t := typeSpec.Type.(type) {
		case *ast.StructType:
			info.IsStruct = true
			// Analyze struct fields and embedded types
			for _, field := range t.Fields.List {
				if field.Names == nil {
					// Embedded type
					info.Embedded = append(info.Embedded, getTypeString(field.Type))
				}
			}

		case *ast.InterfaceType:
			info.IsInterface = true
			// Analyze interface methods
			for _, method := range t.Methods.List {
				if len(method.Names) > 0 {
					info.Methods = append(info.Methods, method.Names[0].Name)
				}
			}
		}

		a.types[info.Name] = info
	}

	return nil
}

// analyzeCalls analyzes function calls
func (a *Analyzer) analyzeCalls() error {
	for _, file := range a.parser.Files() {
		ast.Inspect(file, func(n ast.Node) bool {
			if call, ok := n.(*ast.CallExpr); ok {
				if fun, ok := call.Fun.(*ast.Ident); ok {
					caller := getCurrentFunction(n)
					if caller != "" {
						a.calls[caller] = append(a.calls[caller], fun.Name)
					}
				}
			}
			return true
		})
	}
	return nil
}

// analyzeUses analyzes type usage
func (a *Analyzer) analyzeUses() error {
	for _, file := range a.parser.Files() {
		ast.Inspect(file, func(n ast.Node) bool {
			if typeExpr, ok := n.(ast.Expr); ok {
				if typeName := getTypeString(typeExpr); typeName != "" {
					if info, exists := a.types[typeName]; exists {
						context := getCurrentFunction(n)
						if context != "" {
							a.uses[context] = append(a.uses[context], info.Name)
						}
					}
				}
			}
			return true
		})
	}
	return nil
}

// getCurrentFunction returns the enclosing function name
func getCurrentFunction(node ast.Node) string {
	for n := node; n != nil; n = findParent(n) {
		if fn, ok := n.(*ast.FuncDecl); ok {
			return fn.Name.Name
		}
	}
	return ""
}

// findParent finds the parent AST node
func findParent(node ast.Node) ast.Node {
	var parent ast.Node
	ast.Inspect(node, func(n ast.Node) bool {
		if n == node {
			return false
		}
		parent = n
		return true
	})
	return parent
}

// getTypeString returns a string representation of a type
func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.SelectorExpr:
		return getTypeString(t.X) + "." + t.Sel.Name
	default:
		return ""
	}
}
