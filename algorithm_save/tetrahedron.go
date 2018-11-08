package algorithm

import (
	//"log"
)

type Tetrahedron struct {
	Id					int
	Weight				float64
	TissueId			int
	Vertices			[4]*Vertex3d
	Neighbors			[4]*Tetrahedron
	PlaneLastChecked	[]int
}

/* this function checks if a point that has been selected by the user is
 * contained within tet. it checks all tetrahedron faces (groups of three
 * vertices) to see if the remaining point in the tetrahedron is on the same
 * side as point. it returns false if a face is found where the point is
 * not on the same side as the remaining vertex. see
 * https://stackoverflow.com/questions/25179693/how-to-check-whether-the-point-is-in-the-tetrahedron-or-not
 * for more details about procedure.
 */

func (tet *Tetrahedron) containsPoint(point [3]float64)(bool) {

	for idx, v := range tet.Vertices {
		//contains the three points of the face that does not include v
		var opposite []*Vertex3d
		for idx2, v2 := range tet.Vertices {
			if idx != idx2 {
				opposite = append(opposite, v2)
			}
		}

		//get 3d coordinates of opposite face points
		facePt0, facePt1, facePt2 := opposite[0].Vec, opposite[1].Vec, opposite[2].Vec

		//get two vectors contained in the opposite face
		vec0 := []float64{facePt1[0] - facePt0[0], facePt1[1] - facePt0[1], facePt1[2] - facePt0[2]}

		vec1 := []float64{facePt2[0] - facePt0[0], facePt2[1] - facePt0[1], facePt2[2] - facePt0[2]}

		
		//get normal of opposite face
		normal := []float64{vec0[1] * vec1[2] - vec0[2] * vec1[1], vec0[2] * vec1[0] - vec0[0] * vec1[2],
			vec0[0] * vec1[1] - vec0[1] * vec1[0]}


		v4Diff := []float64{v.Vec[0] - facePt0[0], v.Vec[1] - facePt0[1], v.Vec[2] - facePt0[2]}
		newPointDiff := []float64{point[0] - facePt0[0], point[1] - facePt0[1], point[2] - facePt0[2]}

		v4Dot := normal[0] * v4Diff[0] + normal[1] * v4Diff[1] + normal[2] * v4Diff[2]
		newPointDot := normal[0] * newPointDiff[0] + normal[1] * newPointDiff[1] + normal[2] * newPointDiff[2]

		//log.Println("asdf")
		//log.Println(v4Dot)
		//log.Println(newPointDot)

		if v4Dot * newPointDot < 0.0 {
			return false
		}
	}
	return true
}

/* returns true if any two vertices in tet have the same 3d coordinates, or
 * if all vertices have the same x, y or z values. this accounts for the
 * malformed tetrahedrons.
 */

func (tet *Tetrahedron) isFlat()(bool){
	for idx0, v0 := range tet.Vertices {
		for idx1, v1 := range tet.Vertices {
			if idx0 != idx1 {
				if v0.equals(v1) {
					//log.Println("here")
					//log.Println(tet.Id)
					return true
				}
			}
		}
	}

	singleX := true
	singleY := true
	singleZ := true

	for i := 1; i < len(tet.Vertices); i++ {
		if tet.Vertices[i].Vec[0] != tet.Vertices[i-1].Vec[0] {
			singleX = false
		}
		if tet.Vertices[i].Vec[1] != tet.Vertices[i-1].Vec[1] {
			singleY = false
		}
		if tet.Vertices[i].Vec[2] != tet.Vertices[i-1].Vec[2] {
			singleZ = false
		}
	}
	return singleX || singleY || singleZ
}