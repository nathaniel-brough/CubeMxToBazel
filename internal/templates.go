package cubemxtobazelinternal

import (
	"bytes"
	"log"
	"strconv"
	"text/template"
)

// CC rule string templates
const (
	ccBinaryTemplate  = "{{.Comment}}\ncc_binary({{.Keys}})\n"
	ccLibraryTemplate = "{{.Comment}}\ncc_library({{.Keys}})\n"
	ccImportTemplate  = "{{.Comment}}\ncc_import({{.Keys}})\n"
)

// Operand templates
const (
	operandTemplate = "{{.Key}}={{.Value}},{{.Comment}}"
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

// bBool Bazel Boolean Type
type bBool bool

func (b bBool) String() string {
	if b {
		return `True`
	}
	return `False`
}

// bString Bazel string type
type bString string

func (s bString) String() string {
	return strconv.Quote(string(s))
}

// bStringList Bazel string list type
type bStringList []string

// bStringList Stringer implementation
func (sList bStringList) String() string {
	output := "["
	for _, s := range sList {
		output = output + strconv.Quote(string(s)) + ","
	}
	output = output + "]"
	return output
}

// attribute describes bazel target attributes
type attribute interface {
	Attribute() attributeBase
}

// attributeBase is the most basic raw type that all other types are converted to before being converted to a string
type attributeBase struct {
	comment
	Key   string
	Value string
}

// attributeBase conversion to a string
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

// attributeBString conversion to a string
type attributeBString struct {
	comment
	Key   string
	Value bString
}

// attributeBString conversion to base type
func (at attributeBString) Attribute() attributeBase {
	return attributeBase{Key: at.Key, Value: at.Value.String()}
}

// attributeBStringList Bazel List of strings attribute
type attributeBStringList struct {
	comment
	Key   string
	Value bStringList
}

// attributeBStringList
func (at attributeBStringList) Attribute() attributeBase {
	return attributeBase{Key: at.Key, Value: at.Value.String(), comment: comment{at.Comment}}
}

type attributeBBool struct {
	comment
	Key   string
	Value bBool
}

func (at attributeBBool) Attribute() attributeBase {
	return attributeBase{Key: at.Key, Value: at.Value.String(), comment: comment{at.Comment}}
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
	Keys attributeList
}

type CcLibraryRule struct {
	rule
}

func (r CcLibraryRule) String() string {
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

type ccImportRule struct {
	rule
}

func (r ccImportRule) String() string {
	templateName := "cc_import"
	t := template.Must(template.New(templateName).Parse(ccImportTemplate))
	var output bytes.Buffer
	err := t.Execute(&output, r)
	if err != nil {
		log.Println("Executing template:", err)
	}
	return output.String()
}
