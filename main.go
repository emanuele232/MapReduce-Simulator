package main

import (
	"fmt"
	"math/rand"
)

const nNodes = 3
const maxJobs = 100
const lenQ = 2

var nodes []Node
var jobServed = 0
var nInputSliced = 0
var systemClock = 0
var lastJobSplitted = 0
var inputSplits []string //when a job is splitted this array is populated
var rate = 1.00

/*
timeOfCompletion stores the time needed for a node to complete
its map task (position in array = id node)
*/
var timeOfCompletion [nNodes]float64

type Node struct {
	lenService    int
	lenJoin       int
	serviceTasksQ []string
	joinTasksQ    []string
}

type Job struct {
	id     int
	nTasks int
}

func splitJob(job Job) []string {
	a := make([]string, 0)
	for i := 0; i < job.nTasks; i++ {
		//generates an unique id for the task (idJob_nTask)
		a = append(a, fmt.Sprint(job.id, "_", i))
	}
	return a
}

/*
sendTasksToQueue populate the queues of the nodes until there are no
tasks available at the moment or the queues are full
*/
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
	job := Job{lastJobSplitted, 10}
	inputSplits = splitJob(job)

	//generates the times in which the nodes end the computation of the map tasks
	for i := 0; i < nNodes; i++ {
		timeOfCompletion[i] = rand.ExpFloat64() / rate
	}

}

func main() {

	_ = timeOfCompletion
	_ = inputSplits
	_ = jobServed
	_ = maxJobs
	_ = nInputSliced
	_ = systemClock

	initialize()
	sendTasksToQueues()
	fmt.Println(nodes[0].serviceTasksQ)
	fmt.Println(nodes[1].serviceTasksQ)
	fmt.Println(nodes[2].serviceTasksQ)

	fmt.Println(inputSplits)

	/*
		main iteration, there is only one event that we want to address
		to advance the system clock and that is the completion of a map task

		for jobServed < maxJobs {
		}
	*/
}
