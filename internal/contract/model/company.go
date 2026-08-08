package model

// ServiceCodeType — код услуги из словаря contract_management.services.
type ServiceCodeType string

const (
	ServiceCodeTripCreation     ServiceCodeType = "trip_creation"
	ServiceCodeTripParticipants ServiceCodeType = "trip_participants"
	ServiceCodeNotifications    ServiceCodeType = "notifications"
	ServiceCodePremiumSupport   ServiceCodeType = "premium_support"
)

// GetAvailableServiceByCompanyIDRequest — вход usecase (не HTTP DTO).
type GetAvailableServiceByCompanyIDRequest struct {
	CompanyID   int             `validate:"required,min=1"`
	ServiceCode ServiceCodeType `validate:"required,oneof=trip_creation trip_participants notifications premium_support"`
}

// GetAvailableResultResponse — результат availability usecase (HTTP-форма — в http/dto).
type GetAvailableResultResponse struct {
	CompanyID   int
	ServiceCode ServiceCodeType
	Allowed     bool
	Reason      string
}
