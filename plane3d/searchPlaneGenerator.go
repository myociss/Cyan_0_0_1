package plane3d

import (
	"../graph"
	//"log"
	//"math"
)

/* finds a group of poins that lie on the faces of a tetrahedron within
 * this plane.
 */

type VertexGroup struct {
	//Points		[][2]float64
	FoundVertices	[]*FoundVertex2d
	Weight			float64
	TissueId		int
	TetId			int
}

type PlaneGen struct {
	Id 				int
	Graph 			*graph.Graph
	Plane			*Plane3d
	TargetTets		[]*graph.Tetrahedron
	Stack			*Stack
	FoundVertices	[]*VertexGroup
}

//var allVertices []*VertexGroup

func NewPlaneGenerator(graph *graph.Graph, id int, alpha float64, theta float64, target [3]float64, targetTets []*graph.Tetrahedron)(*PlaneGen){
	plane := newPlane3d(alpha, theta, target)
	var foundVertices []*VertexGroup
	return &PlaneGen{id, graph, plane, targetTets, nil, foundVertices}
}


func (planeGen *PlaneGen) Generate(order string, threadId int){
	planeGen.initStack()
	if (order == "concurrent"){
		planeGen.generateConcurrent(threadId)
	} else {
		planeGen.generateSequential()
	}
	
}

func (planeGen *PlaneGen) generateSequential(){

	tetsChecked := make([]bool, len(planeGen.Graph.Tetrahedrons))

	for i := 0; i < len(tetsChecked); i++ {
		tetsChecked[i] = false
	}

	//var allVertices []*VertexGroup

	for planeGen.Stack.Count > 0 {
		//log.Println("here")
		tet := planeGen.Stack.pop().Tet

		if tet != nil {
			if !tetsChecked[tet.Id] {
				planeGen.formGroup(tet)
				tetsChecked[tet.Id] = true
			}
		}
	}
	//return allVertices
}

func (planeGen *PlaneGen) generateConcurrent(threadId int){
	
	//var allVertices []*VertexGroup

	for planeGen.Stack.Count > 0 {
		tet := planeGen.Stack.pop().Tet

		if tet != nil {

			if tet.PlaneLastChecked[threadId] != planeGen.Id {
				planeGen.formGroup(tet)
				tet.PlaneLastChecked[threadId] = planeGen.Id
			}
		}
	}
	//return allVertices
}

func (planeGen *PlaneGen) formGroup(tet *graph.Tetrahedron){
	//log.Println("HERE")
	group := &VertexGroup{planeGen.Plane.formVertexGroup(planeGen.Stack, tet), tet.Weight, 
		tet.TissueId, tet.Id}
	if len(group.FoundVertices) > 2 {
		//log.Println(len(allVertices))
		planeGen.FoundVertices = append(planeGen.FoundVertices, group)		
		for _, neighbor := range tet.Neighbors {
			planeGen.Stack.push(&Node{neighbor})
		}
	}
}

func (planeGen *PlaneGen) initStack(){
	planeGen.Stack = newStack()

	for _, tet := range planeGen.TargetTets {
		planeGen.Stack.push(&Node{tet})
	}
}