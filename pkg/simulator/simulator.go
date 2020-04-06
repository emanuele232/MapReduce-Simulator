package simulator

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

const nNodes = 3
const maxJobs = 100
const lenQ = 100

var nodes []Node
var jobServed = 0
var nInputSliced = 0
var systemClock = 0.0
var jobSplitted = 0
var inputSplits []string //when a job is splitted this array is populated
var rate = 1.00
var taskCompletion map[string]int

/*
	timeOfCompletion stores the time needed for a node to complete
	its map task (position in array = id node)
*/
var timeOfCompletion [nNodes]float64

func sendTasksToQueues() {
	var nodeID = 0
	var nFullQueues = 0
	var task string

	for range inputSplits {
		task, inputSplits = inputSplits[0], inputSplits[1:]

		for nodeID < len(nodes) {
			if len(nodes[nodeID].serviceTasksQ) < lenQ {
				nodes[nodeID].serviceTasksQ = append(nodes[nodeID].serviceTasksQ, task)
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

func initialize() {

	//initialize nodes
	for i := 0; i < nNodes; i++ {
		nodes = append(nodes, Node{
			lenService:    0,
			lenJoin:       0,
			serviceTasksQ: make([]string, 0),
			joinTasksQ:    make([]string, 0)})
	}
	//create and split the first job
	job := Job{jobSplitted, 3}
	inputSplits = job.splitJob()

	//generates the times in which the nodes end the computation of the map tasks
	for i := 0; i < nNodes; i++ {
		timeOfCompletion[i] = rand.ExpFloat64() / rate
	}

	taskCompletion = make(map[string]int)

}

func remove(s []string, i int) []string {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
}

func Start() {
	initialize()
	sendTasksToQueues()

	for jobServed < maxJobs {

		var servingNode int
		var nextTime float64
		for i := range timeOfCompletion {
			if nextTime == 0 || timeOfCompletion[i] < nextTime {
				nextTime = timeOfCompletion[i]
				servingNode = i
			}
		}

		//advance system clock
		systemClock = systemClock + nextTime

		for i := range timeOfCompletion {
			timeOfCompletion[i] = timeOfCompletion[i] - nextTime
		}

		//remove task from the service Q and adds it to the join Q
		task := nodes[servingNode].serviceTasksQ[0]
		nodes[servingNode].serviceTasksQ = nodes[servingNode].serviceTasksQ[1:]
		nodes[servingNode].joinTasksQ = append(nodes[servingNode].joinTasksQ, task)

		s := strings.Split(task, "_")[0]
		taskCompletion[s] = taskCompletion[s] + 1
		var pattern = regexp.MustCompile(fmt.Sprint(s, "_[0-9]$"))

		if taskCompletion[s] == 3 {
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
		}

		if len(inputSplits) == 0 {
			job := Job{jobSplitted, 3}
			inputSplits = job.splitJob()
		}

		sendTasksToQueues()

		if len(nodes[servingNode].serviceTasksQ) == 0 {
			timeOfCompletion[servingNode] = 0
		} else {
			timeOfCompletion[servingNode] = rand.ExpFloat64() / rate
		}

	}
}
