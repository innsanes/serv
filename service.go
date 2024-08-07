package serv

type Service struct {
}

func (s *Service) BeforeServe() error {
	return nil
}

func (s *Service) AfterServe() error {
	return nil
}

func (s *Service) Serve() error {
	return nil
}

func (s *Service) BeforeStop() error {
	return nil
}
