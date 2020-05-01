package main

import (
	"flag"

	"gitlab.com/emanuele232/mapreduce-simulator/pkg/simulator"
	//"math/rand"
)

func main() {
	rateControl := flag.Bool("rateControl", true, "is the flag control enabled?")
	nNodes := flag.Int("nNodes", 5, "the number of nodes")
	maxJobs := flag.Int("maxJobs", 50, "number of job to be completed to terminate the simulation")
	flag.Parse()

	simulator.Start(*rateControl, *nNodes, *maxJobs)

}
