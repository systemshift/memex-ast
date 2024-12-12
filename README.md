# memex-ast

A Memex module for analyzing Go source code and storing its AST (Abstract Syntax Tree) structure in a Memex repository.

## Features

- Parses Go source files into a graph structure
- Tracks relationships between code elements:
  - Package dependencies
  - Function calls
  - Type implementations
  - Struct composition
- Stores source locations for IDE integration
- Preserves code structure metadata

## Usage

```bash
# Analyze a single file
memex-ast -repo /path/to/repo.mx -source /path/to/file.go

# Analyze a directory
memex-ast -repo /path/to/repo.mx -source /path/to/project
```

## Node Types

- ast.package: Go packages
- ast.function: Functions and methods
- ast.struct: Struct definitions
- ast.interface: Interface definitions
- ast.method: Interface methods
- ast.field: Struct fields
- ast.import: Package imports

## Link Types

- ast.calls: Function calls
- ast.implements: Interface implementations
- ast.contains: Package/struct containment
- ast.imports: Package imports
- ast.extends: Struct embedding
- ast.references: Type references

## Development Status

This is an MVP (Minimum Viable Product) implementation. Future improvements:
- Function call analysis
- Interface implementation detection
- Type dependency tracking
- Code change impact analysis
