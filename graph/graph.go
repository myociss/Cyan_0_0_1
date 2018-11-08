package graph

import (
	"../utils"
)

type Graph struct {
	Vertices3d		[]*Vertex3d
	Tetrahedrons	[]*Tetrahedron
}

func NewGraph(numCPU int)(*Graph){
	ucg := utils.ImportMesh()
	return connectedGraph(ucg, numCPU)
}

func connectedGraph(ucg *utils.UnconnectedGraph, numCPU int)(*Graph){
	vertices := getVertices(ucg)
	tetrahedrons := getTets(ucg, numCPU)

	for idx, v := range vertices {
		vertexFromFile := ucg.Vertices[idx]
		for tetIdx, tetId := range vertexFromFile.TetIds {
			v.Tetrahedrons[tetIdx] = tetrahedrons[tetId]
		}
	}

	for idx, tet := range tetrahedrons {
		tetFromFile := ucg.Tetrahedrons[idx]
		for vertexIdx, vertexId := range tetFromFile.VertexIds {
			tet.Vertices[vertexIdx] = vertices[vertexId]
		}

		for neighborIdx, neighborId := range tetFromFile.NeighborIds {
			tet.Neighbors[neighborIdx] = tetrahedrons[neighborId]
		}
	}

	return &Graph{vertices, tetrahedrons}
}

func getVertices(ucg *utils.UnconnectedGraph)([]*Vertex3d){
	vertices := make([]*Vertex3d, len(ucg.Vertices))

	for idx, v := range ucg.Vertices {
		vertices[idx] = &Vertex3d{idx, v.Vec, make([]*Tetrahedron, len(v.TetIds))}
	}

	return vertices
}

func getTets(ucg *utils.UnconnectedGraph, numCPU int)([]*Tetrahedron){
	tetrahedrons := make([]*Tetrahedron, len(ucg.Tetrahedrons))
	
	for i := 0; i < len(ucg.Tetrahedrons); i++ {
		ucgTet := ucg.Tetrahedrons[i]
		var tetVertices [4]*Vertex3d
		var tetNeighbors [4]*Tetrahedron
		planeLastChecked := make([]int, numCPU)
		for cpu := 0; cpu < numCPU; cpu++ {
			planeLastChecked[cpu] = -1
		}
		tetrahedrons[i] = &Tetrahedron{i, ucgTet.Weight, ucgTet.TissueId, 
			tetVertices, tetNeighbors, planeLastChecked}
	}

	return tetrahedrons
}