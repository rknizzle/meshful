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
			f, err := parseFace(tokens)
			if err != nil {
				return nil, err
			}
			faces = append(faces, f)
		}

		fmt.Println(firstToken)
	}

	return &meshful.Mesh{}, nil
}

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

func parseFace(tokens []string) (meshful.Triangle, error) {
	if len(tokens) != 4 {
		return meshful.Triangle{}, errors.New("Incorrect number of tokens in the face line")
	}
	return meshful.Triangle{}, nil
}
