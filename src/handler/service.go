// Package handler section
package handler

import (
	"caos/service"
)

// ServiceRequester - Service requester interface API
type ServiceRequester interface {
	Start()
}

// Node - Global node service for handler
var Node service.Node
