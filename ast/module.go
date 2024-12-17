package ast

import (
	"fmt"
	"path/filepath"

	"github.com/systemshift/memex/pkg/sdk/types"
)

// Node types
const (
	NodeTypePackage   = "ast.package"
	NodeTypeFunction  = "ast.function"
	NodeTypeStruct    = "ast.struct"
	NodeTypeInterface = "ast.interface"
	NodeTypeField     = "ast.field"
	NodeTypeMethod    = "ast.method"
	NodeTypeImport    = "ast.import"
)

// Link types
const (
	LinkTypeCalls      = "ast.calls"      // Function calls
	LinkTypeImplements = "ast.implements" // Interface implementations
	LinkTypeContains   = "ast.contains"   // Containment
	LinkTypeImports    = "ast.imports"    // Package imports
	LinkTypeEmbeds     = "ast.embeds"     // Type embedding
	LinkTypeUses       = "ast.uses"       // Type usage
)

// Module implements AST analysis
type Module struct {
	repo     types.Repository
	parser   *Parser
	analyzer *Analyzer
	builder  *GraphBuilder
}

// ID returns module identifier
func (m *Module) ID() string {
	return "ast"
}

// Name returns human-readable name
func (m *Module) Name() string {
	return "AST Analysis"
}

// Description returns module description
func (m *Module) Description() string {
	return "Analyzes Go source code structure and relationships"
}

// Commands returns available module commands
func (m *Module) Commands() []types.ModuleCommand {
	return []types.ModuleCommand{
		{
			Name:        "parse",
			Description: "Parse Go source files",
			Usage:       "ast parse <path>",
		},
		{
			Name:        "types",
			Description: "Show type relationships",
			Usage:       "ast types [type-name]",
		},
		{
			Name:        "calls",
			Description: "Show function call graph",
			Usage:       "ast calls [function-name]",
		},
		{
			Name:        "impls",
			Description: "Find interface implementations",
			Usage:       "ast impls <interface-name>",
		},
		{
			Name:        "deps",
			Description: "Show package dependencies",
			Usage:       "ast deps [package-path]",
		},
	}
}

// ValidateNodeType validates node types
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

// ValidateLinkType validates link types
func (m *Module) ValidateLinkType(linkType string) bool {
	switch linkType {
	case LinkTypeCalls, LinkTypeImplements, LinkTypeContains,
		LinkTypeImports, LinkTypeEmbeds, LinkTypeUses:
		return true
	default:
		return false
	}
}

// ValidateMetadata validates module-specific metadata
func (m *Module) ValidateMetadata(meta map[string]interface{}) error {
	// No special validation needed for AST metadata
	return nil
}

// HandleCommand handles module commands
func (m *Module) HandleCommand(cmd string, args []string) error {
	switch cmd {
	case "parse":
		if len(args) < 1 {
			return fmt.Errorf("path required")
		}
		return m.Parse(args[0])

	case "types":
		typeName := ""
		if len(args) > 0 {
			typeName = args[0]
		}
		return m.ShowTypes(typeName)

	case "calls":
		funcName := ""
		if len(args) > 0 {
			funcName = args[0]
		}
		return m.ShowCalls(funcName)

	case "impls":
		if len(args) < 1 {
			return fmt.Errorf("interface name required")
		}
		return m.ShowImplementations(args[0])

	case "deps":
		pkgPath := ""
		if len(args) > 0 {
			pkgPath = args[0]
		}
		return m.ShowDependencies(pkgPath)

	default:
		return fmt.Errorf("unknown command: %s", cmd)
	}
}

// Parse parses Go source files
func (m *Module) Parse(path string) error {
	// Resolve absolute path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("resolving path: %w", err)
	}

	// Parse files
	if err := m.parser.ParsePath(absPath); err != nil {
		return fmt.Errorf("parsing files: %w", err)
	}

	// Analyze relationships
	if err := m.analyzer.Analyze(); err != nil {
		return fmt.Errorf("analyzing code: %w", err)
	}

	// Build graph
	if err := m.builder.Build(); err != nil {
		return fmt.Errorf("building graph: %w", err)
	}

	return nil
}

// ShowTypes shows type relationships
func (m *Module) ShowTypes(typeName string) error {
	// Query type relationships
	var query string
	if typeName != "" {
		query = fmt.Sprintf(`type:ast.struct name:"%s"`, typeName)
	} else {
		query = "type:ast.struct"
	}

	// TODO: Execute query
	fmt.Printf("Types query: %s\n", query)
	return nil
}

// ShowCalls shows function call graph
func (m *Module) ShowCalls(funcName string) error {
	// Query call graph
	var query string
	if funcName != "" {
		query = fmt.Sprintf(`type:ast.function name:"%s" -[ast.calls*]->`, funcName)
	} else {
		query = "type:ast.function -[ast.calls]-> type:ast.function"
	}

	// TODO: Execute query
	fmt.Printf("Calls query: %s\n", query)
	return nil
}

// ShowImplementations shows interface implementations
func (m *Module) ShowImplementations(interfaceName string) error {
	// Query implementations
	query := fmt.Sprintf(`type:ast.struct -[ast.implements]-> {type:ast.interface name:"%s"}`, interfaceName)

	// TODO: Execute query
	fmt.Printf("Implementations query: %s\n", query)
	return nil
}

// ShowDependencies shows package dependencies
func (m *Module) ShowDependencies(pkgPath string) error {
	// Query dependencies
	var query string
	if pkgPath != "" {
		query = fmt.Sprintf(`{type:ast.package path:"%s"} -[ast.imports]-> type:ast.package`, pkgPath)
	} else {
		query = "type:ast.package -[ast.imports]-> type:ast.package"
	}

	// TODO: Execute query
	fmt.Printf("Dependencies query: %s\n", query)
	return nil
}
