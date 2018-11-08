package algorithm

import (
	"../plane3d"
	"sort"
	//"log"
	"../utils"
)

type Plane2d struct {
	Groups 			[]*plane3d.VertexGroup
	Points			[]*plane3d.FoundVertex2d
	LineComponents	[][3]float64
	TargetGroup		*plane3d.VertexGroup
	Distances		[]float64
}

func newPlane2d(groups []*plane3d.VertexGroup)(*Plane2d){
	count := 0

	for _, vertexGroup := range groups {
		count += len(vertexGroup.FoundVertices)
	}

	points := make([]*plane3d.FoundVertex2d, count)
	//groups2d := make([]*plane3d.VertexGroup, len(groups))

	var distances []float64
	var lineComponents [][3]float64

	plane := &Plane2d{groups, points, lineComponents, groups[0], distances}
	plane.sortByAngle()
	plane.getUniqueLines()
	return plane
}

func (plane *Plane2d) findShortestPaths() {
	paths := make([]float64, plane.Points[len(plane.Points) - 1].AngleId + 1)

	for idx, _ := range paths {
		paths[idx] = 0.0
	}

	for i := 1; i < len(plane.Groups); i++ {
		shape := newShape(plane.Groups[i])
		

		shape.calcInternalDistances(plane)
		//if shape.MinAngle.AngleId < shape.MaxAngleId {
			//distances := shape.calcDistances(plane)
			//idx of 0 is shape min angle id + 1
			//for idx, dist := range distances {
			//	paths[idx + shape.MinAngleId + 1] +=  dist
			//}
		//}
	}
}

func (plane *Plane2d) sortByAngle(){
	count := 0

	for _, vertexGroup := range plane.Groups {
		for _, vertex := range vertexGroup.FoundVertices {
			plane.Points[count] = vertex
			count++
		}
	}
	sort.Sort(VertexSorter(plane.Points))
}


func (plane *Plane2d) getUniqueLines(){
	count := 0
	anglePrev := plane.Points[0].Angle

	for idx, vertex := range plane.Points {
		if vertex.Angle != anglePrev || idx == 0 {
			edge := utils.Line([2]float64{0.0, 0.0}, vertex.Vec)
			plane.LineComponents = append(plane.LineComponents, edge)
			count++
		}

		vertex.AngleId = count - 1
		anglePrev = vertex.Angle
	}

	plane.Distances = make([]float64, len(plane.LineComponents))
	for i := 0; i < len(plane.Distances); i++ {
		plane.Distances[i] = 0
	}
}