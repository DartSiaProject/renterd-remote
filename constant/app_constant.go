package constants

// Constant Api
const (
	DefaultConfigFile = "./config.cnf"

	WelcomeMessage               = "Welcome to your renterd-remote interface."
	GetStartMessage              = "To get started, please set the basic application parameters."
	ChangeCredentials            = "To change credentials, please set the following parameters."
	ChangeListenInterface        = "To change ip interface listener, please select one of the following parametres."
	PasswordInput                = "Please type your password to login through the mobile application : "
	PasswordConfirmationInput    = "Please re-type your password : "
	RenterdPasswordInput         = "Please type the renterd's password : "
	IpInterfaceHandlerInput      = "Please select an interface to handle request for mobile app: "
	EmailInput                   = "Please type an email address  : "
	EmailError                   = "You must enter a valid email address"
	PasswordError                = "Your password must be at least 8 characters long, with at least one upper, one lower case letter and one special caracter."
	ErrorSavingEnvFile           = "Error saving .env file"
	InvalidCredentials           = "Invalid credentials"
	OperationSuccessMessage      = "Operation Success"
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
	InvalidCommand               = "invalid command"
	InitAppNeeded                = "Error: You need to first initialise the App.\nPlease start the App without the credentials command"
	Error                        = "error"
	AllIpInterfacesMessage       = "All Interfaces"
	IncorrectParams              = "incorrect parameters"
	SliceError                   = "slice bounds out of range"
)

func DefaultParams() map[string]string {
	return map[string]string{
		//"SERVER_ADDRESS":  "",
		"SERVER_PORT":     "8000",
		"GIN_MODE":        "release",
		"RENTERD_ADDRESS": "localhost:9980",
	}
}
