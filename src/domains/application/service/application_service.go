package service

type ApplicationServicer interface{}

type ApplicationService struct{}

type ApplicationServiceOptions func(s *ApplicationService)

func NewApplication(opts ...ApplicationServiceOptions) ApplicationServicer {
	svc := &ApplicationService{}

	for _, opt := range opts {
		opt(svc)
	}

	return svc
}
