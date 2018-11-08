package algorithm

import (
	"encoding/json"
	"testing"
	"../plane3d"
	//"../graph"
	"../utils"
	"math"
	//"log"
)

type TestShape struct {
	Points				[][2]float64 `json:"points"`
	Weight				float64 `json:"weight"`
	TangentDistances	[][4]float64 `json:"tangent_distances"`
}

type TestGroup struct {
	VertexGroup			*plane3d.VertexGroup
	TangentPoints		[]*plane3d.FoundVertex2d
	ExpectedDistances	[]float64
}


func TestShapeDistances(t *testing.T){
	plane2d, testGroups := getTestData()



	for idx, testGroup := range testGroups {
		if idx == 4 {
			//expected distances
			t.Logf("%+v\n", testGroup.ExpectedDistances)
			shape2d := newShape(testGroup.VertexGroup)
			t.Logf("%+v\n", shape2d)
			t.Logf("%+v\n", shape2d.Points[shape2d.LowEdge[0]])
			t.Logf("%+v\n", shape2d.Points[shape2d.LowEdge[1]])
			t.Logf("%+v\n", shape2d.Points[shape2d.HighEdge[0]])
			t.Logf("%+v\n", shape2d.Points[shape2d.HighEdge[1]])
			shape2d.calcInternalDistances(plane2d)

			/*if shape2d.MinVertex.AngleId == shape2d.StartVertex.AngleId {
				shape2d.calcInternalDistances(plane2d, shape2d.MaxVertex.AngleId)
			} else {
				t.Logf("HERE")
				shape2d.calcInternalDistances(plane2d, shape2d.MaxVertex.AngleId + 1)
				t.Logf("%+v\n", plane2d.Distances)
				shape2d.calcInternalDistances(plane2d, shape2d.StopVertex.AngleId)
			}*/

			//shape2d.calcInternalDistances(plane2d)
			//actual distances
			t.Logf("%+v\n", plane2d.Distances)

			/*for idx, dist := range distances {
				distExpected := testGroup.ExpectedDistances[idx]
				//if round(distExpected) != round(dist) {
				//	t.Errorf("Expected distance %f for path through points [0, 0] and %+v\n intersecting point group %+v\n, got distance %f.",
				//		distExpected, testGroup.TangentPoints[idx], shape2d, dist)
				//}
				t.Log(dist)
				t.Log(distExpected)
				t.Log("\n")
			}*/
			t.Log("\n\n")
			//t.Logf("%+v\n", shape2d.calcDistances(plane2d))
			//t.Log("\n")
		}
	}
}


func TestShiftEdges(t *testing.T){
	group := make([]*plane3d.FoundVertex2d, 4)
	shape := &Shape2d{group, 0, [2]int{1, 0}, [2]int{1, 2}, 0}
	shape.shiftLowEdge()

	if !(shape.LowEdge[0] == 0 && shape.LowEdge[1] == 3) {
		t.Errorf("Expected value of shifted low edge: [0 3], found %+v\n", shape.LowEdge)
	}

	shape.shiftHighEdge()
	shape.shiftHighEdge()
	shape.shiftHighEdge()

	if !(shape.HighEdge[0] == 0 && shape.HighEdge[1] == 1) {
		t.Errorf("Expected value of shifted high edge: [3 0], found %+v\n", shape.HighEdge)
	}
}

func TestColinearStart(t *testing.T){
	x0, y0 := 1.0, 1.0
	x1, y1 := 2.0, 2.0
	x2, y2 := -2.0, 2.0

	v0 := &plane3d.FoundVertex2d{Vec: [2]float64{x0, y0}, Angle: math.Atan2(y0, x0), AngleId: 0}
	v1 := &plane3d.FoundVertex2d{Vec: [2]float64{x1, y1}, Angle: math.Atan2(y1, x1), AngleId: 0}
	v2 := &plane3d.FoundVertex2d{Vec: [2]float64{x2, y2}, Angle: math.Atan2(y2, x2), AngleId: 0}

	group := &plane3d.VertexGroup{FoundVertices: []*plane3d.FoundVertex2d{v0, v1, v2}, 
		Weight: 0, TissueId: 0, TetId: 0}

	shape := newShape(group)
	
	if !(shape.HighEdge[0] == 1 && shape.HighEdge[1] == 2){
		t.Errorf("Expected value of high edge: [1 2], found %+v\n", shape.HighEdge)
	}

	if !(shape.LowEdge[0] == 0 && shape.LowEdge[1] == 2){
		t.Errorf("Expected value of low edge: [0 2], found %+v\n", shape.LowEdge)
	}
}

func TestCrossesNegativeX(t *testing.T){
	x0, y0 := 0.0, 1.0
	x1, y1 := -1.0, -1.0
	x2, y2 := -1.0, 1.0

	v0 := &plane3d.FoundVertex2d{Vec: [2]float64{x0, y0}, Angle: math.Atan2(y0, x0), AngleId: 0}
	v1 := &plane3d.FoundVertex2d{Vec: [2]float64{x1, y1}, Angle: math.Atan2(y1, x1), AngleId: 0}
	v2 := &plane3d.FoundVertex2d{Vec: [2]float64{x2, y2}, Angle: math.Atan2(y2, x2), AngleId: 0}

	group := &plane3d.VertexGroup{FoundVertices: []*plane3d.FoundVertex2d{v0, v1, v2}, 
		Weight: 0, TissueId: 0, TetId: 0}
	shape := newShape(group)

	if !(shape.HighEdge[0] == 0 && shape.HighEdge[1] == 1){
		t.Errorf("Expected value of high edge: [0 1], found %+v\n", shape.HighEdge)
	}

	if !(shape.LowEdge[0] == 0 && shape.LowEdge[1] == 2){
		t.Errorf("Expected value of low edge: [0 2], found %+v\n", shape.LowEdge)
	}
}

func getTestData()(*Plane2d, []*TestGroup){
	var testShapes []*TestShape
	json.Unmarshal(utils.OpenJsonFile("../mesh/test_3.json"), &testShapes)

	var testGroups []*TestGroup
	var vertexGroups []*plane3d.VertexGroup

	for _, testShape := range testShapes {
		vertices2d := make([]*plane3d.FoundVertex2d, len(testShape.Points))
		for i := 0; i < len(testShape.Points); i++ {
			point := testShape.Points[i]
			vertices2d[i] = &plane3d.FoundVertex2d{Vec: [2]float64{point[0], point[1]},
				Angle: math.Atan2(point[1], point[0]), AngleId: -1}
			//log.Printf("%+v\n", vertices2d[i])
		}

		
		//tetrahedron tissue and id don't matter here
		vertexGroup := &plane3d.VertexGroup{FoundVertices: vertices2d, Weight: testShape.Weight, 
			TissueId: 0, TetId: 0}
		vertexGroups = append(vertexGroups, vertexGroup)

		var tangentPoints []*plane3d.FoundVertex2d
		var distances []float64

		testGroups = append(testGroups, &TestGroup{vertexGroup, tangentPoints, 
			distances})
	}

	plane2d := newPlane2d(vertexGroups)

	for idx, testShape := range testShapes {
		for _, dist := range testShape.TangentDistances{
			for _, point := range plane2d.Points {
				if point.Vec[0] == dist[0] && point.Vec[1] == dist[1] {
					testGroups[idx].TangentPoints = append(testGroups[idx].TangentPoints,
						point)
					testGroups[idx].ExpectedDistances = append(testGroups[idx].ExpectedDistances,
						dist[2])
					break
				}
			}
		}
	}
	return plane2d, testGroups
}