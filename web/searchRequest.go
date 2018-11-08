package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"log"
)

/* handles requests to generate and search epsilon^2 planes through a
 * target point. responds with an error if request is invalid; responds
 * with page containing mesh if valid.
 */

func SearchSpaceRequestHandler(w http.ResponseWriter, r *http.Request) {
	epsilon, _ := r.URL.Query()["epsilon"]
	target, _ := r.URL.Query()["target"]
	tetId, _ := r.URL.Query()["tetId"]
	ssr := &SearchSpaceRequest{strings.Join(epsilon, ""), strings.Join(target, ""),
			strings.Join(tetId, "")} 
	
	searchSpace, ssrErrors := ssr.validate()
	if ssrErrors != nil {
		errorMsg, _ := json.Marshal(ssrErrors)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(errorMsg)
	} else {
		log.Printf("%+v\n", searchSpace)
		w.WriteHeader(http.StatusOK)
	}
}


/* the validate function checks that the user entered a valid number for
 * epsilon and selected a target point, then creates a SearchSpace
 */

func (ssr *SearchSpaceRequest) validate() (*SearchSpace, error) {
	var allErrors []string

	epsilon, err := strconv.Atoi(ssr.EpsilonStr)

	if err != nil || epsilon < 1 || epsilon > 64 {
		allErrors = append(allErrors, "Please enter a value for epsilon between 1 and 64, inclusive.")
	}

	if len(ssr.TargetStr) == 0 {
		allErrors = append(allErrors, "Please select a target point.")
	}

	if len(allErrors) > 0 {
		return nil, errors.New(strings.Join(allErrors, "\n"))
	}

	targetStrComponents := strings.Split(ssr.TargetStr, ",")
	//target := make([]float64, 3)
	var target [3]float64

	for idx, strComp := range targetStrComponents {
		floatComp, _ := strconv.ParseFloat(strComp, 64)
		target[idx] = floatComp
	}

	tetId, _ := strconv.Atoi(ssr.TetIdStr)

	return &SearchSpace{epsilon, target, tetId}, nil
}