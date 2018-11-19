package cubemxtobazelinternal

import "fmt"

func mxFileToBazelString(file MxFile) bString {
	return bString(file.Name)
}

func mxFilesToBazelStringList(files MxFiles) bStringList {
	var result bStringList
	for _, file := range files {
		result = append(result, string(file.Name))
	}
	return result
}

func mxComponentToCcLibraryRule(comp MxComponent) ccLibraryRule {
	var files MxFiles = comp.Files
	sourceFiles := files.SourceFiles().Files()
	headerFiles := files.HeaderFiles().Files()
	asmFiles := files.AssemblyFiles().Files()

	bazelTargetComment := fmt.Sprintf("# %s, %s:%s, version:%s", comp.Description, comp.Class, comp.Group, comp.Version)
	bazelSourceFiles := append(mxFilesToBazelStringList(sourceFiles), mxFilesToBazelStringList(asmFiles)...)
	bazelSourceFiles = append(bazelSourceFiles, mxFilesToBazelStringList(headerFiles)...)
	bazelHeaderFiles := mxFilesToBazelStringList(headerFiles)

	// Generated attributes
	bazelSourceAttr := attributeBStringList{Operand: attSrcs, Value: bazelSourceFiles}
	bazelHeaderAttr := attributeBStringList{Operand: attHdrs, Value: bazelHeaderFiles}

	// Additional attributes
	// Static linking only
	linkStaticAttr := attributeBBool{Operand: attLinkStatic, Value: true}

	// Combination of all attributes
	allAttr := attributeList{bazelSourceAttr, bazelHeaderAttr, linkStaticAttr}
	libraryRule := ccLibraryRule{rule{Operands: allAttr, comment: comment{Comment: bazelTargetComment}}}
	return libraryRule
}

func mxComponentToCcBinaryRule(comp MxComponent) ccBinaryRule {
	var files MxFiles = comp.Files
	sourceFiles := files.SourceFiles().Files()
	headerFiles := files.HeaderFiles().Files()
	asmFiles := files.AssemblyFiles().Files()

	bazelTargetComment := fmt.Sprintf("# %s, %s:%s, version:%s", comp.Description, comp.Class, comp.Group, comp.Version)
	bazelSourceFiles := append(mxFilesToBazelStringList(sourceFiles), mxFilesToBazelStringList(asmFiles)...)
	bazelSourceFiles = append(bazelSourceFiles, mxFilesToBazelStringList(headerFiles)...)
	bazelHeaderFiles := mxFilesToBazelStringList(headerFiles)

	// Generated attributes
	bazelSourceAttr := attributeBStringList{Operand: attSrcs, Value: bazelSourceFiles}
	bazelHeaderAttr := attributeBStringList{Operand: attHdrs, Value: bazelHeaderFiles}

	// Additional attributes

	// Combination of all attributes
	allAttr := attributeList{bazelSourceAttr, bazelHeaderAttr}
	binaryRule := ccBinaryRule{rule{Operands: allAttr, comment: comment{Comment: bazelTargetComment}}}
	return binaryRule
}
