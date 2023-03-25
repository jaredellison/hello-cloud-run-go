package main

import (
	"fmt"
	app "go-service-template"
	"log"
	"net/http"

	"github.com/kelseyhightower/envconfig"
)

type EnvConfig struct {
	Port        int
	Environment string
	Project     string
}

func main() {
	env := EnvConfig{}
	err := envconfig.Process("service", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	logger := app.GetLogger(env.Environment)

	service := app.NewService(env.Environment, env.Project)

	logger.Info().Msgf("🚀 Server up and listening on http://localhost:%d 🚀", env.Port)

	logger.
		Fatal().
		Err(http.ListenAndServe(fmt.Sprintf(":%d", env.Port), service.Router)).
		Msg("Server shut down")
}
