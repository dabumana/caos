// Package handler section
package handler

import (
	"caos/service"
)

// Node - Global node service for handler
var Node service.Node

// ServiceRequester - Service requester interface API
type ServiceRequester interface {
	Init()
}
