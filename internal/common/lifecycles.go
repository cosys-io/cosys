package common

type Event struct {
	Params *QEParams
	Result any
	State  any
}

func NewEvent(params *QEParams) *Event {
	return &Event{
		params,
		nil,
		nil,
	}
}

type LifecycleFunc struct {
	EventType string
	Func      func(*Event) error
}

type Lifecycles struct {
	LifecycleFuncs []*LifecycleFunc
}

func NewLifecycles(lifecycleFuncs ...*LifecycleFunc) *Lifecycles {
	return &Lifecycles{
		lifecycleFuncs,
	}
}

func (l *Lifecycles) Act(eventType string, event *Event) error {
	for _, lifecycleFunc := range l.LifecycleFuncs {
		if lifecycleFunc.EventType == eventType {
			if err := lifecycleFunc.Func(event); err != nil {
				return err
			}
		}
	}

	return nil
}

func BeforeFindOne(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeFindOne",
		f,
	}
}

func AfterFindOne(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterFindOne",
		f,
	}
}

func BeforeFindMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeFindMany",
		f,
	}
}

func AfterFindMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterFindMany",
		f,
	}
}

func BeforeCreate(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeCreate",
		f,
	}
}

func AfterCreate(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterCreate",
		f,
	}
}

func BeforeCreateMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeCreateMany",
		f,
	}
}

func AfterCreateMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterCreateMany",
		f,
	}
}

func BeforeUpdate(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeUpdate",
		f,
	}
}

func AfterUpdate(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterUpdate",
		f,
	}
}

func BeforeUpdateMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeUpdateMany",
		f,
	}
}

func AfterUpdateMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterUpdateMany",
		f,
	}
}

func BeforeDelete(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeDelete",
		f,
	}
}

func AfterDelete(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterDelete",
		f,
	}
}

func BeforeDeleteMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"beforeDeleteMany",
		f,
	}
}

func AfterDeleteMany(f func(*Event) error) *LifecycleFunc {
	return &LifecycleFunc{
		"afterDeleteMany",
		f,
	}
}
