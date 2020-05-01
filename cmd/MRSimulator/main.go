package main

import (
	"flag"
	"fmt"

	"gitlab.com/emanuele232/mapreduce-simulator/pkg/simulator"
	//"math/rand"
)

func main() {
	rateControl := flag.Bool("rateControl", true, "is the flag control enabled?")
	nNodes := flag.Int("nNodes", 5, "the number of nodes")
	maxJobs := flag.Int("maxJobs", 50, "number of job to be completed to terminate the simulation")
	h := flag.Bool("h", false, "display help")
	help := flag.Bool("help", false, "display help")
	flag.Parse()

	if *h || *help {
		fmt.Println("display menu")
	} else {
		simulator.Start(*rateControl, *nNodes, *maxJobs)

	}

}
