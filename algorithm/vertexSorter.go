package algorithm

import (
	"../plane3d"
)

type VertexSorter []*plane3d.FoundVertex2d

func (a VertexSorter) Len() int           { return len(a) }
func (a VertexSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VertexSorter) Less(i, j int) bool { return a[i].Angle < a[j].Angle }