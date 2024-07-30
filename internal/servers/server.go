package servers

import "context"

// Server for interaction with remote clipboard service
type Server interface {
	// Starts the server and waits until server will be stopped with context.
	// Returns error if error happened otherwise nil.
	Start(ctx context.Context) error
}
