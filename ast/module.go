package ast

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/systemshift/memex/pkg/memex/core"
)

// Node types for AST nodes
const (
	NodeTypePackage   = "ast.package"
	NodeTypeFunction  = "ast.function"
	NodeTypeStruct    = "ast.struct"
	NodeTypeInterface = "ast.interface"
	NodeTypeMethod    = "ast.method"
	NodeTypeField     = "ast.field"
	NodeTypeImport    = "ast.import"
)

// Link types for AST relationships
const (
	LinkTypeCalls      = "ast.calls"      // Function calls another function
	LinkTypeImplements = "ast.implements" // Struct implements interface
	LinkTypeContains   = "ast.contains"   // Package contains type/func, struct contains field
	LinkTypeImports    = "ast.imports"    // Package imports another package
	LinkTypeExtends    = "ast.extends"    // Struct embeds another struct
	LinkTypeReferences = "ast.references" // Type references another type
)

// Module implements the memex module interface for AST analysis
type Module struct {
	repo core.Repository
}

// NewModule creates a new AST module
func NewModule(repo core.Repository) *Module {
	return &Module{
		repo: repo,
	}
}

// ID returns the module identifier
func (m *Module) ID() string {
	return "ast"
}

// Name returns the human-readable name
func (m *Module) Name() string {
	return "AST Analysis"
}

// Description returns the module description
func (m *Module) Description() string {
	return "Analyzes Go source code and stores AST structure in memex"
}

// Capabilities returns the module capabilities
func (m *Module) Capabilities() []core.ModuleCapability {
	return []core.ModuleCapability{} // Empty slice instead of nil
}

// ValidateNodeType checks if a node type is valid for this module
func (m *Module) ValidateNodeType(nodeType string) bool {
	switch nodeType {
	case NodeTypePackage, NodeTypeFunction, NodeTypeStruct,
		NodeTypeInterface, NodeTypeMethod, NodeTypeField,
		NodeTypeImport:
		return true
	default:
		return false
	}
}

// ValidateLinkType checks if a link type is valid for this module
func (m *Module) ValidateLinkType(linkType string) bool {
	switch linkType {
	case LinkTypeCalls, LinkTypeImplements, LinkTypeContains,
		LinkTypeImports, LinkTypeExtends, LinkTypeReferences:
		return true
	default:
		return false
	}
}

// ValidateMetadata validates module-specific metadata
func (m *Module) ValidateMetadata(meta map[string]interface{}) error {
	return nil // No special validation yet
}

// ParseFile parses a Go source file and stores its AST structure
func (m *Module) ParseFile(filename string) error {
	// Create file set for position information
	fset := token.NewFileSet()

	// Parse the Go source file
	file, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
	if err != nil {
		return fmt.Errorf("parsing file: %w", err)
	}

	// Create package node
	pkgMeta := map[string]interface{}{
		"module":   m.ID(),
		"filename": filename,
	}
	pkgID, err := m.repo.AddNode([]byte(file.Name.Name), NodeTypePackage, pkgMeta)
	if err != nil {
		return fmt.Errorf("adding package node: %w", err)
	}

	// Process imports
	for _, imp := range file.Imports {
		importPath := imp.Path.Value
		importMeta := map[string]interface{}{
			"module": m.ID(),
			"path":   importPath,
		}
		if imp.Name != nil {
			importMeta["alias"] = imp.Name.Name
		}

		// Create import node
		importID, err := m.repo.AddNode([]byte(importPath), NodeTypeImport, importMeta)
		if err != nil {
			return fmt.Errorf("adding import node: %w", err)
		}

		// Link import to package
		err = m.repo.AddLink(pkgID, importID, LinkTypeImports, nil)
		if err != nil {
			return fmt.Errorf("adding import link: %w", err)
		}
	}

	// Process declarations
	for _, decl := range file.Decls {
		if err := m.processDecl(pkgID, fset, decl); err != nil {
			return fmt.Errorf("processing declaration: %w", err)
		}
	}

	return nil
}

// processDecl processes an AST declaration
func (m *Module) processDecl(pkgID string, fset *token.FileSet, decl ast.Decl) error {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		// Process function declaration
		funcMeta := map[string]interface{}{
			"module": m.ID(),
			"name":   d.Name.Name,
			"pos":    fset.Position(d.Pos()).String(),
		}
		if d.Recv != nil {
			funcMeta["receiver"] = getTypeString(d.Recv.List[0].Type)
		}

		funcID, err := m.repo.AddNode([]byte(d.Name.Name), NodeTypeFunction, funcMeta)
		if err != nil {
			return fmt.Errorf("adding function node: %w", err)
		}

		// Link function to package
		err = m.repo.AddLink(pkgID, funcID, LinkTypeContains, nil)
		if err != nil {
			return fmt.Errorf("adding function link: %w", err)
		}

	case *ast.GenDecl:
		for _, spec := range d.Specs {
			switch s := spec.(type) {
			case *ast.TypeSpec:
				// Process type declaration
				switch t := s.Type.(type) {
				case *ast.StructType:
					// Process struct
					structMeta := map[string]interface{}{
						"module": m.ID(),
						"name":   s.Name.Name,
						"pos":    fset.Position(s.Pos()).String(),
					}
					structID, err := m.repo.AddNode([]byte(s.Name.Name), NodeTypeStruct, structMeta)
					if err != nil {
						return fmt.Errorf("adding struct node: %w", err)
					}

					// Link struct to package
					err = m.repo.AddLink(pkgID, structID, LinkTypeContains, nil)
					if err != nil {
						return fmt.Errorf("adding struct link: %w", err)
					}

					// Process struct fields
					for _, field := range t.Fields.List {
						fieldMeta := map[string]interface{}{
							"module": m.ID(),
							"type":   getTypeString(field.Type),
							"pos":    fset.Position(field.Pos()).String(),
						}
						if len(field.Names) > 0 {
							fieldMeta["name"] = field.Names[0].Name
						}

						fieldID, err := m.repo.AddNode([]byte(fieldMeta["type"].(string)), NodeTypeField, fieldMeta)
						if err != nil {
							return fmt.Errorf("adding field node: %w", err)
						}

						// Link field to struct
						err = m.repo.AddLink(structID, fieldID, LinkTypeContains, nil)
						if err != nil {
							return fmt.Errorf("adding field link: %w", err)
						}
					}

				case *ast.InterfaceType:
					// Process interface
					interfaceMeta := map[string]interface{}{
						"module": m.ID(),
						"name":   s.Name.Name,
						"pos":    fset.Position(s.Pos()).String(),
					}
					interfaceID, err := m.repo.AddNode([]byte(s.Name.Name), NodeTypeInterface, interfaceMeta)
					if err != nil {
						return fmt.Errorf("adding interface node: %w", err)
					}

					// Link interface to package
					err = m.repo.AddLink(pkgID, interfaceID, LinkTypeContains, nil)
					if err != nil {
						return fmt.Errorf("adding interface link: %w", err)
					}
				}
			}
		}
	}

	return nil
}

// getTypeString returns a string representation of an AST type
func getTypeString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return "*" + getTypeString(t.X)
	case *ast.SelectorExpr:
		return getTypeString(t.X) + "." + t.Sel.Name
	case *ast.ArrayType:
		return "[]" + getTypeString(t.Elt)
	case *ast.MapType:
		return "map[" + getTypeString(t.Key) + "]" + getTypeString(t.Value)
	default:
		return fmt.Sprintf("%T", expr)
	}
}
