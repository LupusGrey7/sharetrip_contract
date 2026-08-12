package model

// ServiceCodeType — code of service from contract_management.services dictionary.
type ServiceCodeType string

const (
	ServiceCodeTripCreation     ServiceCodeType = "trip_creation"
	ServiceCodeTripParticipants ServiceCodeType = "trip_participants"
	ServiceCodeNotifications    ServiceCodeType = "notifications"
	ServiceCodePremiumSupport   ServiceCodeType = "premium_support"
)

// GetAvailableOfferingByCompanyIDRequest — enter for Fiber ParamsParser and usecase.
type GetAvailableOfferingByCompanyIDRequest struct {
	CompanyID   int             `params:"companyId" validate:"required,min=1"`
	ServiceCode ServiceCodeType `params:"serviceCode" validate:"required,oneof=trip_creation trip_participants notifications premium_support"` // TODO: add validation for service code
}

// GetAvailableResultResponse — result of availability usecase (HTTP-form — in http/dto).
type GetAvailableResultResponse struct {
	CompanyID   int
	ServiceCode ServiceCodeType
	Allowed     bool
	Reason      string
}
