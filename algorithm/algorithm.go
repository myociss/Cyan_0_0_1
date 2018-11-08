package algorithm

import (
	"../graph"
	"../plane3d"
	//"log"
	"math"
)

var g *graph.Graph

func Start(graph *graph.Graph){
	g = graph
}

func GetZRange()(float64, float64){
	minZ, maxZ := math.Inf(1), math.Inf(-1)

	for _, vertex := range g.Vertices3d {
		z := vertex.Vec[2]
		if z < minZ {
			minZ = z
		}
		if z > maxZ {
			maxZ = z
		}
	}
	return minZ, maxZ
}

func GetSelectionSlice(zVal float64)([]*plane3d.VertexGroup){
	var planeOriginTets []*graph.Tetrahedron
	target := [3]float64{0.0, 0.0, zVal}

	for _, tet := range g.Tetrahedrons {
		if tet.ContainsPoint(target) && !tet.IsFlat() {
			planeOriginTets = append(planeOriginTets, tet)
		}
	}
	planeGen := plane3d.NewPlaneGenerator(g, 0, 0.0, 0.0, [3]float64{0.0, 0.0, zVal}, planeOriginTets)
	//log.Printf("%+v\n", planeGen)
	planeGen.Generate("sequential", 1)
	return planeGen.FoundVertices
}

func GetPlane(target [3]float64, tetId int, alphaId int, thetaId int, epsilon int)(*plane3d.PlaneGen){
		alpha := (float64(alphaId) * math.Pi) / float64(epsilon)
		theta := (float64(thetaId) * math.Pi) / float64(epsilon)
		planeGen := plane3d.NewPlaneGenerator(g, 0, alpha, theta, target,
			[]*graph.Tetrahedron{g.Tetrahedrons[tetId]})
		planeGen.Generate("sequential", 1)
		//return planeGen.FoundVertices, planeGen.Plane
		return planeGen
}