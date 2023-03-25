package app

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
)

type Service struct {
	Router *chi.Mux
}

func NewService(environment, gcpProjectID string) Service {
	srvc := Service{}
	srvc.Router = chi.NewRouter()

	srvc.Router.Use(GetRequestLogger(environment, gcpProjectID))
	srvc.routes()

	return srvc
}

func (s *Service) routes() {
	s.Router.Get("/", s.handleRoot())
	s.Router.Get("/info", s.handleInfo())
	s.Router.Get("/warn", s.handleWarn())
	s.Router.Get("/error", s.handleError())
}

func (s *Service) handleRoot() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())

		logger.Info().Msg("logging in root handler")

		w.Write([]byte("hello"))
	}
}

func (s *Service) handleInfo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())

		logger.Info().Msg("logging in info handler")

		w.Write([]byte("info"))
	}
}

func (s *Service) handleWarn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())

		logger.Warn().Msg("logging in warn handler")

		w.Write([]byte("warn"))
	}
}

func (s *Service) handleError() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := zerolog.Ctx(r.Context())

		logger.Error().Msg("logging in error handler")

		w.Write([]byte("error"))
	}
}
