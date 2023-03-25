package app

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"
)

var once sync.Once

var log zerolog.Logger

func GetLogger(environment string) zerolog.Logger {
	once.Do(func() {
		zerolog.LevelFieldName = "severity"
		zerolog.TimeFieldFormat = time.RFC3339Nano

		var output io.Writer = zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}

		if environment != "local" {
			output = os.Stdout
		}

		log = zerolog.New(output).
			Level(zerolog.Level(zerolog.InfoLevel)).
			With().
			Timestamp().
			Logger()
	})

	return log
}

func GetRequestLogger(environment string, gcpProjectID string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := GetLogger(environment)

			var trace string
			if gcpProjectID != "" {
				traceHeader := r.Header.Get("X-Cloud-Trace-Context")
				traceParts := strings.Split(traceHeader, "/")
				if len(traceParts) > 0 && len(traceParts[0]) > 0 {
					trace = fmt.Sprintf("projects/%s/traces/%s", gcpProjectID, traceParts[0])
				}
			}

			ctx := context.WithValue(r.Context(), "trace", trace)

			r = r.WithContext(ctx)

			logger.UpdateContext(func(c zerolog.Context) zerolog.Context {
				return c.Str("logging.googleapis.com/trace", trace)
			})

			h := hlog.NewHandler(logger)

			accessHandler := hlog.AccessHandler(
				func(r *http.Request, status, size int, duration time.Duration) {
					hlog.FromRequest(r).Info().
						Str("method", r.Method).
						Stringer("url", r.URL).
						Int("status_code", status).
						Int("response_size_bytes", size).
						Dur("elapsed_ms", duration).
						Msg("incoming request")
				},
			)

			if environment == "local" {
				h(accessHandler(next)).ServeHTTP(w, r)
			} else {
				h(next).ServeHTTP(w, r)
			}
		})
	}
}
