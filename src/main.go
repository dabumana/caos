// Main package
package main

import "caos/handler"

// Main
func main() {
	// Use the service requester interface to initialize node component
	var service handler.IServiceRequester = &handler.Node
	service.Init()
}
