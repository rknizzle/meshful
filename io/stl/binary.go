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
	"encoding/binary"
	"fmt"
	"github.com/rknizzle/meshful"
	"io"
	"math"
)

func readAllBinary(r io.Reader, first6 []byte) (mesh *meshful.Mesh, err error) {
	header := make([]byte, 84)
	copy(header, first6)
	n, readErr := r.Read(header[6:])
	if readErr == io.EOF && n != 84 {
		err = ErrIncompleteBinaryHeader
		return
	} else if readErr != nil {
		err = readErr
		return
	}

	var meshData meshful.Mesh
	triangleCount := binary.LittleEndian.Uint32(header[80:84])
	meshData.Triangles = make([]meshful.Triangle, triangleCount)

	for i := range meshData.Triangles {
		readErr = readTriangleBinary(r, &meshData.Triangles[i])
		if readErr != nil {
			err = fmt.Errorf("While reading triangle no. %d at byte %d: %s", i, 84+i*50, readErr.Error())
			return
		}
	}

	mesh = &meshData
	return
}

func readTriangleBinary(r io.Reader, t *meshful.Triangle) error {
	tbuf := make([]byte, 50)
	n := 0
	for n < 50 {
		l, readErr := r.Read(tbuf[n:])
		n += l
		if readErr != nil {
			if readErr == io.EOF {
				return ErrUnexpectedEOF
			}
			return readErr
		}
	}

	offset := 0
	readBinaryPoint(tbuf, &offset, &(t.Normal))
	readBinaryPoint(tbuf, &offset, &(t.Vertices[0]))
	readBinaryPoint(tbuf, &offset, &(t.Vertices[1]))
	readBinaryPoint(tbuf, &offset, &(t.Vertices[2]))
	return nil
}

func readBinaryPoint(buf []byte, offset *int, p *meshful.Vec3) {
	p.X = readBinaryFloat32(buf, offset)
	p.Y = readBinaryFloat32(buf, offset)
	p.Z = readBinaryFloat32(buf, offset)
}

func readBinaryFloat32(buf []byte, offset *int) float32 {
	v := binary.LittleEndian.Uint32(buf[*offset : (*offset)+4])
	(*offset) += 4
	return math.Float32frombits(v)
}

func readBinaryUint16(buf []byte, offset *int) uint16 {
	v := binary.LittleEndian.Uint16(buf[*offset : (*offset)+2])
	(*offset) += 2
	return v
}
