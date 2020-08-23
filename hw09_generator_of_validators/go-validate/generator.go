package main

import (
	"bytes"
	"errors"
	"go/format"
	"regexp"
	"strings"
	"text/template"
)

const (
	typeString = "string"
	typeInt    = "int"
)

// TField ...
type TField struct {
	Name     string
	Type     string
	Expr     map[string]string
	IsArray  bool
	NeedCast bool
}

// TData ...
type TData struct {
	Package string
	Structs map[string][]TField
}

func isSupported(t string) bool {
	return t == typeString || t == typeInt
}

func getSupportedType(t string, l map[string]string) (string, bool) {
	if isSupported(t) {
		return t, true
	}
	for k, v := range l {
		if k != t {
			continue
		}
		newType, ok := getSupportedType(v, l)
		if !ok {
			break
		}

		return newType, true
	}

	return t, false
}

func parseTag(t string) (map[string]string, error) {
	result := make(map[string]string)
	r := regexp.MustCompile(`.*validate:"(.*)"`)
	if !r.MatchString(t) {
		return result, nil
	}
	val := r.FindStringSubmatch(t)[1]
	for _, expr := range strings.Split(val, "|") {
		keyVal := strings.Split(expr, ":")
		if len(keyVal) != 2 {
			return nil, errors.New("wrong tag expression")
		}
		result[keyVal[0]] = keyVal[1]
	}

	return result, nil
}

func prepareData(vSource *ValidationSource) (*TData, error) {
	data := &TData{
		Package: vSource.Package,
		Structs: make(map[string][]TField),
	}

	for k, v := range vSource.Structs {
		fields := make([]TField, 0, len(v))
		for _, f := range v {
			expr, err := parseTag(f.Tag)
			if err != nil {
				return nil, err
			}
			if len(expr) == 0 {
				continue
			}
			fType, ok := getSupportedType(f.Type, vSource.Types)
			if !ok {
				continue
			}
			fields = append(
				fields,
				TField{
					Name:     f.Name,
					Type:     fType,
					IsArray:  f.IsArray,
					Expr:     expr,
					NeedCast: !isSupported(f.Type),
				},
			)
		}
		if len(fields) == 0 {
			continue
		}
		data.Structs[k] = fields
	}

	return data, nil
}

func buildFile(data *TData) ([]byte, error) {
	tmpl, err := template.New("tmpl").Parse(validationTemplate)
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(make([]byte, 1024))
	err = tmpl.Execute(b, data)
	if err != nil {
		return nil, err
	}
	res, err := format.Source(bytes.Trim(b.Bytes(), "\x00"))
	if err != nil {
		return nil, err
	}

	return res, nil
}
