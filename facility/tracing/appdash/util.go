package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
)

// Returns the remote collector address.
func startAppdashServer(appdashPort, tracingServPort int) {
	store := appdash.NewMemoryStore()

	// Listen on any available TCP port locally.
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: tracingServPort})
	if err != nil {
		log.Fatal(err)
	}

	// Start an Appdash collection server that will listen for spans and
	// annotations and add them to the local collector (stored in-memory).
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go cs.Start()

	// Print the URL at which the web UI will be running.
	appdashURLStr := fmt.Sprintf("http://localhost:%d", appdashPort)
	appdashURL, err := url.Parse(appdashURLStr)
	if err != nil {
		log.Fatalf("Error parsing %s: %s", appdashURLStr, err)
	}
	fmt.Printf("To see your traces, go to %s/traces\n", appdashURL)

	// Start the web UI
	tapp, err := traceapp.New(nil, appdashURL)
	if err != nil {
		log.Fatalf("Error creating traceapp: %v", err)
	}
	tapp.Store = store
	tapp.Queryer = store

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appdashPort), tapp))
}
