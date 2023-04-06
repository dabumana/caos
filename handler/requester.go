// Package handler section
package handler

import (
	"caos/service"
)

// IServiceRequester - Service requester interface API
type IServiceRequester interface {
	Start()
}

// Node - Global node service for handler
var Node service.Node
