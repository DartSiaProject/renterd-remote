package update

import (
	"fmt"
	"log"
	"os"
	constants "renterd-remote/constant"
	"renterd-remote/utils"
	"strings"

	validationHelpers "renterd-remote/utils"

	"github.com/AlecAivazis/survey/v2"
	"github.com/joho/godotenv"
)

func Config() (stopApp bool, Error error) {
	cmd := ""
	if len(os.Args) < 2 {
		return false, nil
	} else if len(os.Args) == 2 {
		cmd = os.Args[1]
	} else if len(os.Args) > 2 {
		fmt.Println(constants.InvalidCommand)
		return true, fmt.Errorf(constants.InvalidCommand)
	}

	if cmd == "credentials" {
		if os.Getenv("USER_EMAIL") != "" || os.Getenv("USER_KEY") != "" {
			resetCredentials()
			return true, nil
		} else {
			fmt.Println(constants.InitAppNeeded)
			return true, fmt.Errorf(constants.InitAppNeeded)
		}
	} else {
		fmt.Println(constants.InvalidCommand)
		return true, fmt.Errorf(constants.InvalidCommand)
	}
}

func resetCredentials() {
	fmt.Println(constants.WelcomeMessage)
	fmt.Println(constants.ChangeCredentials)
	answers := struct {
		Email    string
		Password string
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
	}}, &answers)

	if err != nil {
		log.Fatal(err)
		return
	}

	myEnv, _ := godotenv.Read(constants.DefaultConfigFile)
	myEnv["USER_EMAIL"] = strings.ToLower(answers.Email)
	myEnv["USER_KEY"] = utils.CreateSecretKey(answers.Email, answers.Password)
	myEnv["USER_IV"] = utils.CreateIV(answers.Email, answers.Password)

	err1 := godotenv.Write(myEnv, constants.DefaultConfigFile)
	if err1 != nil {
		log.Fatal(constants.ErrorSavingEnvFile)
	}

	fmt.Println("\n\n", constants.OperationSuccessMessage)
}
