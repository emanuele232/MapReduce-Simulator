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
		fmt.Println("Usage of MRSimulator\n")
		fmt.Println("-rateControl: [true/false] enable/disable Service rate Control")
		fmt.Println("-nNodes: [Integer] Set the number of nodes")
		fmt.Println("-maxJobs: [Integer] Set the number of jobs to complete to terminate the simulation")
		fmt.Println("-help/-h: show this text")
		fmt.Println()

	} else {
		simulator.Start(*rateControl, *nNodes, *maxJobs, "exp")

	}

}
