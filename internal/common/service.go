package common

type Service struct {
	Functions map[string]Function
}

func (s *Service) Function(uid string) Function {
	return s.Functions[uid]
}

func NewService(functions map[string]Function) *Service {
	return &Service{
		functions,
	}
}

type Function func(Cosys) ServiceFunction

type ServiceFunction any
