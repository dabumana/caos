package handler

import (
	"caos/service"
)

/* API Interface */
type IServiceRequester interface {
	// Start
	Start()
}

// Global node
var Node service.NodeService
