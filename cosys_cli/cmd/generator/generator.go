package gen

type Generator []Action

func (g Generator) Generate() error {
	for _, action := range g {
		if err := action.Act(); err != nil {
			return err
		}
	}

	return nil
}

func NewGenerator(actions ...Action) Generator {
	return actions
}
