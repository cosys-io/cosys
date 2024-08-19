package gen

// Action performs a file system change for code generation.
type Action interface {
	Act() error
}

// genOptions are the configurations for actions.
type genOptions struct {
	deleteIfExists bool
	skipIfExists   bool
	genHeadOnly    bool
}

// genOption is a configuration for actions.
type genOption func(*genOptions)

// DeleteIfExists is a configuration for actions,
// specifying that the given file should be deleted if exists.
var DeleteIfExists genOption = func(options *genOptions) {
	if options == nil {
		return
	}

	options.deleteIfExists = true
}

// SkipIfExists is a configuration for actions,
// specifying that the action should be skipped if the given file exists.
var SkipIfExists genOption = func(options *genOptions) {
	if options == nil {
		return
	}

	options.skipIfExists = true
}

// GenHeadOnly is a configuration for actions,
// specifying that only the head of the path should be created.
var GenHeadOnly genOption = func(options *genOptions) {
	if options == nil {
		return
	}

	options.genHeadOnly = true
}
