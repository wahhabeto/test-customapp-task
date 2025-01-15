package http

func (s *server) MountRouter() error {

	apiV1 := s.app.Group("/api/v1")

	apiV1.Get("/health", s.HealthCheck)

	apiV1.Get("/get", s.GetGeneratedNumber)

	return nil
}
