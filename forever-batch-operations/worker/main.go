package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	forever_batch_operations "github.com/temporalio/samples-go/forever-batch-operations"
)

func main() {

	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "batch", worker.Options{})

	w.RegisterWorkflow(forever_batch_operations.AccumulateAndBatchWorkflow)
	w.RegisterWorkflow(forever_batch_operations.SignalNewValuesWorkflow)
	w.RegisterActivity(forever_batch_operations.WriteBatchToFile)
	w.RegisterActivity(forever_batch_operations.WriteValToFile)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
