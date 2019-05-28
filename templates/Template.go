package templates

import (
	"bytes"
	"errors"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/weibreeze/breeze-generator/core"
)

//CodeTemplate names
const (
	All  = "all"
	Php  = "php"
	Java = "java"
	Go   = "go"
	Cpp  = "cpp"
)

var (
	instances = map[string]core.CodeTemplate{
		Php:  &PHPTemplate{},
		Java: &JavaTemplate{},
		Go:   &GoTemplate{},
		Cpp:  &CppTemplate{},
	}
)

//GetTemplate : get template by name
func GetTemplate(names string) ([]core.CodeTemplate, error) {
	templates := make([]core.CodeTemplate, 0, len(instances))
	if names == All {
		for _, t := range instances {
			templates = append(templates, t)
		}
		return templates, nil
	}
	arr := strings.Split(names, ",")
	for _, name := range arr {
		t := instances[strings.ToLower(strings.TrimSpace(name))]
		if t == nil {
			return nil, errors.New("can not find template: " + name)
		}
		templates = append(templates, t)
	}
	return templates, nil
}

//Register : register a new CodeTemplate
func Register(template core.CodeTemplate) {
	instances[template.Name()] = template
}

func sortFields(message *core.Message) []*core.Field {
	fields := make([]*core.Field, 0, len(message.Fields))
	keys := make([]int, 0, len(message.Fields))
	for key := range message.Fields {
		keys = append(keys, key)
	}
	sort.Ints(keys)
	for _, index := range keys {
		fields = append(fields, message.Fields[index])
	}
	return fields
}

func sortUnique(a []string) []string {
	m := make(map[string]bool, len(a))
	for _, v := range a {
		m[v] = true
	}
	result := make([]string, 0, len(m))
	for key := range m {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}

func firstUpper(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func firstLower(s string) string {
	return strings.ToLower(s[:1]) + s[1:]
}

func writeGenerateComment(buf *bytes.Buffer, schemaName string) {
	buf.WriteString("/*\n * Generated by breeze-generator\n * Schema: " + schemaName + "\n * Date: " + time.Now().Format("2006/1/2") + "\n */\n")
}

func withPackageDir(fileName string, schema *core.Schema) string {
	if schema.Options[core.WithPackageDir] != "" {
		return strings.ReplaceAll(schema.Package, ".", string(os.PathSeparator)) + string(os.PathSeparator) + fileName
	}
	return fileName
}
