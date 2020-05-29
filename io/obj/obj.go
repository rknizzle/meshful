package obj

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/rknizzle/meshful"
	"io"
	"os"
	"path/filepath"
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
	// write the obj file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	bufWriter := bufio.NewWriter(file)
	mtlData, err := writeObj(mesh, bufWriter)
	if err != nil {
		return err
	}
	bufWriter.Flush()

	// write the mtl file
	filename = strings.TrimSuffix(filename, filepath.Ext(filename))
	mtlFile, err := os.Create(filename + ".mtl")
	if err != nil {
		return err
	}
	defer file.Close()
	bufWriter = bufio.NewWriter(mtlFile)
	err = writeMaterial(mtlData, bufWriter)

	return bufWriter.Flush()
}

// write the mesh data to an obj file
func writeObj(mesh *meshful.Mesh, w io.Writer) ([]string, error) {
	// a map used to not duplicate vertices written to the obj file
	vertexTracker := make(map[string]int)

	// list of vertex strings to be written to the obj file
	var vertexList []string

	// faceLists groups lists of faces by their color/material
	faceLists := make(map[string][]string)

	// loop through each triangle in the mesh
	for _, triangle := range mesh.Triangles {
		// format the color into an obj string
		colorStr := formatColor(triangle.Color)

		// check if that color or lack-of color has already been seen in another triangle
		_, exists := faceLists[colorStr]
		// if it has not been seen
		if !exists {
			// initialize a new face list for the new color
			faceLists[colorStr] = []string{}
		}

		// tracks the 3 vertex numbers that make up the face
		verticesInFace := [3]int{}

		// for each vertex in the triangle
		for i := 0; i < 3; i++ {
			// format the vertex to an obj style line
			// also use the string as the key in the vertexTracker map to avoid duplicate vertex lines
			vStr := formatVertex(triangle.Vertices[i])

			number, exists := vertexTracker[vStr]
			// if the vertex hasn't already been added as an obj string
			if !exists {
				// add it to the vertex list
				vertexList = append(vertexList, vStr)

				// store the vertex number in the vertexTracker map to be used by upcoming faces that share
				// this vertex
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
		faceLists[colorStr] = append(faceLists[colorStr], formatFace(verticesInFace))
	}

	// write all the lines to the obj file

	// write comment header
	header := "# meshful OBJ export (github.com/rknizzle/meshful)\n\n"
	_, err := w.Write([]byte(header))
	if err != nil {
		return nil, err
	}

	// write all the vertex lines to the file
	for _, v := range vertexList {
		_, err := w.Write([]byte(v))
		if err != nil {
			return nil, err
		}
	}

	// store the lines to write to the accompanying mtl file
	mtlData := []string{}

	// write each face grouped by color
	var counter int = 1
	for colorStr, list := range faceLists {
		// if there is no color specified for some faces just set the color to grey
		if colorStr == "" {
			colorStr = "Kd 0.3 0.3 0.3"
		}

		// for each color add the lines for the color in the mtl file
		material := fmt.Sprintf("mtl%d", counter)
		// newMaterial declares the new material in the file
		// and colorStr is the line that defines the materials RGB values
		newMaterial := fmt.Sprintf("newmtl %s", material)
		// append both lines to the list of lines to be written to the mtl file
		mtlData = append(mtlData, newMaterial+"\n")
		mtlData = append(mtlData, colorStr+"\n")
		mtlData = append(mtlData, "\n")

		// specify that the next faces will be using the newly created material
		_, err := w.Write([]byte(fmt.Sprintf("usemtl %s\n", material)))
		if err != nil {
			return nil, err
		}

		// write the face lines belonging to this color/material
		for _, f := range list {
			_, err := w.Write([]byte(f))
			if err != nil {
				return nil, err
			}
		}
		counter++
	}

	// return the color/material info to write to an accompanying mtl file
	return mtlData, nil
}

// format mesh data into obj lines
func formatVertex(vertex meshful.Vec3) string {
	return fmt.Sprintf("v %f %f %f\n", vertex.X, vertex.Y, vertex.Z)
}

func formatFace(verticesInFace [3]int) string {
	return fmt.Sprintf("f %d %d %d\n", verticesInFace[0], verticesInFace[1], verticesInFace[2])
}

func formatColor(color *meshful.Color) string {
	if color == nil {
		return ""
	}
	return fmt.Sprintf("Kd %f %f %f\n", color.Red, color.Green, color.Blue)
}

// write the material/color contents of the mesh to the mtl file
func writeMaterial(mtlData []string, w io.Writer) error {
	// write comment header
	header := "# meshful mtl export (github.com/rknizzle/meshful)\n\n"
	_, err := w.Write([]byte(header))
	if err != nil {
		return err
	}

	// write material data
	for _, x := range mtlData {
		_, err := w.Write([]byte(x))
		if err != nil {
			return err
		}
	}
	return nil
}
