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
	if options == nil {
		return
	}

	options.deleteIfExists = true
}

func SkipIfExists(options *genOptions) {
	if options == nil {
		return
	}

	options.skipIfExists = true
}

func GenHeadOnly(options *genOptions) {
	if options == nil {
		return
	}

	options.genHeadOnly = true
}
