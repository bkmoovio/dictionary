package dictionary

import "github.com/moov-io/base"

type File struct {
	// FilePath
	FilePath string
	// File Name being parsed
	FileName string
	// FileType being parsed
	FileType string

	errors base.ErrorList
}
