package simulator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

const nPartsOfJob = 6

var nodes []Node
var servedJobs int
var nInputSliced int
var systemClock float64
var jobSplitted int
var inputSplits []string //when a job is splitted this array is populated
var taskCompletion map[string]int
var currentTask string
var servingNode int
var nextTime float64
var lambdas []float64
var timeOfCompletion []float64
var lastNodes int
var nodeSendQ int
var rateControl bool
var nNodes int
var maxJobs int

func initialize() {
	servedJobs = 0
	systemClock = 0
	nodeSendQ = 0
	servedJobs = 0
	servedJobs = 0
	nInputSliced = 0
	systemClock = 0
	jobSplitted = 0
	taskCompletion = make(map[string]int)
	arrivalTimes = make(map[string]float64)
	avgJoinLen = make([]float64, nNodes)
	lambdas = make([]float64, nNodes)
	timeOfCompletion = make([]float64, nNodes)

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
			nk:                     0,
			lambda:                 rand.Float64()})
	}
	//create and split the first job
	job := Job{jobSplitted, nPartsOfJob}
	inputSplits = job.splitJob()

	//generates the times in which the nodes end the computation of the map tasks
	for i := 0; i < nNodes; i++ {
		lambdas[i] = rand.Float64()
		//lambdas[i] = 1

	}

}

/*
Sends all the tasks of a job to queues in order of nodeID, starting from
the node where it stopped last usage.
*/
func sendTasksToQueues() {
	var task string

	for range inputSplits {
		task, inputSplits = inputSplits[0], inputSplits[1:]
		if nodeSendQ == len(nodes) {
			nodeSendQ = 0
		}
		if len(nodes[nodeSendQ].serviceTasksQ) == 0 {
			timeOfCompletion[nodeSendQ] = rand.ExpFloat64() / lambdas[nodeSendQ]
		}
		nodes[nodeSendQ].serviceTasksQ = append(nodes[nodeSendQ].serviceTasksQ, task)
		arrivalTimes[task] = systemClock
		nodeSendQ++
	}
}

/*
	Controls is every task of a job is comlpeted, if it is the case
	removes every task from the join queues
*/
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

//remove an element from an array
func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

/*
	calculates a new timeOfCompletion for the servingNode,
	it considers :
		-if the servingQ is empty
		-if ratecontrol is enabled
		-if the parameter nk is >= 0
*/
func newServiceTime() {
	if len(nodes[servingNode].serviceTasksQ) == 0 {
		timeOfCompletion[servingNode] = 0
	} else {
		var rate float64
		if nodes[servingNode].nk >= 0 && rateControl {
			rate = float64(nodes[servingNode].nk + 1)
			var t = (rand.ExpFloat64() / lambdas[servingNode])

			timeOfCompletion[servingNode] = t * rate
		} else {
			timeOfCompletion[servingNode] = (rand.ExpFloat64() / lambdas[servingNode])

		}

	}

}

//Start the main cycle of the simulator
func Start(rc bool, n int, jobs int) {
	rateControl = rc
	nNodes = n
	maxJobs = jobs

	initialize()
	sendTasksToQueues()

	//debug purposes

	for servedJobs < maxJobs {

		nextTime = 0
		for i := range timeOfCompletion {
			if (nextTime == 0 || timeOfCompletion[i] < nextTime) && timeOfCompletion[i] != 0 {
				nextTime = timeOfCompletion[i]
				servingNode = i
			}
		}

		currentTask = nodes[servingNode].serviceTasksQ[0]
		nodes[servingNode].lenJoin = len(nodes[servingNode].joinTasksQ)

		updateDelay()

		//remove task from the service Q and adds it to the join Q
		nodes[servingNode].serviceTasksQ = nodes[servingNode].serviceTasksQ[1:]
		nodes[servingNode].joinTasksQ = append(nodes[servingNode].joinTasksQ, currentTask)
		nodes[servingNode].taskCompleted++
		nodes[servingNode].nk++

		if servingNode == 0 {
			nodes[len(nodes)-1].nk--
		} else {
			var nk = servingNode - 1
			nodes[nk].nk--
		}

		//advance system clock
		systemClock = systemClock + nextTime

		/*
			since the computation is cuncurrent between the nodes we
			update the time that the other nodes need to complete their
			tasks
		*/
		for i := range timeOfCompletion {
			if timeOfCompletion[i] != 0 {
				timeOfCompletion[i] = timeOfCompletion[i] - nextTime
			}
		}

		//if a job finished this tick, reduce it
		reduce()

		//update of statistical counters
		updateAvgLen()

		if len(inputSplits) == 0 {
			job := Job{jobSplitted, nPartsOfJob}
			inputSplits = job.splitJob()
		}

		//sends tasks to all serviceQs
		sendTasksToQueues()

		//next service time for Serving node
		newServiceTime()

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
