package handler

import (
	"caos/service"
)

/* API Interface */
type IServiceRequester interface {
	// Start
	Start(sandboxMode bool)
}

// Global node
var Node service.NodeService
