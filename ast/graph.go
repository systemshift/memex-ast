package ast

import (
	"fmt"

	"github.com/systemshift/memex/pkg/module"
)

// GraphBuilder builds the Memex graph from analysis
type GraphBuilder struct {
	repo     module.Repository
	analyzer *Analyzer
	// Node ID tracking
	packages  map[string]string
	types     map[string]string
	functions map[string]string
}

// NewGraphBuilder creates a new graph builder
func NewGraphBuilder(repo module.Repository) *GraphBuilder {
	return &GraphBuilder{
		repo:      repo,
		packages:  make(map[string]string),
		types:     make(map[string]string),
		functions: make(map[string]string),
	}
}

// SetAnalyzer sets the analyzer
func (g *GraphBuilder) SetAnalyzer(a *Analyzer) {
	g.analyzer = a
}

// Build builds the graph from analysis
func (g *GraphBuilder) Build() error {
	if g.analyzer == nil {
		return fmt.Errorf("analyzer not set")
	}

	// Build package nodes
	if err := g.buildPackages(); err != nil {
		return fmt.Errorf("building packages: %w", err)
	}

	// Build type nodes
	if err := g.buildTypes(); err != nil {
		return fmt.Errorf("building types: %w", err)
	}

	// Build function nodes
	if err := g.buildFunctions(); err != nil {
		return fmt.Errorf("building functions: %w", err)
	}

	// Build relationships
	if err := g.buildRelationships(); err != nil {
		return fmt.Errorf("building relationships: %w", err)
	}

	return nil
}

// buildPackages creates package nodes
func (g *GraphBuilder) buildPackages() error {
	for _, pkg := range g.analyzer.parser.GetPackages() {
		meta := map[string]interface{}{
			"module": "ast",
			"name":   pkg,
		}
		id, err := g.repo.AddNode([]byte(pkg), NodeTypePackage, meta)
		if err != nil {
			return fmt.Errorf("adding package node: %w", err)
		}
		g.packages[pkg] = id
	}
	return nil
}

// buildTypes creates type nodes
func (g *GraphBuilder) buildTypes() error {
	for name, info := range g.analyzer.types {
		nodeType := NodeTypeStruct
		if info.IsInterface {
			nodeType = NodeTypeInterface
		}

		meta := map[string]interface{}{
			"module":   "ast",
			"name":     name,
			"methods":  info.Methods,
			"embedded": info.Embedded,
		}

		id, err := g.repo.AddNode([]byte(name), nodeType, meta)
		if err != nil {
			return fmt.Errorf("adding type node: %w", err)
		}
		g.types[name] = id
	}
	return nil
}

// buildFunctions creates function nodes
func (g *GraphBuilder) buildFunctions() error {
	for _, fn := range g.analyzer.parser.GetFunctions() {
		meta := map[string]interface{}{
			"module": "ast",
			"name":   fn.Name.Name,
		}
		if fn.Recv != nil {
			meta["receiver"] = getTypeString(fn.Recv.List[0].Type)
		}

		id, err := g.repo.AddNode([]byte(fn.Name.Name), NodeTypeFunction, meta)
		if err != nil {
			return fmt.Errorf("adding function node: %w", err)
		}
		g.functions[fn.Name.Name] = id
	}
	return nil
}

// buildRelationships creates relationships between nodes
func (g *GraphBuilder) buildRelationships() error {
	// Build function calls
	for caller, callees := range g.analyzer.calls {
		callerID, exists := g.functions[caller]
		if !exists {
			continue
		}

		for _, callee := range callees {
			calleeID, exists := g.functions[callee]
			if !exists {
				continue
			}

			if err := g.repo.AddLink(callerID, calleeID, LinkTypeCalls, nil); err != nil {
				return fmt.Errorf("adding call link: %w", err)
			}
		}
	}

	// Build type relationships
	for typeName, info := range g.analyzer.types {
		typeID, exists := g.types[typeName]
		if !exists {
			continue
		}

		// Add embedded type relationships
		for _, embedded := range info.Embedded {
			embeddedID, exists := g.types[embedded]
			if !exists {
				continue
			}

			if err := g.repo.AddLink(typeID, embeddedID, LinkTypeEmbeds, nil); err != nil {
				return fmt.Errorf("adding embed link: %w", err)
			}
		}
	}

	// Build type usage relationships
	for context, uses := range g.analyzer.uses {
		contextID, exists := g.functions[context]
		if !exists {
			continue
		}

		for _, used := range uses {
			usedID, exists := g.types[used]
			if !exists {
				continue
			}

			if err := g.repo.AddLink(contextID, usedID, LinkTypeUses, nil); err != nil {
				return fmt.Errorf("adding use link: %w", err)
			}
		}
	}

	return nil
}
