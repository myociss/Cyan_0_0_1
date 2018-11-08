package algorithm

/* maintains a stack containing tetrahedrons to be tested for intersection
 * with a plane. code source: https://gist.github.com/moraes/2141121
 */

type Node struct {
	Tet *Tetrahedron
}

func newStack() *Stack {
	return &Stack{}
}

type Stack struct {
	Nodes []*Node
	Count int
}

func (s *Stack) push(n *Node) {
	s.Nodes = append(s.Nodes[:s.Count], n)
	s.Count++
}

func (s *Stack) pop() *Node {
	if s.Count == 0 {
		return nil
	}
	s.Count--
	return s.Nodes[s.Count]
}