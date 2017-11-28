### txsub
`txsub` provides the machinery that can be used to submit transactions directly to the core and track their progress.  It also helps to hide some of the complex asynchronous nature of transaction submission, waiting to respond to submitters when no definitive state is known.

### Usage
`txsub.System` the struct that ties all the interfaces together.
1. Init it with your own implementation of required interface or use default instead. Example:
```
System{
		Pending:           NewDefaultSubmissionList(),
		Submitter:         NewDefaultSubmitter(httpClient, "<URL_TO_CORE>"),
        Results:           usageSpecificResultsProvider,
        NetworkPassphrase: "network_passphrase",
}
```
2. Start `system.Tick` in separate goroutine.