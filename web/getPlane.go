package web

import (
	"../algorithm"
	"encoding/json"
	"log"
	//"math"
	"net/http"
	"strconv"
	"strings"
	//"../plane3d"
)

type Group3d struct {
	Points		[][3]float64
	Weight		float64
	TissueId	int
	TetId		int
}

type PlaneRes struct {
	Groups		[]*Group3d
	Normal		[3]float64
	NormalDist	float64
	/*AxisX		[3]float64
	AxisY		[3]float64
	Normal		[3]float64
	NormalDist	float64*/
}

/*type PlaneRes struct {
	Shapes 		[]*Shape
	Alpha		float64
	Theta		float64
	Origin		[3]float64
	NormalDist	float64
}

type Shape struct {
	Points [][2]float64
}*/

func PlaneRequestHandler(w http.ResponseWriter, r *http.Request){
	//log.Println(r)

	reqAlpha, _ := r.URL.Query()["alpha"]
	reqTheta, _ := r.URL.Query()["theta"]
	reqTarget, _ := r.URL.Query()["target"]
	reqTetId, _ := r.URL.Query()["tetId"]
	reqEpsilon, _ := r.URL.Query()["epsilon"]

	//log.Println(reqTarget)

	alphaId, _ := strconv.Atoi(strings.Join(reqAlpha, ","))
	thetaId, _ := strconv.Atoi(strings.Join(reqTheta, ","))
	targetStrComponents := strings.Split(strings.Join(reqTarget, ","), ",")
	tetId, _ := strconv.Atoi(strings.Join(reqTetId, ","))
	epsilon, _ := strconv.Atoi(strings.Join(reqEpsilon, ","))
	
	var target [3]float64
	for idx, strComp := range targetStrComponents {
		//log.Println(strComp)
		floatComp, _ := strconv.ParseFloat(strComp, 64)
		target[idx] = floatComp
	}
	
	
	//points, normal, dist := algorithm.GetPlane(target, tetId, alphaId, thetaId, epsilon)
	//points, normal, dist := algorithm.GetPlane(target, tetId, alphaId, thetaId, epsilon)
	planeGen := algorithm.GetPlane(target, tetId, alphaId, thetaId, epsilon)
	groups := make([]*Group3d, len(planeGen.FoundVertices))
	plane := planeGen.Plane
	//axisX, axisY, axisZ := plane.AxisX, plane.AxisY, plane.Normal

	/*alpha := (float64(alphaId) * math.Pi) / float64(epsilon)
	theta := (float64(thetaId) * math.Pi) / float64(epsilon)

	axisX := [3]float64{math.Cos(-theta), 0, -math.Sin(-theta)}
	axisY := [3]float64{-math.Sin(-alpha) * math.Sin(-theta), math.Cos(-alpha), -math.Sin(-alpha) * math.Cos(-theta)}
	axisZ := [3]float64{math.Sin(-theta) * math.Cos(-alpha), math.Sin(-alpha), math.Cos(-alpha) * math.Cos(-theta)}
	*/
	/*axisX := [3]float64{math.Cos(-theta), -math.Sin(-alpha) * math.Sin(-theta), math.Cos(-alpha) * -math.Sin(-theta)}
	axisY := [3]float64{0.0, math.Cos(-alpha), math.Sin(-alpha)}
	axisZ := [3]float64{math.Sin(-theta), -math.Sin(-alpha) * math.Cos(-theta), math.Cos(-alpha) * math.Cos(-theta)}
	*/
	//if axisZ[2] < 0 {
	//	axisZ = [3]float64{-axisZ[0], -axisZ[1], -axisZ[2]}
	//}
	///target := plane.Target
	//log.Println(axisZ)

	axisX, axisY, axisZ := plane.Inverse[0], plane.Inverse[1], plane.Inverse[2]
	//log.Println("here")
	log.Println(len(planeGen.FoundVertices))

	for idx, pointGroup := range planeGen.FoundVertices{
		points := make([][3]float64, len(pointGroup.FoundVertices))
		for pointIdx, vertex := range pointGroup.FoundVertices {
			//log.Println()
			//pointNormal := [3]float64{point[0] - target[0], point[1] - target[1], point[2] - target[2]}
			//log.Println(pointNormal)
			//points[pointIdx] = plane.Get2DCoords(point)
			
			//newPoint := plane.Get2DCoords(point)
			point := vertex.Vec
			x3d := point[0] * axisX[0] + point[1] * axisX[1]
			y3d := point[0] * axisY[0] + point[1] * axisY[1]
			z3d := point[0] * axisZ[0] + point[1] * axisZ[1]
			points[pointIdx] = [3]float64{x3d + target[0], y3d + target[1], z3d + target[2]}
		}
		//log.Println(points)
		groups[idx] = &Group3d{points, pointGroup.Weight, pointGroup.TissueId, pointGroup.TetId}
	}

	a, err := json.Marshal(&PlaneRes{groups, plane.Normal, plane.TargetDist})
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
	}
    w.Write(a)
}