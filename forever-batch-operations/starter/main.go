package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"

	forever_batch_operations "github.com/temporalio/samples-go/forever-batch-operations"
)

func main() {
	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	// start forever batching workflow
	workflowOptions := client.StartWorkflowOptions{
		ID:        forever_batch_operations.GetAccumulateAndBatchWorkflowID(),
		TaskQueue: "batch",
	}
	_, err = c.ExecuteWorkflow(context.Background(), workflowOptions, forever_batch_operations.AccumulateAndBatchWorkflow, nil)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	// start signaling workflow
	workflowOptions = client.StartWorkflowOptions{
		ID:        forever_batch_operations.GetSignalNewValuesWorkflowID(),
		TaskQueue: "batch",
	}
	_, err = c.ExecuteWorkflow(context.Background(), workflowOptions, forever_batch_operations.SignalNewValuesWorkflow)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}
}
