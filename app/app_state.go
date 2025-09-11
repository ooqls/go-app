package app

type AppState struct {
	RegistryInitialized bool
	JWTInitialized bool
	RSAInitialized bool
	LoggingAPIInitialized bool
	HTTPInitialized bool
	GinInitialized bool
	DocsInitialized bool
	TLSInitialized bool
	SQLInitialized bool
	SQLSeeded bool
	Healthy bool
	Running bool
}