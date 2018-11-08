package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Vertex struct {
	Id 				int `json:"id"`
	Vec 			[3]float64 `json:"vec"`
	TetIds 			[]int `json:"tetrahedrons"`
}

type Tetrahedron struct {
	Id				int `json:"id"`
	Weight			float64 `json:"weight"`
	TissueId		int	`json:"tissueId"`
	VertexIds		[4]int `json:"vertexIds"`
	NeighborIds		[4]int `json:"neighborIds"`
}

type UnconnectedGraph struct {
	Vertices		[]*Vertex
	Tetrahedrons 	[]*Tetrahedron
}

func ImportMesh()(*UnconnectedGraph){
	verticesFromFile := readVerticesFromFile()
	tetsFromFile := readTetrahedronsFromFile()
	return &UnconnectedGraph{verticesFromFile,  tetsFromFile}
}

func readVerticesFromFile()([]*Vertex){
	var vertices []*Vertex

	err := json.Unmarshal(OpenJsonFile("./mesh/vertices3d.json"), &vertices)

	if err != nil {
		log.Fatal(err)
	}

	return vertices
}

func readTetrahedronsFromFile()([]*Tetrahedron){
	var tetrahedrons []*Tetrahedron

	err := json.Unmarshal(OpenJsonFile("./mesh/tetrahedrons.json"), &tetrahedrons)

	if err != nil {
		log.Fatal(err)
	}

	return tetrahedrons
}

func OpenJsonFile(fileName string)([]byte){
	jsonFile, err := os.Open(fileName)

	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	fileBytes, err := ioutil.ReadAll(jsonFile)

	if err != nil {
		log.Fatal(err)
	}

	return fileBytes
}