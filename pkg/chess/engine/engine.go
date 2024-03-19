package engine

func Engine() (chan bool, chan string) {
	frEngine := make(chan string)
	toEngine := make(chan bool)

	return toEngine, frEngine
}
