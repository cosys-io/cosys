package common

type BootstrapHook struct {
	uid  string
	hook func(*Cosys) error
}

func (b *BootstrapHook) Call(cosys *Cosys) error {
	return b.hook(cosys)
}

type CleanupHook struct {
	uid  string
	hook func(*Cosys) error
}

func (c *CleanupHook) Call(cosys *Cosys) error {
	return c.hook(cosys)
}
