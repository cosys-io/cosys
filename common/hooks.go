package common

// BootstrapHook is a hook called during the bootstrap stage of the cosys app.
type BootstrapHook func(*Cosys) error

// CleanupHook is a hook called during the cleanup stage of the cosys app.
type CleanupHook func(*Cosys) error
