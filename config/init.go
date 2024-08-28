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
	err := godotenv.Load()

	if err != nil {
		// Generate a JWT secret that is 64 characters long with 10 digits, 10 symbols,
		// allowing upper and lower case letters, disallowing repeat characters.
		secret, err := password.Generate(16, 4, 0, false, false)
		if err != nil {
			log.Fatal(err)
		}

		defaultParams := constants.DefaultParams()
		defaultParams["JWT_SECRET"] = secret

		err1 := godotenv.Write(defaultParams, "./.env")
		if err1 != nil {
			log.Fatal(constants.ErrorSavingEnvFile)
		}
	}
}

func InitApp() {
	fmt.Println(constants.WelcomeMessage)
	answers := struct {
		Email           string
		Password        string
		RenterdPassword string
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
	}}, &answers)

	fmt.Println("")

	if err != nil {
		log.Fatal(err)
		return
	}

	myEnv, _ := godotenv.Read()
	myEnv["USER_EMAIL"] = strings.ToLower(answers.Email)
	myEnv["USER_KEY"] = utils.CreateSecretKey(answers.Email, answers.Password)
	myEnv["USER_IV"] = utils.CreateIV(answers.Email, answers.Password)
	myEnv["RENTERD_KEY"] = answers.RenterdPassword

	err1 := godotenv.Write(myEnv, "./.env")
	if err1 != nil {
		log.Fatal(constants.ErrorSavingEnvFile)
	}
}
