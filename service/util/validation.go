package service

import (
	"log"
)

// IsContextValid - Client context validation
func IsContextValid(service Agent) bool {
	if service.ctx == nil {
		log.Fatalln("Context NOT found")
		return false
	} else if service.client == nil {
		log.Fatalln("Client NOT found")
		return false
	}

	return true
}
