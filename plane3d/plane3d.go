package plane3d

import (
	"../graph"
	//"log"
	"math"
	"github.com/gonum/matrix/mat64"
)

type Plane3d struct {
	//Id int
	Target 		[3]float64
	AxisX 		[3]float64
	AxisY 		[3]float64
	Normal 		[3]float64
	TargetDist 	float64
	Inverse		[3][3]float64
	//TargetTets	[]*graph.Tetrahedron
}

func newPlane3d(alpha float64, theta float64, target [3]float64)(*Plane3d){
	u := [3]float64{1, 0, 0}
	v := [3]float64{0, math.Cos(alpha), math.Sin(alpha)}
	w := [3]float64{0, -math.Sin(alpha), math.Cos(alpha)}

	data := []float64{u[0], u[1], u[2], v[0], v[1], v[2], w[0], w[1], w[2]}

	rotationX := mat64.NewDense(3, 3, data)

	u = [3]float64{math.Cos(theta), 0, math.Sin(theta)}
	v = [3]float64{0, 1, 0}
	w = [3]float64{-math.Sin(theta), 0, math.Cos(theta)}

	data = []float64{u[0], u[1], u[2], v[0], v[1], v[2], w[0], w[1], w[2]}

	rotationY := mat64.NewDense(3, 3, data)

	totalRotation := mat64.NewDense(3, 3, nil)

	totalRotation.Mul(rotationX, rotationY)
	//log.Printf("%+v\n", totalRotation)
	axisX := [3]float64{totalRotation.At(0, 0), totalRotation.At(0, 1), totalRotation.At(0, 2)}
	axisY := [3]float64{totalRotation.At(1, 0), totalRotation.At(1, 1), totalRotation.At(1, 2)}
	normal := [3]float64{totalRotation.At(2, 0), totalRotation.At(2, 1), totalRotation.At(2, 2)}
	//log.Println(normal)

	inverse := mat64.NewDense(3, 3, nil)
	inverse.Inverse(totalRotation)
	//log.Printf("%+v\n", inverse)

	//axisX := [3]float64{math.Cos(theta), -math.Sin(alpha) * math.Sin(theta), math.Cos(alpha) * -math.Sin(theta)}
	//axisY := [3]float64{0.0, math.Cos(alpha), -math.Sin(alpha)}
	inverseX := [3]float64{inverse.At(0, 0), inverse.At(0, 1), inverse.At(0, 2)}
	inverseY := [3]float64{inverse.At(1, 0), inverse.At(1, 1), inverse.At(1, 2)}
	inverseZ := [3]float64{inverse.At(2, 0), inverse.At(2, 1), inverse.At(2, 2)}

	//log.Println(inverseX)
	//log.Println(inverseY)
	//log.Println(inverseZ)
	
	targetDist := target[0] * normal[0] + target[1] * normal[1] + target[2] * normal[2]


	return &Plane3d{target, axisX, axisY, normal, targetDist, 
		[3][3]float64{inverseX, inverseY, inverseZ}}
}

//func newPlane3d(alpha float64, theta float64, target [3]float64)(*Plane3d){
	/*axisX := [3]float64{math.Cos(theta), -math.Sin(alpha) * math.Sin(theta), math.Cos(alpha) * math.Sin(theta)}
	axisY := [3]float64{0.0, math.Cos(alpha), -math.Sin(alpha)}
	axisZ := [3]float64{-math.Sin(theta), math.Sin(alpha) * math.Cos(theta), math.Cos(alpha) * math.Cos(theta)}*/
	// := [3]float64{math.Cos(theta), -math.Sin(theta) * math.Sin(alpha), -math.Sin(theta) * math.Cos(alpha)}
	//axisY := [3]float64{0, math.Cos(alpha), -math.Sin(alpha)}
	//axisZ := [3]float64{math.Sin(theta), math.Sin(alpha) * math.Cos(alpha), math.Cos(alpha) * math.Cos(theta)}
	//rotationX := []float64{1, 0, 0, 0, math.Cos(alpha), math.Sin(alpha), 0, 
	//	-math.Sin(alpha), math.Cos(alpha)}
	//rotationY := []float64{math.Cos(theta), 0, math.Sin(theta), 0, 1, 0, -math.Sin(theta),
	//	0, math.Cos(theta)}

	//data := []float64{axisX[0], axisX[1], axisX[2], axisY[0], axisY[1], axisY[2],
	//	axisZ[0], axisZ[1], axisZ[2]}
	//x := mat64.NewDense(3, 3, rotationX)
	//y := mat64.NewDense(3, 3, rotationY)

	/*x.MulElem(x, y)
	rotationMatrix := x.RawMatrix().Data

	axisX := [3]float64{rotationMatrix[0], rotationMatrix[1], rotationMatrix[2]}
	axisY := [3]float64{rotationMatrix[3], rotationMatrix[4], rotationMatrix[5]}
	normal := [3]float64{rotationMatrix[6], rotationMatrix[7], rotationMatrix[8]}
	
	log.Printf("%+v\n", rotationMatrix)

	//a := mat64.NewDense(3, 3, data)
	log.Printf("%+v\n", x)
	x.Inverse(x)
	log.Printf("%+v\n", x)
	rm := x.RawMatrix().Data
	reverseX := [3]float64{rm[0], rm[1], rm[2]}
	reverseY := [3]float64{rm[3], rm[4], rm[5]}
	reverseZ := [3]float64{rm[6], rm[7], rm[8]}
	inverse := [3][3]float64{reverseX, reverseY, reverseZ}
	log.Println(inverse)

	//normalMag := math.Sqrt( axisZ[0] * axisZ[0] + axisZ[1] * axisZ[1] + axisZ[2] * axisZ[2] )

	//reverseX := [3]float64{math.Cos(theta), 0, -math.Sin(-theta)}
	//reverseY := [3]float64{-math.Sin(-alpha) * math.Sin(-theta), math.Cos(-alpha), -math.Sin(-alpha) * math.Cos(-theta)}
	//reverseZ := [3]float64{math.Sin(-theta) * math.Cos(-alpha), math.Sin(-alpha), math.Cos(-alpha) * math.Cos(-theta)}
	//reverse := [3][3]float64{reverseX, reverseY, reverseZ}
	
	//log.Println(normalMag)
	//normalUnit := [3]float64{axisZ[0] / normalMag, axisZ[1] / normalMag, axisZ[2] / normalMag}
	//log.Println(normalUnit)
	targetDist := target[0] * normal[0] + target[1] * normal[1] + target[2] * normal[2]

	return &Plane3d{target, axisX, axisY, normal, targetDist, inverse}
}*/


func (plane *Plane3d) formVertexGroup(stack *Stack, tet *graph.Tetrahedron)([]*FoundVertex2d) {
	//log.Println("this happens")
	dotProducts := plane.calcDotProducts(tet)

	var vertices []*FoundVertex2d
	vertexPairs := make([][]int, 4)
	count := 0

	//log.Println(dotProducts)

	for i := 0; i < len(dotProducts); i++ {
		if dotProducts[i] == 0.0 {
            vertex := tet.Vertices[i]
			newVec := plane.get2DCoords([3]float64{vertex.Vec[0], vertex.Vec[1], vertex.Vec[2]})
			angle := math.Atan2(newVec[1], newVec[0])
            vertices = append(vertices, &FoundVertex2d{newVec, angle, -1})
            
            for _, tet := range vertex.Tetrahedrons {
                stack.push(&Node{tet})
            }
		} else {

			for j := i + 1; j < len(dotProducts); j++ {
				if dotProducts[i] * dotProducts[j] < 0.0 {
					//newVec := plane.calcIntersect(tet.Vertices[i], tet.Vertices[j])
					newVec := plane.calcIntersect(tet.Vertices[i], tet.Vertices[j])
					angle := math.Atan2(newVec[1], newVec[0])
					vertices = append(vertices, &FoundVertex2d{newVec, angle, -1})
					vertexPairs[count] = []int{i, j}
					count++
				}
			}
		}
	}
	if len(vertices) == 4 {
		vertices = []*FoundVertex2d{vertices[0], vertices[1], vertices[3], vertices[2]}
	}
	return vertices
}


/* calculates intersection point of plane and edge between two vertices.
 * these vertices are known to be intersected by the plane; there will not
 * be two vertices that are on the same side of the plane, or any that lie
 * in the plane.
 */ 

 func (plane *Plane3d) calcIntersect(v1 *graph.Vertex3d, v2 *graph.Vertex3d)([2]float64) {
	target := plane.Target
	w := []float64{v1.Vec[0] - target[0], v1.Vec[1] - target[1], v1.Vec[2] - target[2]}
	u := []float64{v2.Vec[0] - v1.Vec[0], v2.Vec[1] - v1.Vec[1], v2.Vec[2] - v1.Vec[2]}
	normalX, normalY, normalZ := plane.Normal[0], plane.Normal[1], plane.Normal[2]

	N := - (normalX * w[0] + normalY * w[1] + normalZ * w[2])
	D := (normalX * u[0] + normalY * u[1] + normalZ * u[2])
	coords2D := plane.get2DCoords([3]float64{v1.Vec[0] + (N / D) * u[0], 
		v1.Vec[1] + (N / D) * u[1], v1.Vec[2] + (N / D) * u[2]})

	return coords2D
	//return [3]float64{v1.Vec[0] + (N / D) * u[0], 
	//	v1.Vec[1] + (N / D) * u[1], v1.Vec[2] + (N / D) * u[2]}
}

func (plane *Plane3d) calcDotProducts(tet *graph.Tetrahedron)([]float64) {
	dotProducts := make([]float64, 4)
	normalX, normalY, normalZ := plane.Normal[0], plane.Normal[1], plane.Normal[2]

	for idx, vertex := range tet.Vertices {
		x, y, z := vertex.Vec[0] - plane.Target[0], vertex.Vec[1] - plane.Target[1], 
			vertex.Vec[2] - plane.Target[2]
		dotProducts[idx] = normalX * x + normalY * y + normalZ * z
	}
	return dotProducts
}

/* takes a 3d vertex that is known to lie in the plane and returns
 * the 2d coordinates of the vertex in that plane, where the origin of the
 * plane is the target point.
 */
func (plane *Plane3d) get2DCoords(I [3]float64)([2]float64){
	dist := []float64{I[0] - plane.Target[0], I[1] - plane.Target[1], I[2] - plane.Target[2]}
	axisX, axisY := plane.AxisX, plane.AxisY

	newX := (dist[0] * axisX[0]) + (dist[1] * axisX[1]) + (dist[2] * axisX[2])
	newY := (dist[0] * axisY[0]) + (dist[1] * axisY[1]) + (dist[2] * axisY[2])

	return [2]float64{newX, newY}
}

 /*func (plane *Plane3d) Get2DCoords(I [3]float64)([2]float64){
	//target, axisX, axisY, axisZ := plane.Target, plane.Reverse[0], plane.Reverse[1], plane.Reverse[2]
	//target, axisX, axisY, axisZ := plane.Target, plane.AxisX, plane.AxisY, plane.Normal
	target := plane.Target
	axisX := plane.Inverse[0]
	axisY := plane.Inverse[1]
	axisZ := plane.Inverse[2]
	log.Println()
	log.Println(axisX)
	log.Println(axisY)
	log.Println(axisZ)
	
	dist := []float64{I[0] - target[0], I[1] - target[1], I[2] - target[2]}
	log.Println(dist)

	xNew := axisX[0] * dist[0] + axisX[1] * dist[1] + axisX[2] * dist[2]
	yNew := axisY[0] * dist[0] + axisY[1] * dist[1] + axisY[2] * dist[2]
	zNew := axisZ[0] * dist[0] + axisZ[1] * dist[1] + axisZ[2] * dist[2]

	log.Println([3]float64{xNew, yNew, zNew})
	//angle := math.Atan2(yNew, xNew)

	//return &Vertex2d{xNew, yNew, angle, 0}
	return [2]float64{xNew, yNew}
}
*/