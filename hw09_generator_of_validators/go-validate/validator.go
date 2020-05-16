package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"text/template"

	"github.com/pkg/errors"
)

const (
	min   string = "min"
	max   string = "max"
	regex string = "regexp"
	in    string = "in"
	_len  string = "len"
)

type structTemplate struct {
	Name   string
	Prefix string
	Fields []Field
}

type Field struct {
	FieldName string
	FieldType string
	Regexp    string
	Len       int
	Min       int
	Max       int
	In        []string
}

func NewField(_name string, _type string) Field {
	var f = Field{
		FieldName: _name,
		FieldType: _type,
		Len:       -1,
		Min:       -1,
		Max:       -1,
	}
	return f
}

func minusFunc(x, y int) int {
	return x - y
}

var (
	minus = template.FuncMap{"minus": minusFunc}

	structValidatorTpl = template.Must(template.New("structValidatorTpl").Funcs(minus).Parse(`
func ({{.Prefix}} {{.Name}}) Validate() ([]ValidationError,error){
ve:=[]ValidationError{}

{{ range .Fields}}

//{{- .FieldName -}}

{{- if (eq .FieldType "string") -}}

{{- if (ne .Len -1)}}
if len({{$.Prefix}}.{{.FieldName}})!={{.Len}}{
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("len field {{.FieldName}} with value:%v not equal with validate value %v",{{$.Prefix}}.{{.FieldName}},{{.Len}}),
})
}
{{- end -}}

{{- if .In}}
{{$FieldName:= .FieldName}}
{{$cnt:= .In|len}}
{{$last:= (minus $cnt 1)}}
if {{range $i,$k:=.In}}{{if (eq $i $last)}}{{$.Prefix}}.{{$FieldName}}!="{{$k}}"{{else}}{{$.Prefix}}.{{$FieldName}}!="{{$k}}" && {{end}}{{end}} {
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("field {{.FieldName}} with value:%v is not one of  the values %v",{{$.Prefix}}.{{.FieldName}},"{{.In}}"),
})
}
{{- end -}}

{{- if (ne .Regexp "")}}
r:=regexp.MustCompile(` + "`" + `{{.Regexp}}` + "`" + `)
    if !r.MatchString({{$.Prefix}}.{{.FieldName}}){
    	ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("field {{.FieldName}} with value:%v not match the regex %v",{{$.Prefix}}.{{.FieldName}},` + "`" + `{{.Regexp}}` + "`" + `),
})
	}
{{- end -}}

{{- end -}}

{{- if (eq .FieldType "[]string")}}

for _,s:=range {{$.Prefix}}.{{.FieldName}}{
{{- if (ne .Len -1)}}
if len(s)!={{.Len}}{
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("len element of {{.FieldName}} with value:%v not equal with validate value %v",s,{{.Len}}),
})
}
{{- end -}}

{{- if .In}}
{{$FieldName:= .FieldName}}
{{$cnt:= .In|len}}
{{$last:= (minus $cnt 1)}}
if {{range $i,$k:=.In}}{{if (eq $i $last)}}s!="{{$k}}"{{else}}s!="{{$k}}" && {{end}}{{end}} {
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("element of  {{.FieldName}} with value:%v is not one of  the values %v",s,"{{.In}}"),
})
}
{{- end -}}

{{- if (ne .Regexp "")}}
r:=regexp.MustCompile(` + "`" + `{{.Regexp}}` + "`" + `)
    if !r.MatchString(s){
    	ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("element of  {{.FieldName}} with value:%v not match the regex %v",s,` + "`" + `{{.Regexp}}` + "`" + `),
})
	}
{{- end -}}
}
{{- end -}}

{{- if (eq .FieldType "int") }}

{{- if (ne .Min -1)}}
if {{$.Prefix}}.{{.FieldName}}<{{.Min}}{
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("field {{.FieldName}} with value:%v smaller than min value %v",{{$.Prefix}}.{{.FieldName}},{{.Min}}),
})
}
{{- end -}}

{{- if (ne .Max -1)}}
if {{$.Prefix}}.{{.FieldName}}>{{.Max}}{
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("field {{.FieldName}} with value:%v bigger than max value %v",{{$.Prefix}}.{{.FieldName}},{{.Max}}),
})
}
{{- end -}}

{{- if .In}}
{{$FieldName:= .FieldName}}
{{$cnt:= .In|len}}
{{$last:= (minus $cnt 1)}}
if {{range $i,$k:=.In}}{{if (eq $i $last)}}{{$.Prefix}}.{{$FieldName}}!={{$k}}{{else}}{{$.Prefix}}.{{$FieldName}}!={{$k}} && {{end}}{{end}} {
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("field {{.FieldName}} with value:%v is not one of  the values %v",{{$.Prefix}}.{{.FieldName}},"{{.In}}"),
})
}
{{- end -}}
{{- end -}}

{{- if (eq .FieldType "[]int") }}

for _,s:=range {{$.Prefix}}.{{.FieldName}}{
{{- if (ne .Min -1)}}
if s<{{.Min}}{
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("field {{.FieldName}} with value:%v smaller than min value %v",s,{{.Min}}),
})
}
{{- end -}}

{{- if (ne .Max -1)}}
if s>{{.Max}}{
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("field {{.FieldName}} with value:%v bigger than min value %v",s,{{.Max}}),
})
}
{{- end -}}

{{- if .In}}
{{$FieldName:= .FieldName}}
{{$cnt:= .In|len}}
{{$last:= (minus $cnt 1)}}
if {{range $i,$k:=.In}}{{if (eq $i $last)}}s!="{{$k}}"{{else}}s!="{{$k}}" && {{end}}{{end}} {
ve=append(ve,ValidationError{
Field:"{{.FieldName}}",
Err:fmt.Errorf("element of  {{.FieldName}} with value:%v is not one of  the values %v",s,"{{.In}}"),
})
}
{{- end -}}
}
{{- end -}}
{{- end}}
return ve,nil
}
`))
)

//GenValidate generate Validate functions for struct in validating .go file
func GenValidate(fname string) error {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	vname := strings.ReplaceAll(fname, filepath.Ext(fname), "_validation_generated.go")
	out, e := os.Create(vname)
	defer out.Close() //nolint:staticcheck
	if e != nil {
		return fmt.Errorf("can't create validation file. Error: %v", e)
	}

	fmt.Fprintln(out, `/*
* CODE GENERATED AUTOMATICALLY WITH go-validate
* THIS FILE SHOULD NOT BE EDITED BY HAND
*/ 
//nolint:gomnd,gofmt,goimports`)
	fmt.Fprintln(out, `package `+node.Name.Name)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `import (`)
	fmt.Fprintln(out, ` "fmt"`)
	fmt.Fprintln(out, ` "regexp"`)
	fmt.Fprintln(out, `)`)
	fmt.Fprintln(out) // empty line
	fmt.Fprintln(out, `type ValidationError struct{
Field string
Err error
}`)

	s, e := getValidatingStruct(node)
	if e != nil {
		return errors.Wrapf(e, "can't generate validation file for %v", fname)
	}
	for _, str := range s {
		if len(str.Fields) > 0 {
			if e := structValidatorTpl.Execute(out, str); e != nil {
				return errors.Wrapf(e, "can't generate validation file for %v", fname)
			}
		}
	}

	return nil
}

//getValidatingStruct parse the .go file and extract all structs to []structTemplate
func getValidatingStruct(node *ast.File) ([]structTemplate, error) {
	validateStructs := []structTemplate{}

	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			currStruct, ok := currType.Type.(*ast.StructType)
			if !ok {
				continue
			}
			structT := structTemplate{
				Name:   currType.Name.Name,
				Prefix: strings.ToLower(currType.Name.Name),
				Fields: []Field{},
			}

		FieldsLoop:
			for _, field := range currStruct.Fields.List {
				if field.Tag != nil {
					tag := reflect.StructTag(field.Tag.Value[1 : len(field.Tag.Value)-1])
					if tag.Get("validate") == "" {
						continue FieldsLoop
					}
					fieldName := field.Names[0].Name
					fieldType := "string"
					switch v := field.Type.(type) {
					case *ast.Ident:
						fieldType = v.Name
						fieldType = getTrueType(fieldType, node)
					case *ast.ArrayType:
						elemType := v.Elt.(*ast.Ident).Name
						elemType = getTrueType(elemType, node)
						fieldType = "[]" + elemType
					}
					field, err := fieldToValidator(fieldName, fieldType, tag.Get("validate"))
					if err != nil {
						return validateStructs, errors.Wrapf(err, "can't parse %v struct", currType.Name.Name)
					}
					structT.Fields = append(structT.Fields, field)
				}
			}
			validateStructs = append(validateStructs, structT)
		}
	}

	return validateStructs, nil
}

//getTrueType returns base type for type
func getTrueType(fieldType string, node *ast.File) string {
	if fieldType == "int" || fieldType == "string" {
		return fieldType
	}

	for _, f := range node.Decls {
		g, ok := f.(*ast.GenDecl)
		if !ok {
			continue
		}
		for _, spec := range g.Specs {
			currType, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			currStruct, ok := currType.Type.(*ast.Ident)
			if !ok {
				continue
			}
			if fieldType != currType.Name.Name {
				continue
			}
			fieldType = currStruct.Name
			return fieldType
		}
	}

	return fieldType
}

//fieldToValidator returns  Field type from tag
func fieldToValidator(fieldName, fieldType, tag string) (Field, error) {
	field := NewField(fieldName, fieldType)

	for _, param := range strings.Split(tag, `|`) {
		kv := strings.Split(param, `:`)
		var key string
		var value string
		switch len(kv) {
		case 1: //nolint:gomnd
			key = kv[0]
		case 2: //nolint:gomnd
			key = kv[0]
			value = kv[1]
		default:
			key = kv[0]
		}
		err := paramToValidator(&field, key, value)
		if err != nil {
			return field, errors.Wrapf(err, "can't parse %v field", fieldName)
		}
	}

	return field, nil
}

//paramToValidator parse tag
func paramToValidator(field *Field, key, value string) error {
	switch key {
	case _len:
		l, e := strconv.Atoi(value)
		if e != nil {
			return fmt.Errorf(`type validate tag "len"" should be int, but it has value %v`, value)
		}
		field.Len = l
	case min:
		m, e := strconv.Atoi(value)
		if e != nil {
			return fmt.Errorf(`type validate tag "min"" should be int, but it has value %v`, value)
		}
		field.Min = m
	case max:
		m, e := strconv.Atoi(value)
		if e != nil {
			return fmt.Errorf(`type validate tag "max"" should be int, but it has value %v`, value)
		}
		field.Max = m
	case in:
		field.In = append(field.In, strings.Split(value, ",")...)
	case regex:
		field.Regexp = value
	default:
		return fmt.Errorf("%v", "unknown validator tag key")
	}

	return nil
}
