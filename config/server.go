package config

// ======================================================================================
// Server represents the HTTP server network binding options.
// ======================================================================================.
type Server struct {
	// Host network interface host address to bind to
	Host string

	// Port network port number for the HTTP server to listen on
	Port string
}

// ======================================================================================
// CookiesDefault applies default fallback values to empty fields in the Cookies configuration.
// ======================================================================================.
func ServerDefault(srv Server) *Server {
	// Set default server port to 8080 if not specified
	if srv.Port == "" {
		srv.Port = "8080"
	}

	return &srv
}
