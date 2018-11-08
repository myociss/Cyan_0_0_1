package web

/* SearchSpaceRequest contains the values entered for epsilon and the target
 * selected by the user. epsilon is not guaranteed to be valid and target
 * is not guaranteed not to be nil.
 */

type SearchSpaceRequest struct {
    EpsilonStr 	string
	TargetStr 	string
	TetIdStr	string
}

/* this is the description of the plane space to search: it can be described
 * by epsilon, where epsilon^2 is the number of discrete planes to search.
 * the target is a 3-element float array containing the 3d coordinates of
 * the target point selected by the user.
 */

type SearchSpace struct {
	Epsilon int
	Target 	[3]float64
	TetId	int
}