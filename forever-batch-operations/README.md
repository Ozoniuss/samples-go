This sample shows how to implement a workflow that continuously listens for incoming tasks via signals and batches them periodically. It also illustrates how to simulate a ticker using timers.

### Steps to run this sample:
1) Run a [Temporal service](https://github.com/temporalio/samples-go/tree/main/#how-to-use).
2) Run the following command to start the worker
```
go run forever-batch-operations/worker/main.go
```
3) Run the following command to start the example
```
go run forever-batch-operations/starter/main.go
```

Note that the workflows will continue running even after you stop the worker.

Compare the values_received.txt and values_sent.txt files. Values should be written in the same order.
