package simulator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

const nNodes = 3
const maxJobs = 1000
const nPartsOfJob = 9

var nodes []Node
var servedJobs = 0
var nInputSliced = 0
var systemClock = 0.0
var jobSplitted = 0
var inputSplits []string //when a job is splitted this array is populated

var taskCompletion map[string]int
var currentTask string
var servingNode int
var nextTime float64
var lambdas [nNodes]float64

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
			timeStationaryLen:      0,
			nk:                     make(map[string]int),
			nk2:                    0,
			lambda:                 rand.Float64()})
	}
	//create and split the first job
	job := Job{jobSplitted, nPartsOfJob}
	inputSplits = job.splitJob()

	//generates the times in which the nodes end the computation of the map tasks
	for i := 0; i < nNodes; i++ {
		lambdas[i] = float64(i+1) * 5
		timeOfCompletion[i] = rand.ExpFloat64() / lambdas[i]

	}

	taskCompletion = make(map[string]int)
	arrivalTimes = make(map[string]float64)
	avgJoinLen = make([]float64, nNodes)

}

func sendTasksToQueues() {
	var nodeID = 0
	var task string

	for range inputSplits {
		task, inputSplits = inputSplits[0], inputSplits[1:]
		if nodeID == len(nodes) {
			nodeID = 0
		}
		nodes[nodeID].serviceTasksQ = append(nodes[nodeID].serviceTasksQ, task)
		arrivalTimes[task] = systemClock
		nodeID++
		/*
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
		*/
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
		var pattern = regexp.MustCompile(fmt.Sprint(s, "_[0-9]+$"))
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

	//debug purposes

	for servedJobs < maxJobs {

		nextTime = 0
		for i := range timeOfCompletion {
			if nextTime == 0 || timeOfCompletion[i] < nextTime {
				nextTime = timeOfCompletion[i]
				servingNode = i
			}
		}

		/*
			fmt.Println("---- ITERATION ----")
			for i2 := range nodes {
				fmt.Println(nodes[i2].serviceTasksQ)

			}
			fmt.Println(fmt.Sprint("node serving:", servingNode))

			for i2 := range nodes {
				fmt.Println(nodes[i2].joinTasksQ)

			}
		*/

		currentTask = nodes[servingNode].serviceTasksQ[0]
		nodes[servingNode].lenJoin = len(nodes[servingNode].joinTasksQ)

		updateCounters()

		//remove task from the service Q and adds it to the join Q
		nodes[servingNode].serviceTasksQ = nodes[servingNode].serviceTasksQ[1:]
		nodes[servingNode].joinTasksQ = append(nodes[servingNode].joinTasksQ, currentTask)
		nodes[servingNode].taskCompleted++
		nodes[servingNode].nk2++

		if servingNode == 0 {
			nodes[len(nodes)-1].nk2--
		} else {
			var nk2 = servingNode - 1
			nodes[nk2].nk2--
		}

		s := strings.Split(currentTask, "_")[0]
		nodes[servingNode].nk[s]++

		if servingNode == 0 {
			nodes[len(nodes)-1].nk[s]--
		} else {
			var nk2 = servingNode - 1
			nodes[nk2].nk[s]--
		}

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
			var rate float64
			if nodes[servingNode].nk2 >= 0 {
				rate = float64(nodes[servingNode].nk2 + 1)
				//var stime = rand.ExpFloat64() / lambda
				//var srate = 1 / stime

				var t = (rand.ExpFloat64() / lambdas[servingNode])
				_ = rate
				_ = t
				//timeOfCompletion[servingNode] = t * rate
				timeOfCompletion[servingNode] = (rand.ExpFloat64() / lambdas[servingNode])
			} else {
				timeOfCompletion[servingNode] = (rand.ExpFloat64() / lambdas[servingNode])

			}

		}
		fmt.Print()

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
