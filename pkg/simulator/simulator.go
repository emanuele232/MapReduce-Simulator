package simulator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

const version = 1.0

var nPartsOfJob int
var distribution string
var nodes []Node
var servedJobs int
var servedTasks int
var nInputSliced int
var systemClock float64
var jobSplitted int
var inputSplits []string //when a job is splitted this array is populated
var taskCompletion map[string]int
var currentTask string
var servingNode int
var nextTime float64
var timeOfCompletion []float64
var expRate float64
var arrivalRate float64

//var lastNodes int
var nodeSendQ int
var rateControl string
var nNodes int
var maxJobs int
var lambdas [5]float64
var totalEnergyConsumed float64
var nextJobTime float64
var isServingIteration bool

func initialize() {
	servedJobs = 0
	systemClock = 0
	nodeSendQ = 0
	servedJobs = 0
	servedTasks = 0
	nInputSliced = 0
	systemClock = 0
	jobSplitted = 0
	servedJobs = 0
	servedTasks = 0
	nInputSliced = 0
	totalEnergyConsumed = 0
	nextJobTime = 0

	expRate = 5.0
	arrivalRate = 10.0

	taskCompletion = make(map[string]int)
	arrivalTimes = make(map[string]float64)
	avgJoinLen = make([]float64, nNodes)
	energeticConsumption = make([]float64, nNodes)
	timeOfCompletion = make([]float64, nNodes)

	nPartsOfJob = nNodes

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
	/*create and split the first job
	job := Job{jobSplitted, nPartsOfJob}
	inputSplits = job.splitJob()
	*/

	/*generates the times in which the nodes end the computation of the map tasks
	for i := 0; i < nNodes; i++ {
		lambdas[i] = rand.Float64()
		//lambdas[i] = 1

	}
	*/
	messages = 0

	totalEnergyConsumed = 0

	nextJobTime = rand.ExpFloat64() * 20

}

func checkpercentage() {
	if servedJobs == maxJobs/10 {
		fmt.Println("10%")
	}
	if servedJobs == maxJobs/5 {
		fmt.Println("20%")
	}
	if servedJobs == maxJobs/2 {
		fmt.Println("50%")
	}
}

/*
Sends all the tasks of a job to queues in order of nodeID, starting from
the node where it stopped last usage.
*/
func sendTasksToQueues() {
	var task string

	job := Job{jobSplitted, nPartsOfJob}
	inputSplits = job.splitJob()

	arrivalTimes[string(job.id)] = systemClock

	/*
		for range inputSplits {
			task, inputSplits = inputSplits[0], inputSplits[1:]
			if nodeSendQ == len(nodes) {
				nodeSendQ = 0
			}
			if len(nodes[nodeSendQ].serviceTasksQ) == 0 {
				timeOfCompletion[nodeSendQ] = rand.ExpFloat64() / lambdas[nodeSendQ]
				//!TODO cambiare con diversa distribuizione

				//timeOfCompletion[nodeSendQ] = getDistrInstance()
			}
			nodes[nodeSendQ].serviceTasksQ = append(nodes[nodeSendQ].serviceTasksQ, task)
			arrivalTimes[task] = systemClock
			nodeSendQ++
		}
	*/

	for i := range inputSplits {
		task, inputSplits = inputSplits[0], inputSplits[1:]

		if len(nodes[i].serviceTasksQ) == 0 {
			timeOfCompletion[i] = getDistrInstance()
			//!TODO cambiare con diversa distribuizione

			//timeOfCompletion[nodeSendQ] = getDistrInstance()
		}
		nodes[i].serviceTasksQ = append(nodes[i].serviceTasksQ, task)
		arrivalTimes[task] = systemClock

	}

	nextJobTime = rand.ExpFloat64() / 10

}

/*
	Controls is every task of a job is completed, if it is the case
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
		jobTotalDelay = jobTotalDelay + (systemClock - arrivalTimes[s+"_0"])

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

//Start the main cycle of the simulator
func Start(rc string, n int, jobs int, distr string) {
	rateControl = rc
	nNodes = n
	maxJobs = jobs
	distribution = distr

	initialize()

	for servedJobs < maxJobs {
		lenCheck := 1
		for n := range nodes {
			lenCheck = lenCheck * len(nodes[n].serviceTasksQ)
		}
		if lenCheck == 0 {
			sendTasksToQueues()
		}

		nextTime = 0
		for i := range timeOfCompletion {
			if (nextTime == 0 || timeOfCompletion[i] < nextTime) && timeOfCompletion[i] != 0 {
				nextTime = timeOfCompletion[i]
				servingNode = i
				isServingIteration = true
			}
		}

		if nextJobTime < nextTime {
			isServingIteration = false
			nextTime = nextJobTime
			sendTasksToQueues()
			nextJobTime = rand.ExpFloat64() / arrivalRate
		}

		// fmt.Println(fmt.Sprint("serving: ", isServingIteration))
		// fmt.Println(fmt.Sprint("Systemclock: ", systemClock))
		// fmt.Println(fmt.Sprint("Working node:", servingNode))
		// fmt.Println(fmt.Sprint("Next time:", nextTime))
		// fmt.Println(fmt.Sprint("Served tasks:", servedTasks))
		//
		// fmt.Println(fmt.Sprint("time of completion:", timeOfCompletion))
		// for i := range nodes {
		// 	fmt.Println(fmt.Sprint("queue node ", i, ": ", len(nodes[i].serviceTasksQ)))
		// }

		//advance system clock
		systemClock = systemClock + nextTime

		//actions that are exclusive of when we are serving a task
		if isServingIteration {
			currentTask = nodes[servingNode].serviceTasksQ[0]
			nodes[servingNode].lenJoin = len(nodes[servingNode].joinTasksQ)

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

			servedTasks++

			updateDelay()

			updateEnergeticConsumption()

			//if a job finished this tick, reduce it
			reduce()

			//update of statistical counters
			updateAvgLen()

			//next service time for Serving node
			newServiceTime()
		}

		// fmt.Println("------------------------------\n")

		/*
			since the computation is cuncurrent between the nodes we
			update the time that the other nodes need to complete their
			tasks
		*/
		var bb = false

		for i := range timeOfCompletion {
			if i != servingNode {
				timeOfCompletion[i] = timeOfCompletion[i] - nextTime
				if timeOfCompletion[i] < 0 {
					bb = true
				}
			}
		}
		if bb {
			fmt.Println(fmt.Sprint("next:", nextTime))

			fmt.Println(timeOfCompletion)
		}

		if isServingIteration {
			nextJobTime = nextJobTime - nextTime
		}

		//sends tasks to all serviceQs
		//sendTasksToQueues()

	}
	computeStatistics()
	printResults()
	populateTemplate()
	reset()

}
