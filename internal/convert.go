package cubemxtobazelinternal

import (
	"encoding/xml"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
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

// ProjectInit parses a raw gpdsc file and initialises the project structure
func ProjectInit(gpdsc []byte) Project {
	project := projectImpl{}

	// Used for multithreading unmarshall
	var wg sync.WaitGroup
	unmarshal := func(f interface{}) {
		defer wg.Done()
		err := xml.Unmarshal(gpdsc, f)
		if err != nil {
			log.Fatal("gpdsc unmarshall failed:\n", err)
		}
	}

	// Unmarshal project info
	wg.Add(1)
	go unmarshal(&project.info)
	// Unmarshal project requirements
	wg.Add(1)
	go unmarshal(&project.requirements)
	// Unmarshal project options
	wg.Add(1)
	go unmarshal(&project.options)
	// Unmarshal project generator
	wg.Add(1)
	go unmarshal(&project.generator)
	// Unmarshal project components
	wg.Add(1)
	go unmarshal(&project.components)
	// Unmarshal project conditions
	wg.Add(1)
	go unmarshal(&project.conditions)

	wg.Wait()

	// Resolve target naming conflicts
	project.components = project.components.resolveComponentConflict()
	// Homogenise descriptions for each component
	project.components.homogeniseDescriptions()

	return project
}

func ccLibraryTargetName(comp MxComponent) string {
	name := stripWhiteSpace(strings.Join([]string{comp.Class, comp.Group, comp.Subsection}, "_"))
	if name[len(name)-1] == '_' {
		name = name[:len(name)-1]
	}
	return name
}

func getLibraryIncludePaths(files MxFiles) []string {
	headerFiles := files.HeaderFiles().Files()
	// map[directories]no_of_occurances
	includeDir := make(map[string]int)
	for _, file := range headerFiles {
		directory := filepath.Dir(string(mxFileToBazelString(file)))
		count, exists := includeDir[directory]
		if !exists {
			includeDir[directory] = 1
		} else {
			includeDir[directory] = count + 1
		}
	}
	result := []string{}
	for directory := range includeDir {
		result = append(result, directory)
		fmt.Println(directory)
	}
	return result
}

// MxProjectToCcLibraryRules converts the project components into bazel cc_library rules
func MxProjectToCcLibraryRules(proj Project) []CcLibraryRule {
	rules := []CcLibraryRule{}
	components := proj.Components()
	for _, comp := range components {
		var files MxFiles = comp.Files
		// TODO: Make this generic so that IAR compiler is supported
		gccFiltered := append(files.Condition("GCC Toolchain").Files(), files.Condition("").Files()...)
		sourceFiles := gccFiltered.SourceFiles().Files()
		headerFiles := gccFiltered.HeaderFiles().Files()
		asmFiles := gccFiltered.AssemblyFiles().Files()
		includeDirectories := getLibraryIncludePaths(gccFiltered)

		bazelTargetComment := fmt.Sprintf("# %s  %s:%s:%s, version:%s", comp.Description, comp.Class, comp.Group, comp.Subsection, comp.Version)
		bazelSourceFiles := append(mxFilesToBazelStringList(sourceFiles), mxFilesToBazelStringList(asmFiles)...)
		bazelSourceFiles = append(bazelSourceFiles, mxFilesToBazelStringList(headerFiles)...)
		bazelHeaderFiles := mxFilesToBazelStringList(headerFiles)

		// Generated attributes
		name := bString(ccLibraryTargetName(comp))
		bazelNameAttr := attributeBString{Key: attName, Value: name}
		bazelSourceAttr := attributeBStringList{Key: attSrcs, Value: bazelSourceFiles}
		bazelHeaderAttr := attributeBStringList{Key: attHdrs, Value: bazelHeaderFiles}
		bazelIncludeAttr := attributeBStringList{Key: attIncludes, Value: bStringList(includeDirectories)}

		// Additional attributes
		// Static linking only
		linkStaticAttr := attributeBBool{Key: attLinkStatic, Value: true}
		// TODO: Remove this attribute when ARM_GCC_NONE can use -system flag without extern "C" guards, ETA:"Q4 2019"
		bazelStripIncludeAttr := attributeBString{Key: attStripIncludePrefix, Value: bString(".")}

		// Combination of all attributes
		allAttr := attributeList{bazelNameAttr, bazelSourceAttr, bazelHeaderAttr, bazelIncludeAttr, bazelStripIncludeAttr, linkStaticAttr}
		libraryRule := CcLibraryRule{rule{Keys: allAttr, comment: comment{Comment: bazelTargetComment}}}
		rules = append(rules, libraryRule)
	}
	return rules
}

// MxProjectToCcBinaryRule converts the project files into a bazel cc_binary rule
func MxProjectToCcBinaryRule(proj Project) CcBinaryRule {
	var files MxFiles = proj.ProjectFiles()
	// TODO: Make this generic so that IAR can be used
	gccFiltered := append(files.Condition("GCC Toolchain").Files(), files.Condition("").Files()...)
	sourceFiles := gccFiltered.SourceFiles().Files()
	headerFiles := gccFiltered.HeaderFiles().Files()
	asmFiles := gccFiltered.AssemblyFiles().Files()

	bazelTargetComment := "# Main target"
	bazelSourceFiles := append(mxFilesToBazelStringList(sourceFiles), mxFilesToBazelStringList(asmFiles)...)
	bazelSourceFiles = append(bazelSourceFiles, mxFilesToBazelStringList(headerFiles)...)

	// Generated attributes
	name := bString("main")
	bazelNameAttr := attributeBString{Key: attName, Value: name}
	bazelSourceAttr := attributeBStringList{Key: attSrcs, Value: bazelSourceFiles}

	components := proj.Components()
	DependantTargetNames := []string{}
	for _, component := range components {
		dep := ":" + ccLibraryTargetName(component)
		DependantTargetNames = append(DependantTargetNames, dep)
	}
	bazelDepsAttr := attributeBStringList{Key: attDeps, Value: bStringList(DependantTargetNames)}

	// Additional attributes

	// Combination of all attributes
	allAttr := attributeList{bazelNameAttr, bazelSourceAttr, bazelDepsAttr}
	binaryRule := CcBinaryRule{rule{Keys: allAttr, comment: comment{Comment: bazelTargetComment}}}
	return binaryRule
}

// ResolveComponentConflict
func (comp MxComponents) resolveComponentConflict() MxComponents {
	components := comp.Components
	TargetNames := make(map[string]MxComponent)
	for _, comp := range components {
		name := ccLibraryTargetName(comp)
		val, exist := TargetNames[name]
		if !exist {
			TargetNames[name] = comp
		} else {
			// Remove version seperators
			currentVersion := strings.Replace(val.Version, ".", "", -1)
			updatedVersion := strings.Replace(comp.Version, ".", "", -1)
			// Convert to integers
			currVerInt, err := strconv.Atoi(currentVersion)
			if err != nil {
				log.Fatal("Malformed version number: ", currentVersion, "Name:", name)
			}
			updVerInt, err := strconv.Atoi(currentVersion)
			if err != nil {
				log.Fatal("Malformed version number: ", updatedVersion, "Name:", name)
			}
			// Compare version numbers
			if updVerInt > currVerInt {
				TargetNames[name] = comp

			}
		}
	}
	result := []MxComponent{}
	for _, val := range TargetNames {
		result = append(result, val)
	}

	return MxComponents{Components: result}
}
