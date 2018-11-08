package algorithm

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)
var numCPU = 8

type TetrahedronJson struct {
	Id				int `json:"id"`
	Weight			float64 `json:"weight"`
	TissueId		int	`json:"tissueId"`
	VertexIds		[4]int `json:"vertexIds"`
	NeighborIds		[4]int `json:"neighborIds"`
}

func newGraph()(*Graph){
	vertices3d := readVerticesFromFile()
	tetrahedronJsons := readTetrahedronsFromFile()
	tetrahedrons := make([]*Tetrahedron, len(tetrahedronJsons))
	var plane []*Plane
	graph := &Graph{vertices3d, tetrahedrons, plane}
	graph.connect(tetrahedronJsons)
	return graph
}

func (graph *Graph) connect(tetJsons []*TetrahedronJson){
	for i := 0; i < len(tetJsons); i++ {
		var tetVertices [4]*Vertex3d
		var tetNeighbors [4]*Tetrahedron
		planeLastChecked := make([]int, numCPU)
		for cpu := 0; cpu < numCPU; cpu++ {
			planeLastChecked[cpu] = -1
		}
		graph.Tetrahedrons[i] = &Tetrahedron{i, tetJsons[i].Weight, tetJsons[i].TissueId, 
			tetVertices, tetNeighbors, planeLastChecked}
	}

	for tetIdx, tetJson := range tetJsons {
		tet := graph.Tetrahedrons[tetIdx]

		for idx, vertexId := range tetJson.VertexIds {
			vertex := graph.Vertices3d[vertexId]
			tet.Vertices[idx] = vertex
			vertex.Tetrahedrons = append(vertex.Tetrahedrons, tet)
		}

		for idx, neighborId := range tetJson.NeighborIds {
			tet.Neighbors[idx] = graph.Tetrahedrons[neighborId]
		}
	}
}

func readVerticesFromFile()([]*Vertex3d){
	var vertices3d []*Vertex3d

	err := json.Unmarshal(openJsonFile("./mesh/vertices3d.json"), &vertices3d)

	if err != nil {
		log.Fatal(err)
	}

	return vertices3d
}

func readTetrahedronsFromFile()([]*TetrahedronJson){
	var tetrahedrons []*TetrahedronJson

	err := json.Unmarshal(openJsonFile("./mesh/tetrahedrons.json"), &tetrahedrons)

	if err != nil {
		log.Fatal(err)
	}

	return tetrahedrons
}

func openJsonFile(fileName string)([]byte){
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