package algorithm

import (
	//"log"
	"sort"
	"math"
)

type Plane2d struct {
	Groups 			[]*VertexGroup
	Points			[]*Vertex2d
	LineComponents	[][3]float64
	TargetGroup		*VertexGroup
	Distances		[]float64
}

type VertexGroup struct {
	Vertices	[]*Vertex2d
	Weight		float64
}

func (plane3d *Plane) newPlane2d(groups []*VertexGroup3d)(*Plane2d){
	count := 0
	//shapes := make([]*Shape2d, len(groups))

	for _, vertexGroup := range groups {
		count += len(vertexGroup.Vertices)
		//shapes[idx]
	}

	points := make([]*Vertex2d, count)
	groups2d := make([]*VertexGroup, len(groups))
	//shapes := make([]*Shape2d, len(groups))
	count = 0

	for idx0, vertexGroup := range groups {
		vertices2d := make([]*Vertex2d, len(vertexGroup.Vertices))
		for idx1, vertex := range vertexGroup.Vertices {
			newVec := plane3d.get2DCoords(vertex)
			angle := math.Atan2(newVec[1], newVec[0])
			newVertex := &Vertex2d{newVec, angle, -1}
			vertices2d[idx1] = newVertex
			points[count] = newVertex
			count++
		}
		groups2d[idx0] = &VertexGroup{vertices2d, vertexGroup.Weight}
	}
	sort.Sort(VertexSorter(points))

	var distances []float64
	var lineComponents [][3]float64

	plane := &Plane2d{groups2d, points, lineComponents, groups2d[0], distances}
	//plane.sortByAngle()
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
		if shape.MinVertex.AngleId == shape.StartVertex.AngleId {
			shape.calcInternalDistances(plane, shape.MaxVertex.AngleId)
		} else {
			shape.calcInternalDistances(plane, shape.MaxVertex.AngleId + 1)
			shape.calcInternalDistances(plane, shape.StopVertex.AngleId)
		}

		//shape.calcInternalDistances(plane)
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
		for _, vertex := range vertexGroup.Vertices {
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
			edge := line([2]float64{0.0, 0.0}, vertex.Vec)
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

func line(vec0 [2]float64, vec1 [2]float64)([3]float64) {
	A := vec0[1] - vec1[1]
	B := vec1[0] - vec0[0]
	C := vec0[0] * vec1[1] - vec1[0] * vec0[1]
	return [3]float64{A, B, -C}
}

func findIntersect(line0 [3]float64, line1 [3]float64)([2]float64){

	D  := line0[0] * line1[1] - line0[1] * line1[0]
    Dx := line0[2] * line1[1] - line0[1] * line1[2]
	Dy := line0[0] * line1[2] - line0[2] * line1[0]
	
	return [2]float64{Dx / D, Dy / D}
}

/*func findIntersect(L1 [3]float64, L2 [3]float64)([2]float64){
	//check for line through y axis

	//L1 := line(e0[0], e0[1])
	//L2 := line(e1[0], e1[1])
	D  := L1[0] * L2[1] - L1[1] * L2[0]
    Dx := L1[2] * L2[1] - L1[1] * L2[2]
	Dy := L1[0] * L2[2] - L1[2] * L2[0]
	
	return [2]float64{Dx / D, Dy / D}

}*/

type VertexSorter []*Vertex2d

func (a VertexSorter) Len() int           { return len(a) }
func (a VertexSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VertexSorter) Less(i, j int) bool { return a[i].Angle < a[j].Angle }

//possible optimization:
	//create list of unique vertices
	//list will be sorted by x value
	//if x value exists check all x values for y value
	//if no y value insert vertex
	//for each shape:
		//for each 2d vertex:
			//attempt to add to unique list
			//get vertex in unique list
			//replace pointer in shapes.vertices with that pointer
			//add this shape to that vertex