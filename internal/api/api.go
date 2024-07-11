package api

type Health struct{}

func (h Health) HealthCheck() string {
	return "OK!"
}
