package ticker

type Service struct {
}

func New() *Service {
	return &Service{}
}

func (s *Service) Init() error {

	Init()

	return nil
}

func (s *Service) Run() error {
	// unix sock api
	Run()
	return nil
}

func (s *Service) Close() {
	Close()
}
