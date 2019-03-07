package txsub

type listener chan<- fullResult

func send(l chan fullResult, res fullResult) <-chan fullResult {
	l <- res
	close(l)
	return l
}

func newResultListener() listener {
	ch := make(chan fullResult, 1)

	return ch
}
