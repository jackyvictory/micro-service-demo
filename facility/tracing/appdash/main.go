package main

import "flag"

var (
	appdashPort     = flag.Int("appdash.port", 8700, "Run appdash locally on this port.")
	tracingServPort = flag.Int("tracingServ.port", 8701, "Run tracing server locally on this port.")
)

func main() {
	flag.Parse()
	startAppdashServer(*appdashPort, *tracingServPort)
}
