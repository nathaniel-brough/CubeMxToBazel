package cubemxtobazelinternal

// FileFilter can be used to recursively filter a slice of files
type FileFilter interface {
	AssemblyFiles() FileFilter
	SourceFiles() FileFilter
	HeaderFiles() FileFilter
	Condition(string) FileFilter
	Files() MxFiles
}

// MxFiles slice of CubeMx project and driver files
type MxFiles []MxFile
