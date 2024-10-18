package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"slices"
)

type Emitter struct {
	pkg     string
	imports []string
	code    string
}

func NewEmiiter() *Emitter {
	return &Emitter{
		pkg:     "main",
		imports: []string{},
		code:    "",
	}
}
func (e *Emitter) AddImport(module string) {
	if !slices.Contains(e.imports, module) {
		e.imports = append(e.imports, module)
	}
}

func (e *Emitter) EmitCode(code string) {
	e.code += code
}

func (e *Emitter) SaveToFile(outputPath string) {
	final := fmt.Sprintf("package %v\n\n", e.pkg)
	fmt.Println("Emitting code...")

	if len(e.imports) > 0 {
		final += fmt.Sprintf("import (\n")
		for _, mod := range e.imports {
			final += fmt.Sprintf("\t\"%v\"\n", mod)
		}
		final += fmt.Sprintf(")\n\n")
	}
	final += fmt.Sprintf("%v", e.code)

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		err := os.MkdirAll(path.Dir(outputPath), 0700)

		if err != nil {
			log.Fatalf("Error creating output dir (%s): %v\n", outputPath, err)
		}
	}

	err := os.WriteFile(outputPath, []byte(final), 0664)
	if err != nil {
		log.Fatalf("Error writing code to output file (%s): %v\n", outputPath, err)
	}

	fmt.Println("Emitting has finished")
}
