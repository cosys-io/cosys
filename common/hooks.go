package common

type BootstrapHook struct {
	uid  string
	hook func(*Cosys) error
}

func (b BootstrapHook) Call(cosys *Cosys) error {
	return b.hook(cosys)
}

func (b BootstrapHook) String() string {
	return b.uid
}

type CleanupHook struct {
	uid  string
	hook func(*Cosys) error
}

func (c CleanupHook) Call(cosys *Cosys) error {
	return c.hook(cosys)
}

func (c CleanupHook) String() string {
	return c.uid
}
