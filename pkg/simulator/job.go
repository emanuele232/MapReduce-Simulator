package simulator

import "fmt"

type Job struct {
	id     int
	nTasks int
}

func (j Job) splitJob() []string {
	a := make([]string, 0)
	for i := 0; i < j.nTasks; i++ {
		//generates an unique id for the task (idJob_nTask)
		a = append(a, fmt.Sprint(j.id, "_", i))
	}
	jobSplitted++
	return a
}
