package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io"
)

// Field ...
type Field struct {
	Name    string
	Type    string
	Tag     string
	IsArray bool
}

// ValidationSource ...
type ValidationSource struct {
	Package string
	Structs map[string][]Field
	Types   map[string]string
}

func parseSource(r io.Reader) (*ValidationSource, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", r, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	result := ValidationSource{
		Package: f.Name.Name,
		Structs: make(map[string][]Field),
		Types:   make(map[string]string),
	}

	ast.Inspect(f, func(n ast.Node) bool {
		x, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}
		switch nodeType := x.Type.(type) {
		case *ast.Ident:
			result.Types[x.Name.Name] = nodeType.Name
		case *ast.StructType:
			sFields := make([]Field, 0, len(nodeType.Fields.List))
			for _, field := range nodeType.Fields.List {
				if field.Tag == nil {
					continue
				}
				structField := Field{
					Name: field.Names[0].Name,
					Tag:  field.Tag.Value,
				}
				switch fieldType := field.Type.(type) {
				case *ast.Ident:
					structField.Type = fieldType.Name
				case *ast.ArrayType:
					structField.IsArray = true

					elementType, ok := fieldType.Elt.(*ast.Ident)
					if ok {
						structField.Type = elementType.Name
					}
				}
				sFields = append(sFields, structField)
			}
			if len(sFields) > 0 {
				result.Structs[x.Name.Name] = sFields
			}
		}

		return true
	})

	return &result, nil
}
