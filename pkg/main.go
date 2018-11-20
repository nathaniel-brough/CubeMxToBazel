package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"

	internal "github.com/silvergasp/CubeMxToBazel/internal"
)

var projectFile string

func init() {
	const (
		ProjectFileUsage = "Path to the project file, e.g. project.gpdsc"
	)
	// Project File input
	flag.StringVar(&projectFile, "project_file", findProjectFile(), ProjectFileUsage)
	flag.StringVar(&projectFile, "p", findProjectFile(), ProjectFileUsage+" (shorthand)")
}

func main() {
	if projectFile == "" {
		log.Fatal("No project file found in working directory and none specified")
	}
	gpdscMxFile, err := ioutil.ReadFile(projectFile)
	if err != nil {
		log.Fatalf("Error opening file: %s \n Error:%s", projectFile, err)
	}
	project := internal.ProjectInit(gpdscMxFile)
	components := project.Components()

	ccLibRules := ""
	for _, component := range components {
		ccLibRules = ccLibRules + internal.MxComponentToCcLibraryRule(component).String()
	}
	const (
		BUILD = "BUILD"
	)
	ioutil.WriteFile(BUILD, []byte(string(ccLibRules)), 0664)
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
