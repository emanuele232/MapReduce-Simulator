package simulator

import (
	"fmt"
	"html/template"
	"math"
	"os"
)

func reset() {
	nodes = nil
	inputSplits = nil
	taskCompletion = nil

	arrivalTimes = nil
	averageDelayQueue = nil
	avgJoinLen = nil
	avgServiceDelay = nil
	energeticConsumption = nil
	timeOfCompletion = nil
	jobTotalDelay = 0
	energyOnTime = 0
	energyOnNJob = 0
	totalEnergyConsumed = 0
}

func printResults() {
	fmt.Println("\n\n\n\n------------------------------------------")
	fmt.Println(fmt.Sprint("SIMULATOR VERSION:", version))
	fmt.Println(fmt.Sprint("\nParameters:", "\n", "rate control:", rateControl, "\nnodes:", nNodes, "\nJobs:", maxJobs, "\ndistribution:", distribution))
	fmt.Println("------------------------------------------")

	fmt.Println(fmt.Sprintln("System clock:", systemClock))
	fmt.Println(fmt.Sprintln("Energetic consumption:", energeticConsumption))
	fmt.Println(fmt.Sprintln("Served jobs:", servedJobs))
	fmt.Println(fmt.Sprintln("Served tasks :", servedTasks))
	//fmt.Println(arrivalTimes)

	fmt.Print("The energy consumed by the system:")
	fmt.Println(totalEnergyConsumed)
	fmt.Print("\nEnergy consumed on unit of time: ")
	fmt.Println(energyOnTime)

	fmt.Print("\nEnergy consumed on job: ")
	fmt.Println(energyOnNJob)

	fmt.Print("The number of messages needed for ")
	fmt.Print(rateControl)
	fmt.Print(" rate control")
	fmt.Print(": ")
	fmt.Println(messages)

	fmt.Print("Job Avg delay: ")
	fdelay := jobTotalDelay / float64(maxJobs)
	fmt.Print(fdelay)

	fmt.Print("Job Avg service delay: ")
	fmt.Print(avgServiceDelayf)

	fmt.Println("\n------------------------------")

	fmt.Println("Expected average customers in join queues:")
	for i := range nodes {
		fmt.Print(fmt.Sprint("Node-", i, ":"))
		fmt.Println(avgJoinLen[i])
	}

	fmt.Println("\n------------------------------")

	fmt.Println("Energetic consumption in nodes:")
	for i := range nodes {
		fmt.Print(fmt.Sprint("Node-", i, ":"))
		fmt.Println(energeticConsumption[i])
	}

}

func populateTemplate() {
	type Templatedata struct {
		Clock       float64
		Energy      float64
		EnergyJob   float64
		EnergyTime  float64
		Avgdelay    float64
		AvgSvcDelay float64
		Messages    int
		Node01      float64
		Node02      float64
		Node11      float64
		Node12      float64
		Node21      float64
		Node22      float64
		Node31      float64
		Node32      float64
		Node41      float64
		Node42      float64
	}

	data := Templatedata{
		math.Round(systemClock*100) / 100,
		math.Round(totalEnergyConsumed*100) / 100,
		math.Round(energyOnNJob*100) / 100,
		math.Round(energyOnTime*100) / 100,
		math.Round(jobTotalDelay/float64(maxJobs)*100) / 100,
		math.Round(avgServiceDelayf*100) / 100,
		messages,
		math.Round(avgJoinLen[0]*100) / 100,
		math.Round(energeticConsumption[0]*100) / 100,
		math.Round(avgJoinLen[1]*100) / 100,
		math.Round(energeticConsumption[1]*100) / 100,
		math.Round(avgJoinLen[2]*100) / 100,
		math.Round(energeticConsumption[2]*100) / 100,
		math.Round(avgJoinLen[3]*100) / 100,
		math.Round(energeticConsumption[3]*100) / 100,
		math.Round(avgJoinLen[4]*100) / 100,
		math.Round(energeticConsumption[4]*100) / 100,
	}

	t, err := template.ParseFiles("../../templates/result_table")
	if err != nil {
		panic(err)
	}
	fmt.Println("\n\n\n")
	err = t.Execute(os.Stdout, data)

	reset()
}
