package simulator

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

			}
		case "trimodal":
			{

			}
		}

	}

}
