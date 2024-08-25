package constants

// Constant Api
const (
	WelcomeMessage               = "Welcome to your de renterd-remote interface. \nTo get started, please set the basic application parameters."
	PasswordInput                = "Please type your password : "
	PasswordConfirmationInput    = "Please re-type your password : "
	RenterdPasswordInput         = "Please type your renterd's password : "
	EmailInput                   = "Please type your email address : "
	EmailError                   = "You must enter a valid email address"
	PasswordError                = "Your password must be at least 8 characters long, with at least one upper, one lower case letter and one special caracter."
	ErrorSavingEnvFile           = "Error saving .env file"
	InvalidCredentials           = "Invalid credentials"
	SuccessMessage               = "Success"
	Unauthorized                 = "Unauthorized"
	BadRequest                   = "Bad Request"
	InternalServerError          = "An internal error occur"
	HeaderRequestEncryptionError = "Error during encryption of the header data"
	HeaderRequestDecryptionError = "Error during decryption of the header data"
	ParamsRequestEncryptionError = "Error during encryption of the params data"
	ParamsRequestDecryptionError = "Error during decryption of the params data"
	BodyRequestEncryptionError   = "Error during encryption of the body data"
	BodyRequestDecryptionError   = "Error during decryption of the body data"
)

func DefaultParams() map[string]string {
	return map[string]string{
		"SERVER_ADDRESS":  "localhost",
		"SERVER_PORT":     "8080",
		"GIN_MODE":        "debug",
		"RENTERD_ADDRESS": "localhost:9980",
	}
}
