# Mapreduce Simulator

A Mapreduce Simulator in Go

# Functioning
MRSimulator is an event simulator where the only event happening is the *completion of a task from one of the nodes*.
By design, the arrival rate of the jobs (then splitted and sent to queues) is non existent, because an entire new job is provided and sent to the queues when a task is completed.

When the Simulator Starts, it initializes its variables and then takes the first job and sends the tasks in order of nodeID to the nodes.

then the main cycle starts , where the following operations are completed:

- it finds the node that finishes its task first.
- updates statistical counters (avgDelay)
- removes the first task from the service node's *service's queue* and adds it to the serving node's *join queue*.
- updates the rate control's *nk* parameter
- updates system clock
- updates the time in which every node finish its task (parallel computation)
- checks if all tasks of a job are complete, if so eleminate every task of that job from the join queues.
- update statistical counters (avg Len)
- sends another job splitted in tasks to the nodes
- creates a new service time for the serving node

# Statistical counters
1) In the first statistical counters update it updates the delay of the serving node adding to the total delay the sysclock - the arrival time of the completed task.
2) In the second statistical counters update it updates the avg length of the join queues. It adds to the time in which there are a number n of customers in the queue the system clock - the time since the number of customers changed.

**Run 1** (Rate control Enabled, nodes = 5 jobs completed = 2000)

clock: 4038.4493166194065

Expected average customers in join queues:

- Node-0: 2.6026293689061966
- Node-1: 1.5522582461973133
- Node-2: 2.59549639797326
- Node-3: 2.863583549769582
- Node-4: 3.81332084353698

**Run 2** (Rate control Disabled, nodes = 5 jobs completed = 2000)

clock: 2324.7671910895183

Expected average customers in join queues:

- Node-0: 763.9511836162252
- Node-1: 2.6275376290063615
- Node-2: 166.98908626851687
- Node-3: 80.2378673738745
- Node-4: 279.04602776673966