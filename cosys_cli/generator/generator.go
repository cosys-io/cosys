package gen

// Generator is a sequence of actions for code generation.
type Generator []Action

// Generate performs the sequence of actions.
func (g Generator) Generate() error {
	for _, action := range g {
		if err := action.Act(); err != nil {
			return err
		}
	}

	return nil
}

// NewGenerator returns a new Generator from a sequence of actions.
func NewGenerator(actions ...Action) Generator {
	return actions
}
