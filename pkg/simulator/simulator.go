package simulator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

const nNodes = 5
const maxJobs = 1000
const lenQ = 10
const nPartsOfJob = 60

var nodes []Node
var servedJobs = 0
var nInputSliced = 0
var systemClock = 0.0
var jobSplitted = 0
var inputSplits []string //when a job is splitted this array is populated
var rate = 1.00
var taskCompletion map[string]int
var currentTask string
var servingNode int
var nextTime float64

/*
	timeOfCompletion stores the time needed for a node to complete
	its map task (position in array = id node)
*/
var timeOfCompletion [nNodes]float64

func initialize() {

	//initialize nodes
	for i := 0; i < nNodes; i++ {
		nodes = append(nodes, Node{
			lenService:             0,
			lenJoin:                0,
			serviceTasksQ:          make([]string, 0),
			joinTasksQ:             make([]string, 0),
			totalDelay:             0.0,
			taskCompleted:          0,
			totalTimeStationaryLen: make(map[int]float64),
			timeStationaryLen:      0})
	}
	//create and split the first job
	job := Job{jobSplitted, nPartsOfJob}
	inputSplits = job.splitJob()

	//generates the times in which the nodes end the computation of the map tasks
	for i := 0; i < nNodes; i++ {
		timeOfCompletion[i] = rand.ExpFloat64() / rate
	}

	taskCompletion = make(map[string]int)
	arrivalTimes = make(map[string]float64)
	avgJoinLen = make([]float64, nNodes)

}

func sendTasksToQueues() {
	var nodeID = 0
	var nFullQueues = 0
	var task string

	for range inputSplits {
		task, inputSplits = inputSplits[0], inputSplits[1:]

		for nodeID < len(nodes) {
			if len(nodes[nodeID].serviceTasksQ) < lenQ {
				nodes[nodeID].serviceTasksQ = append(nodes[nodeID].serviceTasksQ, task)
				arrivalTimes[task] = systemClock
				nodeID++
				if nodeID >= len(nodes) {
					nodeID = 0
				}
				break
			} else {
				for n := range nodes {
					if len(nodes[n].serviceTasksQ) >= lenQ {
						nFullQueues++
					} else {
						nFullQueues = 0
						break
					}
				}
			}
			nodeID++
			if nodeID >= len(nodes) {
				nodeID = 0
			}
			if nFullQueues == len(nodes) {
				inputSplits = append([]string{task}, inputSplits...)
				break
			}
		}
		if nFullQueues == len(nodes) {
			break
		}
	}
}

func reduce() {
	/*
		keeps track of how many parts of this job
		are completed
	*/
	s := strings.Split(currentTask, "_")[0]
	taskCompletion[s] = taskCompletion[s] + 1

	/*
		if every split of the job is completed we
		remove every split from the join queues
	*/
	if taskCompletion[s] == nPartsOfJob {
		var pattern = regexp.MustCompile(fmt.Sprint(s, "_[0-9]$"))
		for n := range nodes {
			var e = 0
			for range nodes[n].joinTasksQ {
				var match = pattern.MatchString(nodes[n].joinTasksQ[e])
				if match {
					nodes[n].joinTasksQ = remove(nodes[n].joinTasksQ, e)

				} else {
					e++
				}

			}
		}
		servedJobs++
	}
}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func Start() {
	initialize()
	sendTasksToQueues()

	for servedJobs < maxJobs {

		nextTime = 0
		for i := range timeOfCompletion {
			if nextTime == 0 || timeOfCompletion[i] < nextTime {
				nextTime = timeOfCompletion[i]
				servingNode = i
			}
		}

		currentTask = nodes[servingNode].serviceTasksQ[0]
		nodes[servingNode].lenJoin = len(nodes[servingNode].joinTasksQ)

		updateCounters()

		//remove task from the service Q and adds it to the join Q
		nodes[servingNode].serviceTasksQ = nodes[servingNode].serviceTasksQ[1:]
		nodes[servingNode].joinTasksQ = append(nodes[servingNode].joinTasksQ, currentTask)
		nodes[servingNode].taskCompleted++

		//advance system clock
		systemClock = systemClock + nextTime

		/*
			since the computation is cuncurrent between the nodes we
			update the time that the other nodes need to complete their
			tasks
		*/
		for i := range timeOfCompletion {
			timeOfCompletion[i] = timeOfCompletion[i] - nextTime
		}

		reduce()

		if nodes[servingNode].lenJoin != len(nodes[servingNode].joinTasksQ) {
			nodes[servingNode].totalTimeStationaryLen[nodes[servingNode].lenJoin] +=
				systemClock - nodes[servingNode].timeStationaryLen
			nodes[servingNode].timeStationaryLen = systemClock
		}

		if len(inputSplits) == 0 {
			job := Job{jobSplitted, nPartsOfJob}
			inputSplits = job.splitJob()
		}

		sendTasksToQueues()

		if len(nodes[servingNode].serviceTasksQ) == 0 {
			timeOfCompletion[servingNode] = 0
		} else {
			timeOfCompletion[servingNode] = rand.ExpFloat64() / rate
		}

		//debug purposes
		/*
			fmt.Println("---- ITERATION ----")
			for i2 := range nodes {
				fmt.Println(len(nodes[i2].joinTasksQ))

			}
		*/
	}
	computeStatistics()
	printResults()
}

func printResults() {
	fmt.Println(fmt.Sprintln("System clock:", systemClock))
	//fmt.Println(arrivalTimes)
	fmt.Println("Expected average customers in join queues:")
	for i := range nodes {
		fmt.Print(fmt.Sprint("Node-", i, ":"))
		fmt.Println(avgJoinLen[i])
	}

	fmt.Println("\n------------------------------")

	fmt.Println("Expected average delay in service queues:")
	for i := range nodes {
		fmt.Print(fmt.Sprint("Node-", i, ":"))
		fmt.Println(avgServiceDelay[i])
	}
}
