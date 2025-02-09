package config

import (
	"fmt"
	constants "renterd-remote/constant"
	"renterd-remote/utils"
	validationHelpers "renterd-remote/utils"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/joho/godotenv"
	"github.com/sethvargo/go-password/password"
	log "github.com/sirupsen/logrus"
)

func LoadEnvVariables() {
	//Load .env file
	err := godotenv.Load(constants.DefaultConfigFile)

	if err != nil {
		// Generate a JWT secret that is 64 characters long with 10 digits, 10 symbols,
		// allowing upper and lower case letters, disallowing repeat characters.
		secret, err := password.Generate(16, 4, 0, false, false)
		if err != nil {
			log.Fatal(err)
		}

		defaultParams := constants.DefaultParams()
		defaultParams["JWT_SECRET"] = secret

		err1 := godotenv.Write(defaultParams, constants.DefaultConfigFile)
		if err1 != nil {
			log.Fatal(constants.ErrorSavingEnvFile)
		}
	}
}

func InitApp() {
	fmt.Println(constants.WelcomeMessage)
	fmt.Println(constants.GetStartMessage)
	answers := struct {
		Email              string
		Password           string
		ConfirmPassword    string
		RenterdPassword    string
		IpInterfaceHandler string
	}{}

	err := survey.Ask([]*survey.Question{{
		Name: "email",
		Prompt: &survey.Input{
			Message: constants.EmailInput,
		},
		Validate: validationHelpers.ValidateEmail,
	}, {
		Name: "password",
		Prompt: &survey.Password{
			Message: constants.PasswordInput,
		},
		Validate: validationHelpers.ValidatePassword,
	}, {
		Name: "Renterdpassword",
		Prompt: &survey.Password{
			Message: constants.RenterdPasswordInput,
		},
		Validate: survey.Required,
	}, {
		Name: "IpInterfaceHandler",
		Prompt: &survey.Select{
			Message: constants.IpInterfaceHandlerInput,
			Options: utils.GetLocalIP(),
			Default: constants.AllIpInterfacesMessage,
		},
	}}, &answers)

	if err != nil {
		log.Fatal(err)
		return
	}

	myEnv, _ := godotenv.Read(constants.DefaultConfigFile)
	myEnv["USER_EMAIL"] = strings.ToLower(answers.Email)
	myEnv["USER_KEY"] = utils.CreateSecretKey(answers.Email, answers.Password)
	myEnv["USER_IV"] = utils.CreateIV(answers.Email, answers.Password)
	myEnv["RENTERD_KEY"] = answers.RenterdPassword
	if answers.IpInterfaceHandler == constants.AllIpInterfacesMessage {
		myEnv["SERVER_ADDRESS"] = ""
	} else {
		myEnv["SERVER_ADDRESS"] = answers.IpInterfaceHandler
	}

	err1 := godotenv.Write(myEnv, constants.DefaultConfigFile)
	if err1 != nil {
		log.Fatal(constants.ErrorSavingEnvFile)
	}

	fmt.Printf("\n%s\n", constants.OperationSuccessMessage)
}
