package cubemxtobazelinternal

import (
	"fmt"
	"strings"
)

// winToNixPath converts a relative windows path to a nix style path
func winToNixPath(path string) string {
	// Simply replace all backslash with forward slash
	return strings.Replace(path, "\\", "/", -1)
}

func stripWhiteSpace(str string) string {
	return strings.Replace(str, " ", "", -1)
}

func mxFileToBazelString(file MxFile) bString {
	return bString(winToNixPath(file.Name))
}

func mxFilesToBazelStringList(files MxFiles) bStringList {
	var result bStringList
	for _, file := range files {
		result = append(result, string(winToNixPath(file.Name)))
	}
	return result
}

func ccLibraryTargetName(comp MxComponent) string {
	name := stripWhiteSpace(strings.Join([]string{comp.Class, comp.Group, comp.Subsection}, "_"))
	return name
}

// TODO:Change this to use projects rather than MxComponent i.e. func(proj Project) []CcLibraryRule
func MxProjectToCcLibraryRules(proj Project) []CcLibraryRule {
	rules := []CcLibraryRule{}
	components := proj.Components()
	for _, comp := range components {
		var files MxFiles = comp.Files
		sourceFiles := files.SourceFiles().Files()
		headerFiles := files.HeaderFiles().Files()
		asmFiles := files.AssemblyFiles().Files()

		bazelTargetComment := fmt.Sprintf("# %s  %s:%s:%s, version:%s", comp.Description, comp.Class, comp.Group, comp.Subsection, comp.Version)
		bazelSourceFiles := append(mxFilesToBazelStringList(sourceFiles), mxFilesToBazelStringList(asmFiles)...)
		bazelSourceFiles = append(bazelSourceFiles, mxFilesToBazelStringList(headerFiles)...)
		bazelHeaderFiles := mxFilesToBazelStringList(headerFiles)

		// Generated attributes
		name := bString(ccLibraryTargetName(comp))
		bazelNameAttr := attributeBString{Operand: attName, Value: name}
		bazelSourceAttr := attributeBStringList{Operand: attSrcs, Value: bazelSourceFiles}
		bazelHeaderAttr := attributeBStringList{Operand: attHdrs, Value: bazelHeaderFiles}

		// Additional attributes
		// Static linking only
		linkStaticAttr := attributeBBool{Operand: attLinkStatic, Value: true}

		// Combination of all attributes
		allAttr := attributeList{bazelNameAttr, bazelSourceAttr, bazelHeaderAttr, linkStaticAttr}
		libraryRule := CcLibraryRule{rule{Operands: allAttr, comment: comment{Comment: bazelTargetComment}}}
		rules = append(rules, libraryRule)
	}
	return rules
}

func MxProjectToCcBinaryRule(proj Project) ccBinaryRule {
	var files MxFiles = proj.ProjectFiles()
	sourceFiles := files.SourceFiles().Files()
	headerFiles := files.HeaderFiles().Files()
	asmFiles := files.AssemblyFiles().Files()

	bazelTargetComment := fmt.Sprintf("# %s, Device:%s", proj.Info().Name, proj.DeviceName())
	bazelSourceFiles := append(mxFilesToBazelStringList(sourceFiles), mxFilesToBazelStringList(asmFiles)...)
	bazelSourceFiles = append(bazelSourceFiles, mxFilesToBazelStringList(headerFiles)...)
	bazelHeaderFiles := mxFilesToBazelStringList(headerFiles)

	// Generated attributes
	name := bString(stripWhiteSpace(strings.Join([]string{proj.Info().Name, proj.DeviceName()}, "_")))
	bazelNameAttr := attributeBString{Operand: attName, Value: name}
	bazelSourceAttr := attributeBStringList{Operand: attSrcs, Value: bazelSourceFiles}
	bazelHeaderAttr := attributeBStringList{Operand: attHdrs, Value: bazelHeaderFiles}

	components := proj.Components()
	DependantTargetNames := []string{}
	for _, component := range components {
		dep := ":" + ccLibraryTargetName(component)
		DependantTargetNames = append(DependantTargetNames, dep)
	}
	bazelDepsAttr := attributeBStringList{Operand: attDeps, Value: bStringList(DependantTargetNames)}

	// Additional attributes

	// Combination of all attributes
	allAttr := attributeList{bazelNameAttr, bazelSourceAttr, bazelHeaderAttr, bazelDepsAttr}
	binaryRule := ccBinaryRule{rule{Operands: allAttr, comment: comment{Comment: bazelTargetComment}}}
	return binaryRule
}
