package app

type AppState struct {
	RegistryInitialized bool
	JWTInitialized bool
	RSAInitialized bool
	SQLInitialized bool
	SQLSeeded bool
}