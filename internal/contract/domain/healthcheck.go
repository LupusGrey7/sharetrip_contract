package domain

// HealthcheckOutput — output from usecase/service to API presenter.
type HealthcheckOutput struct {
	Status  string
	Message string
}
