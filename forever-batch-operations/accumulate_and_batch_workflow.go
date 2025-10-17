package forever_batch_operations

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"go.temporal.io/sdk/workflow"
)

const (
	SIGNAL_READ_VALS    = "read_vals"
	SIGNAL_COMMIT_BATCH = "commit_batch"
)

func GetAccumulateAndBatchWorkflowID() string {
	return "AccumulateAndBatchWorkflowID"
}

// AccumulateAndBatchWorkflow keeps accumulating data via signals and when there
// is enough data, sends it to an activity.
func AccumulateAndBatchWorkflow(ctx workflow.Context, vals []string) error {

	logger := workflow.GetLogger(ctx)
	logger.Info("AccumulateAndBatchWorkflow workflow started", "vals", vals)

	// read incoming values here
	readValsCh := workflow.GetSignalChannel(ctx, SIGNAL_READ_VALS)
	// listen for when to commit
	commitBatchCh := workflow.NewChannel(ctx)

	selector := workflow.NewSelector(ctx)

	// start batching timer
	workflow.Go(ctx, func(ctx workflow.Context) {
		for {
			t := workflow.NewTimer(ctx, 10*time.Second)
			err := t.Get(ctx, nil)
			// error handling may be different
			if err != nil {
				continue
			}
			commitBatchCh.Send(ctx, nil)
		}
	})

	selector.AddReceive(readValsCh, func(c workflow.ReceiveChannel, more bool) {
		var incoming string
		c.Receive(ctx, &incoming)
		vals = append(vals, incoming)
	})

	shouldContinueAsNew := false
	selector.AddReceive(commitBatchCh, func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		logger.Info("commiting batch...")
		ao := workflow.ActivityOptions{
			StartToCloseTimeout: 1 * time.Minute,
		}
		actx := workflow.WithActivityOptions(ctx, ao)
		err := workflow.ExecuteActivity(actx, WriteBatchToFile, vals).Get(ctx, nil)
		if err != nil {
			// we couldn't write the batch. Simply continue-as-new with vals.
			logger.Error("failed to write batch", "error", err)
			shouldContinueAsNew = true
			return
		}
		vals = []string{}
		shouldContinueAsNew = true
	})

	for {
		selector.Select(ctx)
		if shouldContinueAsNew {
			return workflow.NewContinueAsNewError(ctx, AccumulateAndBatchWorkflow, vals)
		}
	}
}

func WriteBatchToFile(ctx context.Context, vals []string) error {
	// Write the values to this file. We can compare values here to values sent
	// to the workflow to see the durability.
	f, err := os.OpenFile("values_received.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer f.Close()

	f.WriteString(strings.Join(vals, "\n"))
	f.WriteString("\n")

	return nil
}
