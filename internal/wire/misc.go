package wire

func StopHere() {
	stopper := EmptyG[chan struct{}]()

	stopper <- struct{}{}
}
