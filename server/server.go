package server

import (
	example "expamle"
	"github.com/gin-gonic/gin"
)

// Server holds the dependencies for an HTTP server.
type Server struct {
	Router        *gin.Engine
	config        *example.Configuration
	instrumenting serviceInstrumenting
}
