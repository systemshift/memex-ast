package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"memex-ast/ast"
	"memex/internal/memex/repository"
)

func main() {
	// Parse command line flags
	repoPath := flag.String("repo", "", "Path to memex repository")
	sourcePath := flag.String("source", "", "Path to Go source file or directory")
	flag.Parse()

	if *repoPath == "" || *sourcePath == "" {
		fmt.Println("Usage: memex-ast -repo <repository path> -source <source path>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Open or create repository
	repo, err := repository.Open(*repoPath)
	if err != nil {
		if os.IsNotExist(err) {
			repo, err = repository.Create(*repoPath)
			if err != nil {
				log.Fatalf("Error creating repository: %v", err)
			}
		} else {
			log.Fatalf("Error opening repository: %v", err)
		}
	}
	defer repo.Close()

	// Create and register AST module
	module := ast.NewModule(repo)
	if err := repo.RegisterModule(module); err != nil {
		log.Fatalf("Error registering module: %v", err)
	}

	// Process source path
	sourceInfo, err := os.Stat(*sourcePath)
	if err != nil {
		log.Fatalf("Error accessing source path: %v", err)
	}

	if sourceInfo.IsDir() {
		// Process directory
		err = filepath.Walk(*sourcePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && filepath.Ext(path) == ".go" {
				fmt.Printf("Processing %s...\n", path)
				if err := module.ParseFile(path); err != nil {
					fmt.Printf("Error processing %s: %v\n", path, err)
				}
			}
			return nil
		})
		if err != nil {
			log.Fatalf("Error walking directory: %v", err)
		}
	} else {
		// Process single file
		if filepath.Ext(*sourcePath) != ".go" {
			log.Fatal("Source file must be a .go file")
		}
		fmt.Printf("Processing %s...\n", *sourcePath)
		if err := module.ParseFile(*sourcePath); err != nil {
			log.Fatalf("Error processing file: %v", err)
		}
	}

	fmt.Println("AST analysis complete")
}
