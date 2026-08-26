package domain

// ServiceCode — code from contract_management.services dictionary.
type ServiceCode string

const (
	ServiceCodeTripStart        ServiceCode = "trip_start" // ShareTrip moveTripPublished-ToStarted
	ServiceCodeTripCreation     ServiceCode = "trip_creation"
	ServiceCodeTripParticipants ServiceCode = "trip_participants"
	ServiceCodeNotifications    ServiceCode = "notifications"
	ServiceCodePremiumSupport   ServiceCode = "premium_support"
)

// GetAvailableOfferingByCompanyIDInput — input for service/usecase.
type GetAvailableOfferingByCompanyIDInput struct {
	CompanyID   int
	ServiceCode ServiceCode
}

// AvailabilityEntity — result read from storage; it stays below the usecase boundary.
type AvailabilityEntity struct {
	CompanyID   int
	ServiceCode ServiceCode
	Allowed     bool
	Reason      string
}

// AvailabilityOutput — output from usecase/service to API presenter.
type AvailabilityOutput struct {
	CompanyID   int
	ServiceCode ServiceCode
	Allowed     bool
	Reason      string
}

func AvailabilityEntityToOutput(entity *AvailabilityEntity) *AvailabilityOutput {
	if entity == nil {
		return nil
	}
	return &AvailabilityOutput{
		CompanyID:   entity.CompanyID,
		ServiceCode: entity.ServiceCode,
		Allowed:     entity.Allowed,
		Reason:      entity.Reason,
	}
}
