package simulator

import (
	"math"
	"math/rand"
)

var messages int 

func newServiceTime() {
	if len(nodes[servingNode].serviceTasksQ) == 0 {
		timeOfCompletion[servingNode] = 0
	} else {

		switch r := rateControl; r {
		case "no":
			{
				//timeOfCompletion[servingNode] = (rand.ExpFloat64() / lambdas[servingNode])
				timeOfCompletion[servingNode] = getDistrInstance()
			}
		case "bimodal":
			{
				messages++
				if nodes[servingNode].nk >= 0 {

					var rate = float64(nodes[servingNode].nk + 1)
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					timeOfCompletion[servingNode] = t * rate

				} else {
					//timeOfCompletion[servingNode] = (rand.ExpFloat64() / lambdas[servingNode])
					timeOfCompletion[servingNode] = getDistrInstance()
				}
			}
		case "bimodal-fixed":
			{

				var nextNode int

				if servingNode == nNodes-1 {
					nextNode = 0
				} else {
					nextNode = servingNode + 1
				}

				nextQ := nodes[nextNode].lenJoin
				currentQ := nodes[servingNode].lenJoin

				messages++

				if currentQ <= nextQ {

					timeOfCompletion[servingNode] = getDistrInstance()

				} else {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					timeOfCompletion[servingNode] = t * 2
				}

			}
		case "trimodal":
			{
				var previousNode int
				var nextNode int

				switch servingNode {
				case 0:
					{
						previousNode = nNodes - 1
						nextNode = servingNode + 1
					}
				case nNodes - 1:
					{
						previousNode = servingNode - 1
						nextNode = 0
					}
				default:
					{
						previousNode = servingNode - 1
						nextNode = servingNode + 1
					}
				}

				previousQ := nodes[previousNode].lenJoin
				nextQ := nodes[nextNode].lenJoin
				currentQ := nodes[servingNode].lenJoin

				messages = messages + 2

				if (currentQ >= previousQ) && (currentQ >= nextQ) {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					rate := math.Max(float64(previousQ), float64(nextQ))

					if rate == 0 {
						timeOfCompletion[servingNode] = t
					} else {
						timeOfCompletion[servingNode] = t * rate
					}

				}
				if (currentQ <= previousQ) && (currentQ <= nextQ) {

					//timeOfCompletion[servingNode] = (rand.ExpFloat64() / lambdas[servingNode])
					timeOfCompletion[servingNode] = getDistrInstance()

				}
				if (currentQ > previousQ) && (currentQ < nextQ) {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					rate := (float64(previousQ) + float64(nextQ)) / 2
					timeOfCompletion[servingNode] = t * rate

				}

				if (currentQ < previousQ) && (currentQ > nextQ) {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					rate := (float64(previousQ) + float64(nextQ)) / 2
					timeOfCompletion[servingNode] = t * rate

				}

			}

		case "trimodal-fixed":
			{
				var previousNode int
				var nextNode int

				switch servingNode {
				case 0:
					{
						previousNode = nNodes - 1
						nextNode = servingNode + 1
					}
				case nNodes - 1:
					{
						previousNode = servingNode - 1
						nextNode = 0
					}
				default:
					{
						previousNode = servingNode - 1
						nextNode = servingNode + 1
					}
				}

				previousQ := nodes[previousNode].lenJoin
				nextQ := nodes[nextNode].lenJoin
				currentQ := nodes[servingNode].lenJoin

				messages = messages + 2

				if (currentQ >= previousQ) && (currentQ >= nextQ) {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					timeOfCompletion[servingNode] = t * 3

				}
				if (currentQ <= previousQ) && (currentQ <= nextQ) {

					//timeOfCompletion[servingNode] = (rand.ExpFloat64() / lambdas[servingNode])
					timeOfCompletion[servingNode] = getDistrInstance()

				}
				if (currentQ > previousQ) && (currentQ < nextQ) {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					timeOfCompletion[servingNode] = t * 2

				}

				if (currentQ < previousQ) && (currentQ > nextQ) {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					timeOfCompletion[servingNode] = t * 2

				}

			}
		case "bimodal-random":
			{
				randomNode := rand.Intn(nNodes - 1)
				randQ := nodes[randomNode].lenJoin
				currentQ := nodes[servingNode].lenJoin

				messages++

				if currentQ <= randQ {

					timeOfCompletion[servingNode] = getDistrInstance()

				} else {
					var t = getDistrInstance()
					//var t = (rand.ExpFloat64() / lambdas[servingNode])

					timeOfCompletion[servingNode] = t * 2
				}
			}
		}

	}

}
