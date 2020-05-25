package obj

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/rknizzle/meshful"
	"io"
	"os"
	"strconv"
	"strings"
)

// Readfile reads the contents of a Wavefront OBJ file into a new Mesh object
func ReadFile(filename string) (mesh *meshful.Mesh, err error) {
	file, openErr := os.Open(filename)
	if openErr != nil {
		err = openErr
		return
	}
	defer file.Close()

	return readAll(bufio.NewReader(file))
}

func readAll(r io.Reader) (mesh *meshful.Mesh, err error) {
	scanner := bufio.NewScanner(r)

	// keep a list of all the vertices and faces specified in the file
	vertices := []meshful.Vec3{}
	faces := []meshful.Triangle{}

	// loop through each line of the file
	for scanner.Scan() {
		line := scanner.Text()

		// skip blank lines
		if line == "" {
			continue
		}

		// the first word of each line should be a token specifying the data type of that line
		tokens := strings.Fields(line)
		firstToken := tokens[0]

		if firstToken == "#" {
			// Its just a comment, continue to next line
			continue
		}
		if firstToken == "v" {
			// new vertex -- get the coordinates and add it to the list of vertices
			v, err := parseVertex(tokens)
			if err != nil {
				return nil, err
			}
			vertices = append(vertices, v)
		}
		if firstToken == "f" {
			// new face -- construct the face using the list of vertices
			f, err := parseFace(tokens, vertices)
			if err != nil {
				return nil, err
			}
			faces = append(faces, f)
		}
	}

	return &meshful.Mesh{faces}, nil
}

// parse the line of the OBJ file into a Vec3 data structure
// example value:
// v 0.000000 10.000000 0.000000
func parseVertex(tokens []string) (meshful.Vec3, error) {
	if len(tokens) != 4 {
		return meshful.Vec3{}, errors.New("Incorrect number of tokens in the vertex line")
	}

	x, err := strconv.ParseFloat(tokens[1], 32)
	if err != nil {
		return meshful.Vec3{}, err
	}

	y, err := strconv.ParseFloat(tokens[2], 32)
	if err != nil {
		return meshful.Vec3{}, err
	}

	z, err := strconv.ParseFloat(tokens[3], 32)
	if err != nil {
		return meshful.Vec3{}, err
	}

	return meshful.Vec3{float32(x), float32(y), float32(z)}, nil
}

// parse the line of the OBJ file into a Triangle data structure
// example values:
// f 1/1/1 2/2/2 3/3/3
// f 1 2 3
func parseFace(tokens []string, vertices []meshful.Vec3) (meshful.Triangle, error) {
	if len(tokens) != 4 {
		return meshful.Triangle{}, errors.New("Incorrect number of tokens in the face line")
	}

	faceVerts := [3]meshful.Vec3{}
	// get the data for each vertex
	for i := 1; i <= 3; i++ {
		// example vertex data values: 1/1/1 or 1
		// if the vertex data is seperated by '/', split it and get the first value (the number in the vertices list)
		vertexData := strings.Split(tokens[i], "/")
		// convert vertex number value from string -> int
		vertexNumber, err := strconv.Atoi(vertexData[0])
		if err != nil {
			return meshful.Triangle{}, err
		}

		vertexIndex := vertexNumber - 1
		faceVerts[i-1] = vertices[vertexIndex]
	}

	return meshful.Triangle{Vertices: faceVerts}, nil
}

func WriteFile(filename string, mesh *meshful.Mesh) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	err = writeAll(mesh, bufWriter)
	if err != nil {
		return err
	}
	return bufWriter.Flush()
}

func writeAll(mesh *meshful.Mesh, w io.Writer) error {
	// a map used to not duplicate vertices written to the obj file
	vertexTracker := make(map[string]int)

	// list of strings to be written to the obj file
	var vertexList []string
	var faceList []string

	// loop through each triangle in the mesh
	for _, triangle := range mesh.Triangles {

		// tracks the 3 vertex numbers that make up the face
		verticesInFace := [3]int{}

		// for each vertex in the triangle
		for i := 0; i < 3; i++ {
			// format the vertex to the obj style line
			// also use the string as the key in the vertexTracker map
			vStr := formatVertex(triangle.Vertices[i])

			number, exists := vertexTracker[vStr]
			// if the vertex hasn't already been added as an obj string
			if !exists {
				// add it to the vertex list
				vertexList = append(vertexList, vStr)
				// store the vertex number in the vertexTracker map
				vertexNumberInList := len(vertexList)
				vertexTracker[vStr] = vertexNumberInList
				// add the vertex number to the vertices that make up this face
				verticesInFace[i] = vertexNumberInList
			} else {
				// else the vertex already exists, so no need to add a new vertex line
				// just set that this vertex number is part of the current face
				verticesInFace[i] = number
			}
		}

		// format the face into an obj line
		faceList = append(faceList, formatFace(verticesInFace))
	}

	// write all the vertex lines to the file
	for _, v := range vertexList {
		_, err := w.Write([]byte(v))
		if err != nil {
			return err
		}
	}

	// write all the face lines to the file
	for _, f := range faceList {
		_, err := w.Write([]byte(f))
		if err != nil {
			return err
		}
	}

	return nil
}

func formatVertex(vertex meshful.Vec3) string {
	return fmt.Sprintf("v %f %f %f\n", vertex.X, vertex.Y, vertex.Z)
}

func formatFace(verticesInFace [3]int) string {
	return fmt.Sprintf("f %d %d %d\n", verticesInFace[0], verticesInFace[1], verticesInFace[2])
}
