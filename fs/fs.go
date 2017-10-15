// Copyright 2017 Mathew Robinson <mrobinson@praelatus.io>. All rights reserved.
// Use of this source code is governed by the AGPLv3 license that can be found in
// the LICENSE file.

// Package fs provides definitions for abstracting away filesystem interaction
// in Praelatus
package fs

import (
	"crypto/md5"
	cryp "crypto/rand"
	"fmt"
	"io"
	"math/rand"
	"os"
)

// FS is the interface all filesystem abstractions must implement.
type FS interface {
	Init() error

	Get(path string) (io.ReadCloser, error)
	Save(file *os.File) (string, error)
}

// GenUniqueFileName takes a file and hashes it's name with a salt to avoid
// collisions.
func GenUniqueFileName(file *os.File) (string, error) {
	saltString, err := newSalt(16)
	fn := fmt.Sprintf("%x", md5.Sum([]byte(file.Name()+saltString)))
	return fn, err
}

// newSalt generates a new salt string
func newSalt(length int) (string, error) {
	real := length
	if length == 0 {
		real = rand.Intn(100)
	}

	res := make([]byte, real)

	_, err := cryp.Read(res)
	return string(res), err
}
