package gen

type Action interface {
	Act() error
}

type genOptions struct {
	deleteIfExists bool
	skipIfExists   bool
	genHeadOnly    bool
}

type genOption func(*genOptions)

func DeleteIfExists(options *genOptions) {
	options.deleteIfExists = true
}

func SkipIfExists(options *genOptions) {
	options.skipIfExists = true
}

func GenHeadOnly(options *genOptions) {
	options.genHeadOnly = true
}
