package simulator

//Node definition
type Node struct {
	lenService        int
	lenJoin           int
	serviceTasksQ     []string
	joinTasksQ        []string
	totalDelay        float64
	taskCompleted     int
	timeStationaryLen float64
	nk                int
	lambda            float64

	/*
		time in which the queue has n  elements in it (n = index of the slice)
	*/
	totalTimeStationaryLen map[int]float64
}
