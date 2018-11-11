package cubemxtobazelinternal

import (
	"bytes"
	"log"
	"strconv"
	"text/template"
)

// CC rule string templates
const (
	ccBinaryTemplate  = "{{.Comment}}\ncc_binary(\n{{.Operands}}\n)"
	ccLibraryTemplate = "{{.Comment}}\ncc_binary(\n{{.Operands}}\n)"
	ccImportTemplate  = "{{.Comment}}\ncc_binary(\n{{.Operands}}\n)"
)

// Operand templates
const (
	operandTemplate = "  {{.Operand}}={{.Value}},{{.Comment}}\n"
)

// Common Attribute List
const (
	attName               = "name"
	attData               = "data"
	attVisibility         = "visibility"
	attToolchains         = "toolchains"
	attDeps               = "deps"
	attDeprecation        = "deprecation"
	attFeatures           = "features"
	attLicenses           = "licenses"
	attCompatibleWith     = "compatible_with"
	attDistribs           = "distribs"
	attExecCompatibleWith = "exec_compatible_with"
	attRestrictedTo       = "restricted_to"
)

// CC Attribute List
const (
	attSrcs               = "srcs"
	attHdrs               = "hdrs"
	attCopts              = "copts"
	attDefines            = "defines"
	attIncludes           = "includes"
	attLinkOpts           = "linkopts"
	attLinkShared         = "linkshared"
	attLinkStatic         = "linkstatic"
	attStaticLibrary      = "static_library"
	attSharedLibrary      = "shared_library"
	attMalloc             = "malloc"
	attNoCopts            = "nocopts"
	attStamp              = "stamp"
	attStripIncludePrefix = "strip_include_prefix"
)

// Bazel Types
type bBool bool

func (b bBool) String() string {
	if b {
		return "True"
	}
	return "False"
}

type bString string

func (s bString) String() string {
	return strconv.Quote(string(s))
}

type bStringList []string

func (sList bStringList) String() string {
	var output string
	for _, s := range sList {
		output = "\"" + output + s
	}
	return output
}

type parameter interface {
	Parameter() operandBase
}

type operandBase struct {
	Operand string
	Value   string
}

func (op operandBase) String() string {
	templateName := "operand"
	t := template.Must(template.New(templateName).Parse(operandTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, op)
	if err != nil {
		log.Println("Executing template:", err)
	}
	return output.String()
}

type operandBString struct {
	Operand string
	Value   bString
}

func (op operandBString) Parameter() operandBase {
	return operandBase{Operand: op.Operand, Value: op.Value.String()}
}

type operandBBool struct {
	Operand string
	Value   bBool
}

func (op operandBBool) Parameter() operandBase {
	return operandBase{Operand: op.Operand, Value: op.Value.String()}
}
