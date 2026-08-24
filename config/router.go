package config

// ======================================================================================
// Router represents the routing and CORS (Cross-Origin Resource Sharing) configuration.
// ======================================================================================.
type Router struct {
	// AllowOriginsFunc custom function to dynamically evaluate allowed origins
	AllowOriginsFunc AllowOriginFunc

	// AllowOrigins list of allowed CORS origins
	AllowOrigins []string

	// AllowMethods allowed HTTP methods for CORS requests
	AllowMethods []string

	// AllowHeaders allowed HTTP headers for CORS requests
	AllowHeaders []string

	// ExposeHeaders headers exposed to the client in responses
	ExposeHeaders []string

	// AllowCredentials indicates whether credentials are included in CORS requests
	AllowCredentials bool

	// MaxAge maximum time (in seconds) the preflight response can be cached
	MaxAge int
}

// ======================================================================================
// AllowOriginFunc defines a function signature for evaluating custom CORS origin policies.
// ======================================================================================.
type AllowOriginFunc func(origin string) bool

// ======================================================================================
// RouterDefault applies default fallback values to empty fields in the Router configuration.
// ======================================================================================.
func RouterDefault(router Router) *Router {
	// Set default router AllowOrigins to ["*"] if not specified
	if len(router.AllowOrigins) == 0 {
		router.AllowOrigins = []string{"*"}
	}

	// Set default router AllowMethods if not specified
	if len(router.AllowMethods) == 0 {
		router.AllowMethods = []string{
			"GET",
			"POST",
			"HEAD",
			"PUT",
			"DELETE",
			"PATCH",
		}
	}

	// Set default router AllowHeaders if not specified
	if len(router.AllowHeaders) == 0 {
		router.AllowHeaders = []string{
			"Content-Type",
			"Content-Length",
			"Accept-Encoding",
			"X-CSRF-Token",
			"Authorization",
			"Accept",
			"Origin",
			"Cache-Control",
			"X-Requested-With",
		}
	}

	// Set default router ExposeHeaders if not specified
	if len(router.ExposeHeaders) == 0 {
		router.ExposeHeaders = []string{
			"Content-Length",
		}
	}

	// Set default router MaxAge to 1800 if not specified
	if router.MaxAge == 0 {
		router.MaxAge = 1800
	}

	return &router
}
