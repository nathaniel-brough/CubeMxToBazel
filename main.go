package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	internal "github.com/silvergasp/CubeMxToBazel/internal"
)

func main() {

	// TODO:Fix command line args
	const (
		ProjectFileUsage = "Path to the project file, e.g. project.gpdsc"
	)
	// Project File input
	defaultProject := findProjectFile()
	projectFile := *flag.String("project_file", defaultProject, ProjectFileUsage)
	// End TODO
	if projectFile == "" {
		log.Fatal("No project file found in working directory and none specified")
	}
	gpdscMxFile, err := ioutil.ReadFile(projectFile)
	if err != nil {
		log.Fatalf("Error opening file: %s \n Error:%s", projectFile, err)
	}
	project := internal.ProjectInit(gpdscMxFile)

	ccLibRules := internal.MxProjectToCcLibraryRules(project)
	ccBinRules := internal.MxProjectToCcBinaryRule(project)

	ccLibRulesStr := ""
	for _, libRule := range ccLibRules {
		ccLibRulesStr += libRule.String()
	}
	ccBinRuleStr := ccBinRules.String()
	BUILDString := ccLibRulesStr + ccBinRuleStr

	const (
		BUILD = "BUILD"
	)
	ioutil.WriteFile(BUILD, []byte(string(BUILDString)), 0664)
}

func findProjectFile() string {
	const (
		globPattern = "*.gpdsc"
	)
	projectFiles, err := filepath.Glob(globPattern)
	if err != nil {
		log.Fatal("Error finding project file:", err)
	}
	if len(projectFiles) == 1 {
		return projectFiles[0]
	} else if len(projectFiles) > 1 {
		log.Fatal("Multiple project files, only one valid .gpdsc file can be in the project root")
	}
	return ""
}
