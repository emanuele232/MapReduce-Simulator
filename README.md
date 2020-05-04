# Mapreduce Simulator

A Mapreduce Simulator in Go

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