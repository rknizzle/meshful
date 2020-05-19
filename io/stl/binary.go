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

// Write solid in binary STL into an io.Writer.
// Does not check whether len(mesh.Triangles) fits into uint32.
func writeSolidBinary(w io.Writer, mesh *meshful.Mesh) error {
	headerBuf := make([]byte, 84)
	// write generic header
	copy(headerBuf, []byte("Exported by meshful"))

	// Write triangle count
	binary.LittleEndian.PutUint32(headerBuf[80:84], uint32(len(mesh.Triangles)))
	_, errHeader := w.Write(headerBuf)
	if errHeader != nil {
		return errHeader
	}

	// Write each triangle
	for _, t := range mesh.Triangles {
		tErr := writeTriangleBinary(w, &t)
		if tErr != nil {
			return tErr
		}
	}

	return nil
}

func writeTriangleBinary(w io.Writer, t *meshful.Triangle) error {
	buf := make([]byte, 50)
	offset := 0
	encodePoint(buf, &offset, &t.Normal)
	encodePoint(buf, &offset, &t.Vertices[0])
	encodePoint(buf, &offset, &t.Vertices[1])
	encodePoint(buf, &offset, &t.Vertices[2])
	// NOTE: Just not writing any attributes for now
	//encodeUint16(buf, &offset, t.Attributes)
	_, err := w.Write(buf)
	return err
}

func encodePoint(buf []byte, offset *int, pt *meshful.Vec3) {
	encodeFloat32(buf, offset, pt.X)
	encodeFloat32(buf, offset, pt.Y)
	encodeFloat32(buf, offset, pt.Z)
}

func encodeFloat32(buf []byte, offset *int, f float32) {
	u32 := math.Float32bits(f)
	binary.LittleEndian.PutUint32(buf[*offset:(*offset)+4], u32)
	(*offset) += 4
}

func encodeUint16(buf []byte, offset *int, u uint16) {
	binary.LittleEndian.PutUint16(buf[*offset:(*offset)+2], u)
	(*offset) += 2
}
