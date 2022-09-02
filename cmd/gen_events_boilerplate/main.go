package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"text/template"
)

func main() {
	types := getEventTypes()
	genBoilerplate(types)
}

func getEventTypes() []string {
	var types []string
	fset := token.NewFileSet()
	// Parse the file and create an AST
	file, err := parser.ParseFile(fset, "events_structs.go", nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}
	ast.Inspect(file, func(n ast.Node) bool {
		eventValueSpec, ok := n.(*ast.ValueSpec)
		if ok && eventValueSpec != nil {
			typeIdent, ok := eventValueSpec.Type.(*ast.Ident)
			if ok && typeIdent.Name == "EventTypeName" {
				types = append(types, eventValueSpec.Names[0].Name)
			}
		}
		return true
	})
	return types
}

func genBoilerplate(types []string) {
	boilerplateTemplate := `package scalingo

// Do not edit, generated with 'go generate'
{{ range . }}

func (e *{{ . }}Type) TypeDataPtr() interface{} {
	return &e.TypeData
}{{ end }}
`

	tpl, err := template.New("boilerplate-gen").Parse(boilerplateTemplate)
	if err != nil {
		panic(err)
	}

	fd, err := os.OpenFile("events_boilerplate.go", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	err = tpl.Execute(fd, types)
	if err != nil {
		panic(err)
	}
}
