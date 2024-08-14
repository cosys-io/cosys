package common

// Server is a core service for the external API.
type Server interface {
	Start() error
}
