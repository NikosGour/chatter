package common

import (
	"github.com/NikosGour/chatter/internal/projectpath"
	"github.com/NikosGour/logging/log"
	"github.com/go-playground/validator/v10"
	"github.com/joho/godotenv"
)

var (
	Dotenv   map[string]string
	Validate *validator.Validate
)

const (
	envFile = ".env"
)

const (
	EnvHOST_ADDRESS = "HOST_ADDRESS"
	EnvPORT         = "PORT"

	EnvPOSTGRES_HOST_ADDRESS  = "POSTGRES_HOST_ADDRESS"
	EnvPOSTGRES_PORT          = "POSTGRES_PORT"
	EnvPOSTGRES_USER          = "POSTGRES_USER"
	EnvPOSTGRES_ROOT_PASSWORD = "POSTGRES_ROOT_PASSWORD"
	EnvPOSTGRES_DB            = "POSTGRES_DB"
)

func InitDotenv() {
	dotenv, err := godotenv.Read(projectpath.RootFile(envFile))
	if err != nil {
		log.Fatal("%s", err)
	}
	Validate = validator.New()
	Dotenv = dotenv
}
