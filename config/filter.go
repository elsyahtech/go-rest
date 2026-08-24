package config

// ======================================================================================
// Filter represents global filtering or security mechanism options.
// ======================================================================================.
type Filter struct {
	// AuthenticationType strategy or type used for request authentication
	AuthenticationType string
}

// ======================================================================================
// FilterDefault applies default fallback values to empty fields in the Filter configuration.
// ======================================================================================.
func FilterDefault(filter Filter) *Filter {
	// Set default AuthenticationType to JWT if not specified
	if filter.AuthenticationType == "" {
		filter.AuthenticationType = "JWT"
	}

	return &filter
}
