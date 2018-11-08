package algorithm

import (
	"log"
	//"fmt"
	"math"
	"sync"
	//"../utils"
)

type Graph struct {
	Vertices3d		[]*Vertex3d
	Tetrahedrons	[]*Tetrahedron
	Planes	 		[]*Plane
}

var graph *Graph
var selectionSliceOriginTets map[float64][]*Tetrahedron
var allSlices map[float64][][][2]float64

func GetGraph()(*Graph){
	return graph
}

func GetGraphFromFiles(){
	graph = newGraph()
}

func GetSelectionSlice(zVal float64)([]*VertexGroup3d){
	var planeOriginTets []*Tetrahedron
	target := [3]float64{0.0, 0.0, zVal}

	for _, tet := range graph.Tetrahedrons {
		if tet.containsPoint(target) && !tet.isFlat() {
			planeOriginTets = append(planeOriginTets, tet)
		}
	}
	plane := newPlane(0, 0.0, 0.0, [3]float64{0.0, 0.0, zVal}, planeOriginTets)
	tetsLength := len(graph.Tetrahedrons)
	return plane.findCrossSection(0, GET_SELECTION_SLICES, tetsLength)
}

func GetPlane(target [3]float64, tetId int, alphaId int, thetaId int, epsilon int)([]*VertexGroup3d, *Plane){
	/*log.Println(target)
	log.Println(tetId)
	log.Println(alpha)
	log.Println(theta)*/
	alpha := (float64(alphaId) * math.Pi) / float64(epsilon)
	theta := (float64(thetaId) * math.Pi) / float64(epsilon)
	plane := newPlane(0, alpha, theta, target, []*Tetrahedron{graph.Tetrahedrons[tetId]})
	return plane.findCrossSection(0, GET_SELECTION_SLICES, len(graph.Tetrahedrons)), plane
}

func GetPlaneSpace(target [3]float64, tetId int, epsilon int)(){
	allPlanes := getAllPlanes(target, tetId, epsilon)

	var wg sync.WaitGroup

	sliceSize := (len(allPlanes) + numCPU - 1) / numCPU

	for i := 0; i < len(allPlanes); i += sliceSize {
		end := i + sliceSize
	
		if end > len(allPlanes) {
			end = len(allPlanes)
		}

		wg.Add(1)

		go func(planes []*Plane, threadId int){
			defer wg.Done()
			for _, plane := range planes {
				newPoints := plane.findCrossSection(threadId, GET_PLANE_SPACE, 0)
				log.Println(len(newPoints))
			}
		}(allPlanes[i:end], i / sliceSize)
	}
	wg.Wait()
	log.Println("done finding groups")
}

func getAllPlanes(target [3]float64, tetId int, epsilon int)([]*Plane){
	planes := make([]*Plane, epsilon * epsilon)

	for i := 0; i < epsilon; i++ {
		alpha := (float64(i) * math.Pi) / float64(epsilon)
		for j := 0; j < epsilon; j++ {
			theta := (float64(j) * math.Pi) / float64(epsilon)
			startTets := []*Tetrahedron{graph.Tetrahedrons[tetId]}
			plane := newPlane(i * epsilon + j, alpha, theta, target, startTets)
			planes[i * epsilon + j] = plane
		}
	}
	return planes
}
