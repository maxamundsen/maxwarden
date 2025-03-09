package constants

const (
	SESSION_COOKIE_NAME            = "_maxwarden_session"
	SESSION_COOKIE_EXPIRY_DAYS int = 100
	SESSION_COOKIE_ENTROPY     int = 33

	IDENTITY_COOKIE_NAME        string = "_maxwarden_identity"
	IDENTITY_COOKIE_EXPIRY_DAYS int    = 30
	IDENTITY_TOKEN_EXPIRY_DAYS  int    = 30
	IDENTITY_COOKIE_ENTROPY     int    = 33
	IDENTITY_LOGIN_PATH         string = "/auth/login"
	IDENTITY_LOGOUT_PATH        string = "/auth/logout"
	IDENTITY_DEFAULT_PATH       string = "/app"
	IDENTITY_AUTH_REDIRECT      bool   = true
	IDENTITY_AUTH_KEY           string = "CORRECT_HORSE_BATTERY_STAPLE"

	// This key is NOT used for the hashing of passwords, or secure session data over the wire.
	// It is ONLY used for performing quick file and string hashes, where security is not a factor.
	DATA_HASH_KEY string = "01234567890123456789012345678901"

	PASSWORD_MIN_LENGTH         int = 8
	PASSWORD_REQUIRED_UPPERCASE int = 1
	PASSWORD_REQUIRED_LOWERCASE int = 1
	PASSWORD_REQUIRED_NUMBERS   int = 1
	PASSWORD_REQUIRED_SYMBOLS   int = 0

	MAX_LOGIN_ATTEMPTS int = 5
)