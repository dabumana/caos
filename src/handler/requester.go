// Package handler section
package handler

import "caos/service"

var Node service.Node

// IServiceRequester - Service requester interface API
type IServiceRequester interface {
	Init()
}
