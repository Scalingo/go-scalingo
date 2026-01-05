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
	genSpecialize(types)
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

func genSpecialize(types []string) {
	specializeTemplate := `package scalingo

// Do not edit, generated with 'go generate'

import (
	"encoding/json"

	"github.com/Scalingo/go-scalingo/v9/debug"
)

func (pev *Event) Specialize() DetailedEvent {
	var e DetailedEvent
	ev := *pev
	switch ev.Type {
{{ range . }}	case {{ . }}:
		e = &{{ . }}Type{Event: ev}
{{end}}	default:
		return pev
	}
	err := json.Unmarshal(pev.RawTypeData, e.TypeDataPtr())
	if err != nil {
		debug.Printf("error reading the data: %+v\n", err)
		return pev
	}
	return e
}
`

	tpl, err := template.New("specialize-gen").Parse(specializeTemplate)
	if err != nil {
		panic(err)
	}

	fd, err := os.OpenFile("events_specialize.go", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer fd.Close()

	err = tpl.Execute(fd, types)
	if err != nil {
		panic(err)
	}
}
