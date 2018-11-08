package algorithm

import (
	"../plane3d"
	//"log"
	"sort"
	"math"
	"../utils"
)

type Shape2d struct {
	Points 		[]*plane3d.FoundVertex2d
	Weight		float64
	LowEdge		[2]int
	HighEdge	[2]int
	EndAngleId	int
}

func newShape(group *plane3d.VertexGroup)(*Shape2d){
	points := group.FoundVertices
	lowVertex, highVertex := points[0], points[0]
	lowIdx, highIdx := 0, 0

	for i := 0; i < len(points); i++ {
		if points[i].Angle < lowVertex.Angle {
			lowVertex = points[i]
			lowIdx = i
		}

		if points[i].Angle > highVertex.Angle {
			highVertex = points[i]
			highIdx = i
		}
	}

	var lowEdge, highEdge [2]int

	shape := &Shape2d{points, group.Weight, lowEdge, highEdge, 0}

	shape.findStartingEdges(lowIdx, highIdx)
	return shape
}

func (shape *Shape2d) calcInternalDistancesTarget(plane *Plane2d){
	var direction int
	//angle ids are going down
	if shape.Points[0].AngleId > shape.Points[1].AngleId {
		direction = -1
	} else {
		direction = 1
	}

	edge := [2]int{0, 1}
	edgeAsLine := utils.Line(shape.Points[0].Vec, shape.Points[1].Vec)
	count := 0
	angleId := shape.Points[0].AngleId

	for count < len(plane.LineComponents) {
		intersect := utils.FindIntersect(plane.LineComponents[angleId], edgeAsLine)
		dist := shape.Weight * math.Sqrt(intersect[0] * intersect[0] + intersect[1] * intersect[1])
		plane.Distances[angleId] += dist

		angleId += direction

		if angleId == -1 {
			angleId = len(plane.LineComponents) - 1
		} else if angleId == len(plane.LineComponents){
			angleId = 0
		}

		if shape.Points[edge[1]].AngleId == angleId {
			edge0, edge1 := edge[0] + 1, edge[1] + 1

			if edge0 == len(shape.Points){
				edge[0] = 0
				edge[1] = edge1
			} else if edge1 == len(shape.Points){
				edge[1] = 0
				edge[0] = edge0
			} else {
				edge[0], edge[1] = edge0, edge1
			}
			edgeAsLine = utils.Line(shape.Points[edge[0]].Vec, shape.Points[edge[1]].Vec)
		}
		count++
	}
}

func (shape *Shape2d) calcInternalDistances(plane *Plane2d){
	angleId := shape.Points[shape.LowEdge[0]].AngleId + 1

	v0High, v1High := shape.Points[shape.HighEdge[0]], shape.Points[shape.HighEdge[1]]
	v0Low, v1Low := shape.Points[shape.LowEdge[0]], shape.Points[shape.LowEdge[1]]

	highEdgeAsLine := utils.Line(v0High.Vec, v1High.Vec)
	lowEdgeAsLine := utils.Line(v0Low.Vec, v1Low.Vec)

	for angleId != shape.EndAngleId {
		i0, i1 := utils.FindIntersect(plane.LineComponents[angleId], highEdgeAsLine),
			utils.FindIntersect(plane.LineComponents[angleId], lowEdgeAsLine)

		xDist, yDist := i1[0] - i0[0], i1[1] - i0[1]

		dist := shape.Weight * math.Sqrt(xDist * xDist + yDist * yDist)
		plane.Distances[angleId] += dist

		if v1High.AngleId == angleId {
			shape.shiftHighEdge()
			v0High, v1High = shape.Points[shape.HighEdge[0]], shape.Points[shape.HighEdge[1]]
			highEdgeAsLine = utils.Line(v0High.Vec, v1High.Vec)
		}

		if v1Low.AngleId == angleId {
			shape.shiftLowEdge()
			v0Low, v1Low = shape.Points[shape.LowEdge[0]], shape.Points[shape.LowEdge[1]]
			lowEdgeAsLine = utils.Line(v0Low.Vec, v1Low.Vec)
		}
		angleId++

		if angleId == len(plane.LineComponents){
			angleId = 0
		}
	}
}

func (shape *Shape2d) findStartingEdges(lowIdx int, highIdx int){
	if shape.Points[highIdx].Angle - shape.Points[lowIdx].Angle >= math.Pi {
		shape.findClockwiseBoundaries()
	} else {
		shape.HighEdge = [2]int{lowIdx, (lowIdx + 1) % len(shape.Points)}
		lowEdgeEnd := lowIdx - 1
		if lowEdgeEnd < 0 {
			lowEdgeEnd = len(shape.Points) - 1
		}
		shape.LowEdge = [2]int{lowIdx, lowEdgeEnd}
		shape.EndAngleId = shape.Points[highIdx].AngleId
	}
	shape.checkColinearStart()
}

func (shape *Shape2d) findClockwiseBoundaries(){
	
	pointsCopy := make([]*plane3d.FoundVertex2d, len(shape.Points))
	copy(pointsCopy, shape.Points)
	sort.Sort(VertexSorter(pointsCopy))

	maxAngleDiff := 0.0
	endIdx := 0

	for i := 1; i < len(pointsCopy); i++ {
		angleDiff := pointsCopy[i].Angle - pointsCopy[i - 1].Angle
		if angleDiff > maxAngleDiff {
			maxAngleDiff = angleDiff
			endIdx = i - 1
		}
	}

	endPoint, startPoint := pointsCopy[endIdx], pointsCopy[endIdx + 1]
	shape.EndAngleId = endPoint.AngleId

	for idx, point := range shape.Points {
		if point.Angle == startPoint.Angle {
			shape.HighEdge = [2]int{idx, (idx + 1) % len(shape.Points)}
			lowEdgeEnd := idx - 1
			if lowEdgeEnd < 0 {
				lowEdgeEnd = len(shape.Points) - 1
			}
			shape.LowEdge = [2]int{idx, lowEdgeEnd}
			break
		}
	}
}

func (shape *Shape2d) checkColinearStart(){
	if shape.Points[shape.HighEdge[0]].Angle == shape.Points[shape.HighEdge[1]].Angle {
		shape.shiftHighEdge()
	}
}

func (shape *Shape2d) shiftLowEdge(){
	low0, low1 := shape.LowEdge[0] - 1, shape.LowEdge[1] - 1

	if low1 < 0 {
		shape.LowEdge[1] = len(shape.Points) - 1
		shape.LowEdge[0] = low0
	} else if low0 < 0 {
		shape.LowEdge[0] = len(shape.Points) - 1
		shape.LowEdge[1] = low1
	} else {
		shape.LowEdge[0], shape.LowEdge[1] = low0, low1
	}
}

func (shape *Shape2d) shiftHighEdge(){
	high0, high1 := shape.HighEdge[0] + 1, shape.HighEdge[1] + 1
	if high0 == len(shape.Points) {
		shape.HighEdge[0] = 0
		shape.HighEdge[1] = high1
	} else if high1 == len(shape.Points) {
		shape.HighEdge[0] = high0
		shape.HighEdge[1] = 0
	} else {
		shape.HighEdge[0], shape.HighEdge[1] = high0, high1
	}
}
