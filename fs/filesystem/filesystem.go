// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

package filesystem

import (
	"io"
	"os"
	"path/filepath"

	"github.com/praelatus/praelatus/config"
	"github.com/praelatus/praelatus/fs"
)

// FS implements fs.FS for a local filesystem.
type FS struct {
	baseDir string
}

// New will open a filesystem.FS at the default path
func New() FS {
	return FS{config.DataDir()}
}

// Init will create the baseDir
func (f FS) Init() error {
	return os.MkdirAll(f.baseDir, os.ModePerm)
}

// Get will retrieve the file by joining the given path with FS.baseDir
func (f FS) Get(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(f.baseDir, path))
}

// Save will generate a unique file name and save the file to the baseDir
func (f FS) Save(file *os.File) (string, error) {
	fn, err := fs.GenUniqueFileName(file)
	if err != nil {
		return "", err
	}

	newFile, err := os.Create(filepath.Join(f.baseDir, fn))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(file, newFile)
	return fn, err
}
