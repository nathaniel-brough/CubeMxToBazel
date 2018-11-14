package cubemxtobazelinternal

import (
	"bytes"
	"log"
	"strconv"
	"text/template"
)

// CC rule string templates
const (
	ccBinaryTemplate  = "{{.Comment}}\ncc_binary({{.Operands}})\n"
	ccLibraryTemplate = "{{.Comment}}\ncc_library({{.Operands}})\n"
	ccImportTemplate  = "{{.Comment}}\ncc_import({{.Operands}})\n"
)

// Operand templates
const (
	operandTemplate = "{{.Operand}}={{.Value}},{{.Comment}}"
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

// Comments
type comment struct {
	Comment string
}

// Bazel Types
type bBool bool

func (b bBool) String() string {
	if b {
		return `"True"`
	}
	return `"False"`
}

type bString string

func (s bString) String() string {
	return strconv.Quote(string(s))
}

type bStringList []string

func (sList bStringList) String() string {
	output := "["
	for _, s := range sList {
		output = output + strconv.Quote(string(s)) + ","
	}
	output = output + "]"
	return output
}

type attribute interface {
	Attribute() attributeBase
}

type attributeBase struct {
	comment
	Operand string
	Value   string
}

func (at attributeBase) String() string {
	templateName := "operand"
	t := template.Must(template.New(templateName).Parse(operandTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, at)
	if err != nil {
		log.Println("Executing template:", err)
	}
	return output.String()
}

type attributeBString struct {
	comment
	Operand string
	Value   bString
}

func (at attributeBString) Attribute() attributeBase {
	return attributeBase{Operand: at.Operand, Value: at.Value.String()}
}

type attributeBStringList struct {
	comment
	Operand string
	Value   bStringList
}

func (at attributeBStringList) Attribute() attributeBase {
	return attributeBase{Operand: at.Operand, Value: at.Value.String(), comment: comment{at.Comment}}
}

type attributeBBool struct {
	comment
	Operand string
	Value   bBool
}

func (at attributeBBool) Attribute() attributeBase {
	return attributeBase{Operand: at.Operand, Value: at.Value.String(), comment: comment{at.Comment}}
}

type attributeList []attribute

func (list attributeList) String() string {
	var output string
	for _, s := range list {
		output += s.Attribute().String()
	}
	return output
}

type rule struct {
	comment
	Operands attributeList
}

type ccLibraryRule struct {
	rule
}

func (r ccLibraryRule) String() string {
	templateName := "cc_library"
	t := template.Must(template.New(templateName).Parse(ccLibraryTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, r)
	if err != nil {
		log.Println("Executing template:", err)
	}
	return output.String()
}

type ccBinaryRule struct {
	rule
}

func (r ccBinaryRule) String() string {
	templateName := "cc_binary"
	t := template.Must(template.New(templateName).Parse(ccBinaryTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, r)
	if err != nil {
		log.Println("Executing template:", err)
	}
	return output.String()
}
