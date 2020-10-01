package simulator

import (
	"math"
	"math/rand"
)

func getDistrInstance() float64 {
	//fmt.Println(servingNode)
	switch d := distribution; d {
	case "exp":
		{
			//fmt.Println("exponential distribution instance")
			return rand.ExpFloat64()
		}
	case "hyperexp":
		{
			//returns and hyper
			var p float64 = 0.5
			var rate1 float64 = 3
			var rate2 float64 = 5

			var P float64 = rand.Float64()
			var out float64

			if P <= p {
				out = (math.Log(1-rand.Float64()) / rate1) * -1
			} else {
				out = (math.Log(1-rand.Float64()) / rate2) * -1
			}

			return out
		}

	}
	return 0
}
