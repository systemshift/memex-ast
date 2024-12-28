package ast

import (
	"fmt"
	"path/filepath"

	"github.com/systemshift/memex/pkg/module"
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

// AST implements Go source code analysis
type AST struct {
	*module.Base
	parser   *Parser
	analyzer *Analyzer
	builder  *GraphBuilder
}

// New creates a new AST module
func New() module.Module {
	m := &AST{
		Base: module.NewBase("ast", "AST Analysis", "Analyzes Go source code structure"),
	}

	// Add AST-specific commands
	m.Base.AddCommand(module.Command{
		Name:        "parse",
		Description: "Parse Go source files",
		Usage:       "ast parse <path>",
		Args:        []string{"path"},
	})

	m.Base.AddCommand(module.Command{
		Name:        "types",
		Description: "Show type relationships",
		Usage:       "ast types [type-name]",
	})

	m.Base.AddCommand(module.Command{
		Name:        "calls",
		Description: "Show function call graph",
		Usage:       "ast calls [function-name]",
	})

	m.Base.AddCommand(module.Command{
		Name:        "impls",
		Description: "Find interface implementations",
		Usage:       "ast impls <interface-name>",
		Args:        []string{"interface-name"},
	})

	m.Base.AddCommand(module.Command{
		Name:        "deps",
		Description: "Show package dependencies",
		Usage:       "ast deps [package-path]",
	})

	return m
}

// Init initializes the module
func (m *AST) Init(repo module.Repository) error {
	if err := m.Base.Init(repo); err != nil {
		return err
	}

	// Initialize components
	m.parser = NewParser(repo)
	m.analyzer = NewAnalyzer(repo)
	m.builder = NewGraphBuilder(repo)

	// Connect components
	m.analyzer.SetParser(m.parser)
	m.builder.SetAnalyzer(m.analyzer)

	return nil
}

// HandleCommand handles module commands
func (m *AST) HandleCommand(cmd string, args []string) error {
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
		return m.Base.HandleCommand(cmd, args)
	}
}

// Parse parses Go source files
func (m *AST) Parse(path string) error {
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
func (m *AST) ShowTypes(typeName string) error {
	nodes, err := m.Base.GetNode(typeName)
	if err != nil {
		return fmt.Errorf("querying nodes: %w", err)
	}

	if nodes.Type == NodeTypeStruct || nodes.Type == NodeTypeInterface {
		fmt.Printf("%s: %s\n", nodes.Meta["name"], nodes.Type)
		if methods, ok := nodes.Meta["methods"].([]string); ok {
			fmt.Printf("  Methods: %v\n", methods)
		}
		if embedded, ok := nodes.Meta["embedded"].([]string); ok {
			fmt.Printf("  Embedded: %v\n", embedded)
		}
	}

	return nil
}

// ShowCalls shows function call graph
func (m *AST) ShowCalls(funcName string) error {
	node, err := m.Base.GetNode(funcName)
	if err != nil {
		return fmt.Errorf("getting function: %w", err)
	}

	if node.Type == NodeTypeFunction {
		fmt.Printf("%s:\n", node.Meta["name"])
		links, err := m.Base.GetLinks(node.ID)
		if err != nil {
			return fmt.Errorf("getting links: %w", err)
		}
		for _, link := range links {
			if link.Type == LinkTypeCalls {
				callee, err := m.Base.GetNode(link.Target)
				if err != nil {
					continue
				}
				fmt.Printf("  calls: %s\n", callee.Meta["name"])
			}
		}
	}

	return nil
}

// ShowImplementations shows interface implementations
func (m *AST) ShowImplementations(interfaceName string) error {
	node, err := m.Base.GetNode(interfaceName)
	if err != nil {
		return fmt.Errorf("getting interface: %w", err)
	}

	if node.Type == NodeTypeInterface {
		links, err := m.Base.GetLinks(node.ID)
		if err != nil {
			return fmt.Errorf("getting links: %w", err)
		}
		for _, link := range links {
			if link.Type == LinkTypeImplements {
				impl, err := m.Base.GetNode(link.Source)
				if err != nil {
					continue
				}
				fmt.Printf("%s implements %s\n", impl.Meta["name"], interfaceName)
			}
		}
	}

	return nil
}

// ShowDependencies shows package dependencies
func (m *AST) ShowDependencies(pkgPath string) error {
	node, err := m.Base.GetNode(pkgPath)
	if err != nil {
		return fmt.Errorf("getting package: %w", err)
	}

	if node.Type == NodeTypePackage {
		fmt.Printf("%s:\n", node.Meta["path"])
		links, err := m.Base.GetLinks(node.ID)
		if err != nil {
			return fmt.Errorf("getting links: %w", err)
		}
		for _, link := range links {
			if link.Type == LinkTypeImports {
				dep, err := m.Base.GetNode(link.Target)
				if err != nil {
					continue
				}
				fmt.Printf("  imports: %s\n", dep.Meta["path"])
			}
		}
	}

	return nil
}
