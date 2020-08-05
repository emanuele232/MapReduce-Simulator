package simulator

import (
	"fmt"
	"math"
)

var arrivalTimes map[string]float64
var averageDelayQueue []float64
var avgJoinLen []float64
var avgServiceDelay []float64
var energeticConsumption float64

/*
update statistical counters
*/
func updateDelay() {
	nodes[servingNode].totalDelay = nodes[servingNode].totalDelay + systemClock - arrivalTimes[currentTask]
}

func updateAvgLen() {
	if nodes[servingNode].lenJoin != len(nodes[servingNode].joinTasksQ) {
		nodes[servingNode].totalTimeStationaryLen[nodes[servingNode].lenJoin] +=
			systemClock - nodes[servingNode].timeStationaryLen
		nodes[servingNode].timeStationaryLen = systemClock
	}
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

func computeEnergeticConsumption() {
	energeticConsumption = math.Pow((float64(servedJobs) / systemClock), 3.0)
	fmt.Println(float64(servedJobs) / systemClock)
}

func computeStatistics() {
	computeAvgLen()
	computeAvgDelay()
	computeEnergeticConsumption()
}
