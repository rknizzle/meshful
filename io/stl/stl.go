/*
The MIT License (MIT)

Copyright (c) 2014 Hagen Schendel

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
*/

package stl

import (
	"bufio"
	"bytes"
	"errors"
	"github.com/rknizzle/meshful"
	"io"
	"os"
)

// ErrIncompleteBinaryHeader is used when reading binary STL files with incomplete header.
var ErrIncompleteBinaryHeader = errors.New("Incomplete STL binary header, 84 bytes expected")

// ErrUnexpectedEOF is used by ReadFile and ReadAll to signify an incomplete file.
var ErrUnexpectedEOF = errors.New("Unexpected end of file")

var asciiStart = []byte("solid ")

// ReadFile reads the contents of a file into a new Mesh object. The file
// can be either in STL ASCII format, beginning with "solid ", or in
// STL binary format, beginning with a 84 byte header. Shorthand for os.Open and ReadAll
func ReadFile(filename string) (mesh *meshful.Mesh, err error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		err = openErr
		return
	}
	defer file.Close()

	return ReadAll(bufio.NewReader(file))
}

// ReadAll reads the contents of a file into a new Mesh object. The file
// can be either in STL ASCII format, beginning with "solid ", or in
// STL binary format, beginning with a 84 byte header. Because of this,
// the file pointer has to be at the beginning of the file.
func ReadAll(r io.Reader) (mesh *meshful.Mesh, err error) {
	isASCII, first6, isASCIIErr := isASCIIFile(r)
	if isASCIIErr != nil {
		err = isASCIIErr
		return
	}

	if isASCII {
		return nil, errors.New("ascii not yet implemented")
	} else {
		mesh, err = readAllBinary(r, first6)
	}

	return
}

// isASCIIFile detects if the file is in STL ASCII format or if it is binary otherwise.
// It will consume 6 bytes and return them.
func isASCIIFile(r io.Reader) (isASCII bool, first6 []byte, err error) {
	first6 = make([]byte, 6) // "solid "
	isASCII = false
	n, readErr := r.Read(first6)
	if n != 6 || readErr == io.EOF {
		err = ErrUnexpectedEOF
		return
	} else if readErr != nil {
		err = readErr
		return
	}

	if bytes.Equal(first6, asciiStart) {
		isASCII = true
	}

	return
}

// Extracts an ASCII string from a byte slice. Reads all characters
// from the beginning until a \0 or a non-ASCII character is found.
func extractASCIIString(byteData []byte) string {
	i := 0
	for i < len(byteData) && byteData[i] < byte(128) && byteData[i] != byte(0) {
		i++
	}
	return string(byteData[0:i])
}
