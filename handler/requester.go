// Requester component
package handler

import (
	"caos/service"
)

// IServiceRequester - Service requester interface API
type IServiceRequester interface {
	Start(sandboxMode bool)
}

// Node - Global node service for handler
var Node service.NodeService
