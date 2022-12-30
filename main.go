package main

import "caos/handler"

// Main
func main() {
	// Use the service requester interface to initialize node component
	var hn handler.IServiceRequester = handler.Node
	hn.Start()
}
