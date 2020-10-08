package simulator

import (
	"math"
	//"fmt"
)

var arrivalTimes map[string]float64
var averageDelayQueue []float64
var avgJoinLen []float64
var avgServiceDelay []float64
var energeticConsumption []float64
var jobTotalDelay float64
var energyOnTime float64
var energyOnNJob float64

/*
update statistical counters
*/
func updateDelay() {
	nodes[servingNode].totalDelay = nodes[servingNode].totalDelay + systemClock - arrivalTimes[currentTask]
	//fmt.Println(fmt.Sprintln(systemClock, "-", arrivalTimes[currentTask], "=", systemClock-arrivalTimes[currentTask]))

}

func updateAvgLen() {
	if nodes[servingNode].lenJoin != len(nodes[servingNode].joinTasksQ) {
		nodes[servingNode].totalTimeStationaryLen[nodes[servingNode].lenJoin] +=
			systemClock - nodes[servingNode].timeStationaryLen
		nodes[servingNode].timeStationaryLen = systemClock
	}
}

func updateEnergeticConsumption() {
	var f = math.Pow(1/timeOfCompletion[servingNode], 2)
	energeticConsumption[servingNode] = energeticConsumption[servingNode] + (f * timeOfCompletion[servingNode])
	if servingNode == 0 && (f*timeOfCompletion[servingNode]) < 0 {
		/*
			fmt.Println(fmt.Sprintln("f:", f))
			fmt.Println(fmt.Sprintln("toc:", timeOfCompletion[servingNode]))
		*/
	}

	/*
		fmt.Print("node: ")
		fmt.Println(servingNode)
		fmt.Print("time:")
		fmt.Println(timeOfCompletion[servingNode])
		fmt.Print("partial: ")
		fmt.Println(math.Pow(1/timeOfCompletion[servingNode], 3))
		fmt.Print("added: ")
		fmt.Println((f * timeOfCompletion[servingNode]))
		fmt.Print("total: ")
		fmt.Println(energeticConsumption[servingNode])
		fmt.Println()
	*/

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

func computeTotalEnergy() {
	for i := range nodes {
		totalEnergyConsumed = totalEnergyConsumed + energeticConsumption[i]
	}
	energyOnTime = totalEnergyConsumed / systemClock
	energyOnNJob = totalEnergyConsumed / float64(maxJobs)
}

func computeStatistics() {
	computeAvgLen()
	computeAvgDelay()
	computeTotalEnergy()
}
