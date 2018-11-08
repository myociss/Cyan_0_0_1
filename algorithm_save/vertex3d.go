package algorithm

type Vertex3d struct {
	Id 				int `json:"id"`
	Vec 			[3]float64 `json:"vec"`
	Tetrahedrons 	[]*Tetrahedron `json:"tetrahedrons"`
}

/* checks if this vertex has the same 3d coordinates as another vertex.
 */

func (vertex *Vertex3d) equals(v *Vertex3d)(bool) {
	return vertex.Vec[0] == v.Vec[0] && vertex.Vec[1] == v.Vec[1] && vertex.Vec[2] == v.Vec[2]
}

/* adds all tetrahedrons this vertex is a member of to the stack
 */
func (vertex *Vertex3d) addTets(stack *Stack, planeId int){

	for _, tet := range vertex.Tetrahedrons {
		stack.push(&Node{tet})
	}
}