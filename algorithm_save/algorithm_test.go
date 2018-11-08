package algorithm

import (
	"encoding/json"
	"testing"
	"math"
)

type TestShape struct {
	Points				[][2]float64 `json:"points"`
	Weight				float64 `json:"weight"`
	TangentDistances	[][4]float64 `json:"tangent_distances"`
}

type TestGroup struct {
	VertexGroup			*VertexGroup
	TangentPoints		[]*Vertex2d
	ExpectedDistances	[]float64
}

func TestShapeDistances(t *testing.T){
	plane2d, testGroups := getTestData()
	//for _, group := range testGroups {
	//	t.Logf("%+v\n", group)
	//}
	/*for _, v := range plane2d.Points {
		t.Logf("%+v\n", v)
	}*/
	for idx, testGroup := range testGroups {
		if idx == 2 {
			//expected distances
			t.Logf("%+v\n", testGroup.ExpectedDistances)
			shape2d := newShape(testGroup.VertexGroup)

			if shape2d.MinVertex.AngleId == shape2d.StartVertex.AngleId {
				shape2d.calcInternalDistances(plane2d, shape2d.MaxVertex.AngleId)
			} else {
				t.Logf("HERE")
				shape2d.calcInternalDistances(plane2d, shape2d.MaxVertex.AngleId + 1)
				t.Logf("%+v\n", plane2d.Distances)
				shape2d.calcInternalDistances(plane2d, shape2d.StopVertex.AngleId)
			}

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

/*type ExpectedDistance struct {
	AngleId		int
	Distance	float64
}*/

func getTestData()(*Plane2d, []*TestGroup){
	var testShapes []*TestShape
	json.Unmarshal(openJsonFile("../mesh/test_1.json"), &testShapes)

	var testGroups []*TestGroup
	var vertexGroups []*VertexGroup

	for _, testShape := range testShapes {
		vertices2d := make([]*Vertex2d, len(testShape.Points))
		for i := 0; i < len(testShape.Points); i++ {
			point := testShape.Points[i]
			vertices2d[i] = &Vertex2d{[2]float64{point[0], point[1]}, 0, 0,
				math.Atan2(point[1], point[0]), -1}
		}
		vertexGroup := &VertexGroup{vertices2d, testShape.Weight}
		vertexGroups = append(vertexGroups, vertexGroup)

		var tangentPoints []*Vertex2d
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

	/*for _, point := range plane2d.Points {
		for idx, testShape := range testShapes {
			for _, tangentDist := range testShape.TangentDistances {
				if point.Vec[0] == tangentDist[0] && point.Vec[1] == tangentDist[1] {
					testGroups[idx].TangentPoints = append(testGroups[idx].TangentPoints,
						point)
					testGroups[idx].ExpectedDistances = append(testGroups[idx].ExpectedDistances,
						tangentDist[2])
					break
				}
			}
		}
	}*/
	return plane2d, testGroups
}

func round(num float64)(float64){
	return math.Round(num * 1000)/1000
}