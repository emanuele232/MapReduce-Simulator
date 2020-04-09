package simulator

var arrivalTimes map[string]float64
var averageDelayQueue []float64
var avgJoinLen []float64
var avgServiceDelay []float64

/*
update statistical counters
*/
func updateCounters() {
	nodes[servingNode].totalDelay += systemClock - arrivalTimes[currentTask]
}

func computeAvgLen() {

	for n := range nodes {
		for i := range nodes[n].totalTimeStationaryLen {
			avgJoinLen[n] += float64(i) * nodes[n].totalTimeStationaryLen[i]
		}
		avgJoinLen[n] = avgJoinLen[n] / systemClock
	}
}

func computeAvgDelay() {
	for n := range nodes {
		var avg = nodes[n].totalDelay / float64(nodes[n].taskCompleted)
		avgServiceDelay = append(avgServiceDelay, avg)
	}
}

func computeStatistics() {
	computeAvgLen()
	computeAvgDelay()
}
