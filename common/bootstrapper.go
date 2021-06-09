package common

// start up of the environment
func StartUp() {
	// Initialize private/public keys for JWT authentication
	initKeys()

	// Start a SQL DB session to e used by repositories
	createOracleDbSession()
}
