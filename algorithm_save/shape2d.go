package algorithm

import (
	"log"
	"math"
)

type Shape2d struct {
	Points 		[]*Vertex2d
	Weight		float64
	HighEdge	[2]int
	LowEdge		[2]int
	StartVertex	*Vertex2d
	StopVertex	*Vertex2d
	MinVertex	*Vertex2d
	MaxVertex	*Vertex2d
}

func newShape(group *VertexGroup)(*Shape2d){
	points := group.Vertices
	minAngle, maxAngle := points[0], points[0]
	minAngleIdx, maxAngleIdx := 0, 0
	var start, stop *Vertex2d
	var startIdx, stopIdx int

	for i := 0; i < len(points); i++ {
		if points[i].Angle < minAngle.Angle {
			minAngle = points[i]
			minAngleIdx = i
		}

		if points[i].Angle > maxAngle.Angle {
			maxAngle = points[i]
			maxAngleIdx = i
		}
	}

	if maxAngle.Angle - minAngle.Angle > math.Pi {
		log.Println("THIS HAPPENS")
		startIdx, stopIdx = findBoundingAngleIds(points)
		log.Println(startIdx)
		log.Println(stopIdx)
		start, stop = points[startIdx], points[stopIdx]
	} else {
		start, stop, startIdx, stopIdx = minAngle, maxAngle, minAngleIdx, maxAngleIdx
	}

	highEdge := [2]int{startIdx, (startIdx + 1) % len(points)}
	lowEdge := [2]int{startIdx - 1, startIdx}
	if lowEdge[0] == -1 {
		lowEdge[0] = len(points) - 1
	}

	shape := &Shape2d{points, group.Weight, highEdge, lowEdge, start, stop, minAngle, maxAngle}
	shape.checkColinearStart()
	return shape
}

func (shape *Shape2d) calcInternalDistances(plane *Plane2d, stopAngleId int){
	angleId := shape.Points[shape.LowEdge[1]].AngleId

	v0High, v1High := shape.Points[shape.HighEdge[0]], shape.Points[shape.HighEdge[1]]
	v0Low, v1Low := shape.Points[shape.LowEdge[0]], shape.Points[shape.LowEdge[1]]

	highEdgeAsLine := line(v0High.Vec, v1High.Vec)
	lowEdgeAsLine := line(v0Low.Vec, v1Low.Vec)

	for angleId < stopAngleId {
		log.Println(angleId)
		i0, i1 := findIntersect(plane.LineComponents[angleId], highEdgeAsLine),
			findIntersect(plane.LineComponents[angleId], lowEdgeAsLine)

		xDist, yDist := i1[0] - i0[0], i1[1] - i0[1]

		dist := shape.Weight * math.Sqrt(xDist * xDist + yDist * yDist)

		//distances[count] = dist
		//count++
		plane.Distances[angleId] += dist

		if v1High.AngleId == angleId {
			shape.shiftHighEdge()
			v0High, v1High = shape.Points[shape.HighEdge[0]], shape.Points[shape.HighEdge[1]]
			highEdgeAsLine = line(v0High.Vec, v1High.Vec)
		}

		if v0Low.AngleId == angleId {
			shape.shiftLowEdge()
			v0Low, v1Low = shape.Points[shape.LowEdge[0]], shape.Points[shape.LowEdge[1]]
			lowEdgeAsLine = line(v0Low.Vec, v1Low.Vec)
		}
		angleId++

	}
}

func findBoundingAngleIds(points []*Vertex2d)(int, int){
	maxAngleDiff := 0.0
	startIdx := 0
	endIdx := 0
	//find biggest difference in angles; the point after this is our start point
	for i, _ := range points {
		next := i + 1

		if next == len(points){
			next = 0
		}
		angleDiff := points[next].Angle - points[i].Angle
		if angleDiff > maxAngleDiff {
			maxAngleDiff = angleDiff
			startIdx = next
			endIdx = i
		}
	}

	return startIdx, endIdx
}

func (shape *Shape2d) checkColinearStart(){
	if shape.Points[shape.HighEdge[0]].AngleId == shape.Points[shape.HighEdge[1]].AngleId {
		shape.shiftHighEdge()
	}

	if shape.Points[shape.LowEdge[0]].AngleId == shape.Points[shape.LowEdge[1]].AngleId {
		shape.shiftLowEdge()
	}
}

func (shape *Shape2d) shiftLowEdge(){
	if shape.LowEdge[0] == 0 {
		shape.LowEdge[0] = len(shape.Points) - 1
		shape.LowEdge[1]--
	} else if shape.LowEdge[1] == 0 {
		shape.LowEdge[0]--
		shape.LowEdge[1] = len(shape.Points) - 1
	} else {
		shape.LowEdge[0]--
		shape.LowEdge[1]--
	}
}

func (shape *Shape2d) shiftHighEdge(){
	if shape.HighEdge[0] == len(shape.Points) - 1 {
		shape.HighEdge[0] = 0
		shape.HighEdge[1]++
	} else if shape.HighEdge[1] == len(shape.Points) - 1 {
		shape.HighEdge[0]++
		shape.HighEdge[1] = 0
	} else {
		shape.HighEdge[0]++
		shape.HighEdge[1]++
	}
}