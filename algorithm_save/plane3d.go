package algorithm

import (
	"math"
	"log"
)

type Plane struct {
	Id int
	Target 		[3]float64
	AxisX 		[3]float64
	AxisY 		[3]float64
	Normal 		[3]float64
	TargetDist 	float64
	TargetTets	[]*Tetrahedron
}

type VertexGroup3d struct {
	Vertices	[][3]float64
	Weight		float64
	TissueId	int
	TetId 		int
}

var GET_SELECTION_SLICES = 0
var GET_PLANE_SPACE = 1

func (plane *Plane) findCrossSection(threadId int, mode int, tetsLength int) ([]*VertexGroup3d) {
	var tetsChecked []bool
	if mode == GET_SELECTION_SLICES {
		tetsChecked = make([]bool, tetsLength)
		for i := 0; i < tetsLength; i++ {
			tetsChecked[i] = false
		}
	}
	stack := newStack()

	for _, tet := range plane.TargetTets {
		stack.push(&Node{tet})
	}

	var allVertices []*VertexGroup3d

	for stack.Count > 0 {
		tet := stack.pop().Tet

		if tet != nil {
			checked := (mode == GET_SELECTION_SLICES && tetsChecked[tet.Id]) ||
				(mode == GET_PLANE_SPACE && tet.PlaneLastChecked[threadId] == plane.Id)

			if !checked {
				group := &VertexGroup3d{plane.formVertexGroup(stack, tet), tet.Weight, 
					tet.TissueId, tet.Id}
		
				if len(group.Vertices) > 2 {
					allVertices = append(allVertices, group)
		
					for _, neighbor := range tet.Neighbors {
						stack.push(&Node{neighbor})
					}
				}
		
				if mode == GET_SELECTION_SLICES {
					tetsChecked[tet.Id] = true
				} else {
					tet.PlaneLastChecked[threadId] = plane.Id
				}
			}
		}
	}
	return allVertices
}



/* finds a group of poins that lie on the faces of a tetrahedron within
 * this plane.
 */

 func (plane *Plane) formVertexGroup(stack *Stack, tet *Tetrahedron)([][3]float64) {
	dotProducts := plane.calcDotProducts(tet)

	var vertices3d [][3]float64
	vertexPairs := make([][]int, 4)
	count := 0

	for i := 0; i < len(dotProducts); i++ {
		if dotProducts[i] == 0.0 {
			tet.Vertices[i].addTets(stack, plane.Id)
			vec := tet.Vertices[i].Vec
			//newVec := plane.get2DCoords([]float64{vec[0], vec[1], vec[2]})
			//angle := math.Atan2(newVec[1], newVec[0])
			//vertices3d = append(vertices3d, &Vertex2d{newVec, tet.Id, tet.TissueId,
			//	angle, -1})
			vertices3d = append(vertices3d, vec)
		} else {

			for j := i + 1; j < len(dotProducts); j++ {
				if dotProducts[i] * dotProducts[j] < 0.0 {
					//newVec := plane.calcIntersect(tet.Vertices[i], tet.Vertices[j])
					//angle := math.Atan2(newVec[1], newVec[0])
					//vertices3d = append(vertices3d, &Vertex2d{newVec, tet.Id, tet.TissueId, 
					//	angle, -1})
					vec := plane.calcIntersect(tet.Vertices[i], tet.Vertices[j])
					vertices3d = append(vertices3d, vec)
					vertexPairs[count] = []int{i, j}
					count++
				}
			}
		}
	}
	if len(vertices3d) == 4 {
		vertices3d = [][3]float64{vertices3d[0], vertices3d[1], vertices3d[3], vertices3d[2]}
	}
	return vertices3d
}

/* finds two vectors in the plane through the target and having its x axis 
 * tilted by alpha radians and y axis tilted by theta radians. finds the
 * normal vector to the plane. finds distance from (0,0,0) in the normal
 * direction of the plane to the target point
 */

func newPlane(id int, alpha float64, theta float64, target [3]float64, tets []*Tetrahedron)(*Plane){
	axisX := [3]float64{math.Cos(alpha), 0.0, math.Sin(alpha)}
	log.Println(axisX)
	axisY := [3]float64{0.0, math.Cos(theta), math.Sin(theta)}
	normal := [3]float64{axisX[1] * axisY[2] - axisX[2] * axisY[1], axisX[2] * axisY[0] - axisX[0] * axisY[2],
		axisX[0] * axisY[1] - axisX[1] * axisY[0]}
	normalMag := math.Sqrt( normal[0] * normal[0] + normal[1] * normal[1] + normal[2] * normal[2] )
	normalUnit := [3]float64{normal[0] / normalMag, normal[1] / normalMag, normal[2] / normalMag}
	log.Println(normalUnit)

	/*axisXMag := math.Sqrt( axisX[0] * axisX[0] + axisX[1] * axisX[1] + axisX[2] * axisX[2] )
	axisYMag := math.Sqrt( axisY[0] * axisY[0] + axisY[1] * axisY[1] + axisY[2] * axisY[2] )
	normalMag := math.Sqrt( normal[0] * normal[0] + normal[1] * normal[1] + normal[2] * normal[2] )

	axisXUnit := [3]float64{axisX[0] / axisXMag, axisX[1] / axisXMag, axisX[2] / axisXMag}
	axisYUnit := [3]float64{axisY[0] / axisYMag, axisY[1] / axisYMag, axisY[2] / axisYMag}
	normalUnit := [3]float64{normal[0] / normalMag, normal[1] / normalMag, normal[2] / normalMag}
	log.Println(normalUnit)

	normalUnitPos := [3]float64{math.Sqrt(normalUnit[0] * normalUnit[0]), math.Sqrt(normalUnit[1] * normalUnit[1]), 
		math.Sqrt(normalUnit[2] * normalUnit[2])}

	targetDist := target[0] * normalUnitPos[0] + target[1] * normalUnitPos[1] + target[2] * normalUnitPos[2]*/
	targetDist := target[0] * normalUnit[0] + target[1] * normalUnit[1] + target[2] * normalUnit[2]

	return &Plane{id, target, axisX, axisY, normalUnit, targetDist, tets}
}

/* finds the dot product of every vertex in tet. vertices on one side of
 * the plane will have a positive value for the dot product and vertices 
 * on the other side will have a negative value. vertices on the plane will
 * have a value of 0. this will be used to determine if the plane cuts the
 * edge between two vertices in a tetrahedron
 */

func (plane *Plane) calcDotProducts(tet *Tetrahedron)([]float64) {
	dotProducts := make([]float64, 4)
	normalX, normalY, normalZ := plane.Normal[0], plane.Normal[1], plane.Normal[2]

	for idx, vertex := range tet.Vertices {
		x, y, z := vertex.Vec[0] - plane.Target[0], vertex.Vec[1] - plane.Target[1], 
			vertex.Vec[2] - plane.Target[2]
		dotProducts[idx] = normalX * x + normalY * y + normalZ * z
	}
	return dotProducts
}

/* calculates intersection point of plane and edge between two vertices.
 * these vertices are known to be intersected by the plane; there will not
 * be two vertices that are on the same side of the plane, or any that lie
 * in the plane.
 */ 

func (plane *Plane) calcIntersect(v1 *Vertex3d, v2 *Vertex3d)([3]float64) {
	target := plane.Target
	w := []float64{v1.Vec[0] - target[0], v1.Vec[1] - target[1], v1.Vec[2] - target[2]}
	u := []float64{v2.Vec[0] - v1.Vec[0], v2.Vec[1] - v1.Vec[1], v2.Vec[2] - v1.Vec[2]}
	normalX, normalY, normalZ := plane.Normal[0], plane.Normal[1], plane.Normal[2]

	N := - (normalX * w[0] + normalY * w[1] + normalZ * w[2])
	D := (normalX * u[0] + normalY * u[1] + normalZ * u[2])
	//coords2D := plane.get2DCoords([]float64{v1.Vec[0] + (N / D) * u[0], 
	//	v1.Vec[1] + (N / D) * u[1], v1.Vec[2] + (N / D) * u[2]})
	return [3]float64{v1.Vec[0] + (N / D) * u[0], 
		v1.Vec[1] + (N / D) * u[1], v1.Vec[2] + (N / D) * u[2]}

	//return coords2D
}

/* takes a 3d vertex that is known to lie in the plane and returns
 * the 2d coordinates of the vertex in that plane, where the origin of the
 * plane is the target point.
 */

func (plane *Plane) get2DCoords(I [3]float64)([2]float64){
	target, axisX, axisY := plane.Target, plane.AxisX, plane.AxisY

	dist := []float64{I[0] - target[0], I[1] - target[1], I[2] - target[2]}

	xNew := axisX[0] * dist[0] + axisX[1] * dist[1] + axisX[2] * dist[2]
	yNew := axisY[0] * dist[0] + axisY[1] * dist[1] + axisY[2] * dist[2]

	//angle := math.Atan2(yNew, xNew)

	//return &Vertex2d{xNew, yNew, angle, 0}
	return [2]float64{xNew, yNew}
}
