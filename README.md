# Memex AST Module

A Go source code analysis module for Memex that creates a queryable graph of code structure and relationships.

## Installation

Install directly from git:
```bash
memex module install https://github.com/systemshift/memex-ast.git
memex module enable ast
```

## Usage

### Parse Go files:
```bash
# Parse a single file
memex ast parse file.go

# Parse a directory
memex ast parse ./pkg/...
```

### View relationships:
```bash
# Show type relationships
memex ast types

# Show call graph
memex ast calls

# Find interface implementations
memex ast impls io.Reader

# Show package dependencies
memex ast deps
```

### Query examples:
```bash
# Find all implementations of an interface
memex query "type:ast.struct -[ast.implements]-> {name: 'io.Reader'}"

# Show function call chain
memex query "type:ast.function -[ast.calls*]-> {name: 'main'}"

# Find unused types
memex query "type:ast.struct -[!ast.uses]->"

# Show package dependencies
memex query "type:ast.package -[ast.imports]-> type:ast.package"
```

## Module Structure

The module analyzes Go source code using:
- `go/ast`: AST parsing
- `go/parser`: Go source parsing
- `go/token`: Token handling

### Node Types:
- `ast.package`: Package declarations
- `ast.function`: Functions and methods
- `ast.struct`: Struct definitions
- `ast.interface`: Interface definitions
- `ast.field`: Struct fields
- `ast.method`: Interface methods
- `ast.import`: Package imports

### Link Types:
- `ast.calls`: Function calls
- `ast.implements`: Interface implementations
- `ast.contains`: Package/type containment
- `ast.imports`: Package dependencies
- `ast.embeds`: Type embedding
- `ast.uses`: Type usage

## Development

This is a memex module and requires memex to be installed. It's designed to be installed directly from git rather than run as a standalone binary.

### Testing:
```bash
go test ./...
```

### Local Development:
```bash
# Clone the repository
git clone https://github.com/systemshift/memex-ast.git

# Install locally for testing
memex module install ./memex-ast

# Make changes and reinstall to test
memex module remove ast
memex module install ./memex-ast
