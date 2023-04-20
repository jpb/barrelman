package data

// Hunk describes a section of source code
type Hunk struct {
	Filepath  string
	StartLine int
	EndLine   int
}
