package cookie

func GetDefaultCookiesDomain(domain string) string {
	if domain == "" {
		return "localhost"
	}

	return domain
}

func GetDefaultCookiesSameSite(cookiesSameSite string) string {
	if cookiesSameSite == "true" {
		return "Strict"
	}

	// Default to "None" if empty
	if cookiesSameSite == "" {
		return "None"
	}

	// Otherwise, return the provided value
	return cookiesSameSite
}
