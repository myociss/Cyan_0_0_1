package utils


func Line(vec0 [2]float64, vec1 [2]float64)([3]float64) {
	A := vec0[1] - vec1[1]
	B := vec1[0] - vec0[0]
	C := vec0[0] * vec1[1] - vec1[0] * vec0[1]
	return [3]float64{A, B, -C}
}

func FindIntersect(line0 [3]float64, line1 [3]float64)([2]float64){

	D  := line0[0] * line1[1] - line0[1] * line1[0]
    Dx := line0[2] * line1[1] - line0[1] * line1[2]
	Dy := line0[0] * line1[2] - line0[2] * line1[0]
	
	return [2]float64{Dx / D, Dy / D}
}