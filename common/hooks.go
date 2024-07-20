package common

// BootstrapHook is a hook called during the bootstrap stage of the cosys app.
type BootstrapHook func(*Cosys) error

// Call calls the bootstrap hook.
func (b BootstrapHook) Call(cosys *Cosys) error {
	return b(cosys)
}

// CleanupHook is a hook called during the cleanup stage of the cosys app.
type CleanupHook func(*Cosys) error

// Call calls the cleanup hook.
func (c CleanupHook) Call(cosys *Cosys) error {
	return c(cosys)
}
